package kubernetes

import (
	"context"
	"testing"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestGetCustomResourcesDynamicColumns_OwnerReferencesCollision(t *testing.T) {
	tests := []struct {
		name          string
		specSchema    v1.JSONSchemaProps
		statusSchema  v1.JSONSchemaProps
		wantColumns   []string
		unwantColumns []string
	}{
		{
			name: "spec and status properties named ownerReferences are prefixed",
			specSchema: v1.JSONSchemaProps{
				Properties: map[string]v1.JSONSchemaProps{
					"ownerReferences": {Type: "string"},
					"replicas":        {Type: "integer"},
				},
			},
			statusSchema: v1.JSONSchemaProps{
				Properties: map[string]v1.JSONSchemaProps{
					"ownerReferences": {Type: "string"},
				},
			},
			wantColumns:   []string{"spec_owner_references", "status_owner_references", "replicas"},
			unwantColumns: []string{"owner_references"},
		},
		{
			name:         "no collision when schemas don't declare ownerReferences",
			specSchema:   v1.JSONSchemaProps{Properties: map[string]v1.JSONSchemaProps{"replicas": {Type: "integer"}}},
			statusSchema: v1.JSONSchemaProps{Properties: map[string]v1.JSONSchemaProps{}},
			wantColumns:  []string{"replicas"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columns := getCustomResourcesDynamicColumns(context.Background(), tt.specSchema, tt.statusSchema)

			got := make(map[string]bool, len(columns))
			for _, c := range columns {
				got[c.Name] = true
			}

			for _, want := range tt.wantColumns {
				if !got[want] {
					t.Errorf("getCustomResourcesDynamicColumns() columns = %v, want to contain %q", got, want)
				}
			}
			for _, unwant := range tt.unwantColumns {
				if got[unwant] {
					t.Errorf("getCustomResourcesDynamicColumns() columns = %v, want NOT to contain %q", got, unwant)
				}
			}
		})
	}
}

func TestCRDResourceInfo_OwnerReferences(t *testing.T) {
	tests := []struct {
		name       string
		references []metav1.OwnerReference
		wantLen    int
		wantKind   string
	}{
		{
			name: "controller owner reference is preserved",
			references: []metav1.OwnerReference{
				{APIVersion: "platform.totvs.app/v1", Kind: "ComponentInstallation", Name: "cert-manager", UID: "abc-123", Controller: boolPtr(true)},
			},
			wantLen:  1,
			wantKind: "ComponentInstallation",
		},
		{
			name:       "no owner references",
			references: nil,
			wantLen:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &unstructured.Unstructured{}
			u.SetOwnerReferences(tt.references)

			info := &CRDResourceInfo{OwnerReferences: u.GetOwnerReferences()}

			refs, ok := info.OwnerReferences.([]metav1.OwnerReference)
			if !ok {
				t.Fatalf("CRDResourceInfo.OwnerReferences type = %T, want []metav1.OwnerReference", info.OwnerReferences)
			}
			if len(refs) != tt.wantLen {
				t.Fatalf("len(OwnerReferences) = %d, want %d", len(refs), tt.wantLen)
			}
			if tt.wantLen == 0 {
				return
			}
			if refs[0].Kind != tt.wantKind {
				t.Errorf("OwnerReferences[0].Kind = %q, want %q", refs[0].Kind, tt.wantKind)
			}
		})
	}
}

func boolPtr(b bool) *bool { return &b }
