package kubernetes

import (
	"os"
	"path/filepath"
	"testing"
)

const inlineKubeconfig = `apiVersion: v1
kind: Config
clusters:
  - name: my-cluster
    cluster:
      server: https://127.0.0.1:6443
contexts:
  - name: my-context
    context:
      cluster: my-cluster
      user: my-user
  - name: other-context
    context:
      cluster: my-cluster
      user: my-user
current-context: my-context
users:
  - name: my-user
    user:
      token: my-token
`

func TestPathOrContents(t *testing.T) {
	// A real file on disk, to exercise the existing-path branch.
	existing := filepath.Join(t.TempDir(), "kubeconfig")
	if err := os.WriteFile(existing, []byte(inlineKubeconfig), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	tests := []struct {
		name       string
		input      string
		wantInline bool
		wantErr    bool
	}{
		{
			name:       "empty string is not inline",
			input:      "",
			wantInline: false,
		},
		{
			name:       "inline kubeconfig YAML is detected",
			input:      inlineKubeconfig,
			wantInline: true,
		},
		{
			name:       "existing file path is not inline",
			input:      existing,
			wantInline: false,
		},
		{
			name:       "missing absolute path falls through to loader, not inline",
			input:      "/no/such/kubeconfig",
			wantInline: false,
		},
		{
			name:       "missing relative path is not inline",
			input:      "kubeconfig.yaml",
			wantInline: false,
		},
		{
			name:       "home-relative path is not inline",
			input:      "~/.kube/config",
			wantInline: false,
		},
		{
			name:       "multiline text without apiVersion is not inline",
			input:      "line one\nline two\n",
			wantInline: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, isInline, err := pathOrContents(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("pathOrContents(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if isInline != tt.wantInline {
				t.Errorf("pathOrContents(%q) isInline = %v, want %v", tt.input, isInline, tt.wantInline)
			}
		})
	}
}

func TestTryInlineKubeconfig(t *testing.T) {
	// A file path should not be treated as inline, so the caller falls through
	// to the standard loader.
	if _, ok, err := tryInlineKubeconfig("~/.kube/config", nil); err != nil || ok {
		t.Fatalf("tryInlineKubeconfig(file path) = ok %v, err %v; want ok false, err nil", ok, err)
	}

	// Inline kubeconfig should parse and use its own current-context.
	cfg, ok, err := tryInlineKubeconfig(inlineKubeconfig, nil)
	if err != nil || !ok {
		t.Fatalf("tryInlineKubeconfig(inline) = ok %v, err %v; want ok true, err nil", ok, err)
	}
	raw, err := cfg.RawConfig()
	if err != nil {
		t.Fatalf("RawConfig() error = %v", err)
	}
	if raw.CurrentContext != "my-context" {
		t.Errorf("CurrentContext = %q, want %q", raw.CurrentContext, "my-context")
	}

	// An explicit config_context should override the current-context.
	override := "other-context"
	cfg, ok, err = tryInlineKubeconfig(inlineKubeconfig, &override)
	if err != nil || !ok {
		t.Fatalf("tryInlineKubeconfig(inline, override) = ok %v, err %v; want ok true, err nil", ok, err)
	}
	raw, err = cfg.RawConfig()
	if err != nil {
		t.Fatalf("RawConfig() error = %v", err)
	}
	if raw.CurrentContext != override {
		t.Errorf("CurrentContext = %q, want %q", raw.CurrentContext, override)
	}
}
