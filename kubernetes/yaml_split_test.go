package kubernetes

import (
	"reflect"
	"testing"
)

func TestYamlDocSeparator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple two documents",
			input:    "apiVersion: v1\nkind: Pod\n---\napiVersion: v1\nkind: Service",
			expected: []string{"apiVersion: v1\nkind: Pod\n", "\napiVersion: v1\nkind: Service"},
		},
		{
			name:     "separator at start of file",
			input:    "---\napiVersion: v1\nkind: Pod",
			expected: []string{"", "\napiVersion: v1\nkind: Pod"},
		},
		{
			name:     "separator with trailing whitespace",
			input:    "apiVersion: v1\nkind: Pod\n---  \napiVersion: v1\nkind: Service",
			expected: []string{"apiVersion: v1\nkind: Pod\n", "\napiVersion: v1\nkind: Service"},
		},
		{
			name:     "separator with trailing tab",
			input:    "apiVersion: v1\nkind: Pod\n---\t\napiVersion: v1\nkind: Service",
			expected: []string{"apiVersion: v1\nkind: Pod\n", "\napiVersion: v1\nkind: Service"},
		},
		{
			name:  "indented separator inside multiline string is NOT a split point",
			input: "apiVersion: v1\ndata:\n  install_info: |\n    ---\n    install_method:\n      tool: helm\nkind: ConfigMap",
			expected: []string{
				"apiVersion: v1\ndata:\n  install_info: |\n    ---\n    install_method:\n      tool: helm\nkind: ConfigMap",
			},
		},
		{
			name: "issue 341 exact scenario - ConfigMap with indented separator",
			input: `---
apiVersion: v1
data:
  install_info: |
    ---
    install_method:
      tool: helm
      tool_version: Helm
      installer_version: datadog-3.135.4
kind: ConfigMap
metadata:
  annotations:
    checksum/install_info: 22bff6a15fb7a4521a3b6a06f55f1fe8ca1570dab4d8bca9f437f28e5301e89a
  labels:
    app.kubernetes.io/instance: datadog
---
apiVersion: v1
kind: Service
metadata:
  name: my-service`,
			expected: []string{
				"",
				"\napiVersion: v1\ndata:\n  install_info: |\n    ---\n    install_method:\n      tool: helm\n      tool_version: Helm\n      installer_version: datadog-3.135.4\nkind: ConfigMap\nmetadata:\n  annotations:\n    checksum/install_info: 22bff6a15fb7a4521a3b6a06f55f1fe8ca1570dab4d8bca9f437f28e5301e89a\n  labels:\n    app.kubernetes.io/instance: datadog\n",
				"\napiVersion: v1\nkind: Service\nmetadata:\n  name: my-service",
			},
		},
		{
			name:     "three documents",
			input:    "---\napiVersion: v1\nkind: Pod\n---\napiVersion: v1\nkind: Service\n---\napiVersion: v1\nkind: Deployment",
			expected: []string{"", "\napiVersion: v1\nkind: Pod\n", "\napiVersion: v1\nkind: Service\n", "\napiVersion: v1\nkind: Deployment"},
		},
		{
			name:     "no separator - single document",
			input:    "apiVersion: v1\nkind: Pod\nmetadata:\n  name: test",
			expected: []string{"apiVersion: v1\nkind: Pod\nmetadata:\n  name: test"},
		},
		{
			name:     "dashes inside a value are not separators",
			input:    "apiVersion: v1\nkind: ConfigMap\ndata:\n  key: ---some-value---",
			expected: []string{"apiVersion: v1\nkind: ConfigMap\ndata:\n  key: ---some-value---"},
		},
		{
			name:     "four dashes on own line are not a separator",
			input:    "apiVersion: v1\nkind: Pod\n----\napiVersion: v1\nkind: Service",
			expected: []string{"apiVersion: v1\nkind: Pod\n----\napiVersion: v1\nkind: Service"},
		},
		{
			name:  "multiple indented separators inside strings",
			input: "apiVersion: v1\nkind: ConfigMap\ndata:\n  a: |\n    ---\n    foo: bar\n  b: |\n    ---\n    baz: qux",
			expected: []string{
				"apiVersion: v1\nkind: ConfigMap\ndata:\n  a: |\n    ---\n    foo: bar\n  b: |\n    ---\n    baz: qux",
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "only separator",
			input:    "---",
			expected: []string{"", ""},
		},
		{
			name:     "consecutive separators",
			input:    "---\n---\n---",
			expected: []string{"", "\n", "\n", ""},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := yamlDocSeparator.Split(tc.input, -1)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("yamlDocSeparator.Split() mismatch\n  input:    %q\n  got:      %q\n  expected: %q", tc.input, result, tc.expected)
			}
		})
	}
}
