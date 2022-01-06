package kubernetes

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	CRDPlural   string = "certificates"
	CRDGroup    string = "cert-manager.io"
	CRDVersion  string = "v1alpha2"
	FullCRDName string = CRDPlural + "." + CRDGroup
)

var SchemeGroupVersion = schema.GroupVersion{Group: CRDGroup, Version: CRDVersion}

type CRDConfig struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               CRDConfigSpec   `json:"spec"`
	Status             CRDConfigStatus `json:"status,omitempty"`
}

type CRDConfigSpec struct {
	CommonName interface{} `json:"commonName"`
	DnsNames   interface{} `json:"dnsNames"`
	Duration   interface{} `json:"duration"`
}

type CRDConfigStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

type CRDConfigList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []CRDConfig `json:"items"`
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&CRDConfig{},
		&CRDConfigList{},
	)
	meta_v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
