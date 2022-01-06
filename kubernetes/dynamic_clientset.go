package kubernetes

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/internalversion/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

type CRDConfigInterface interface {
	List(opts meta_v1.ListOptions) (*CRDConfigList, error)
	// ...
}

type crdConfigClient struct {
	client rest.Interface
}

func (c *crdConfigClient) List(opts meta_v1.ListOptions) (*CRDConfigList, error) {
	result := CRDConfigList{}
	err := c.client.Get().Namespace("steampipe-cloud").Resource(CRDPlural).VersionedParams(&opts, scheme.ParameterCodec).Do(context.TODO()).Into(&result)
	// Resource("projects").
	// VersionedParams(&opts, scheme.ParameterCodec).
	// Do().
	// Into(&result)

	return &result, err
}

type CRDConfigV1Alpha1Interface interface {
	CRDConfigs() CRDConfigInterface
}

type CRDConfigV1Alpha1Client struct {
	restClient rest.Interface
}

func (c *CRDConfigV1Alpha1Client) CRDConfigs() CRDConfigInterface {
	return &crdConfigClient{
		client: c.restClient,
	}
}

func NewClient(cfg *rest.Config) (*CRDConfigV1Alpha1Client, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}
	config := *cfg
	config.GroupVersion = &SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme)
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &CRDConfigV1Alpha1Client{restClient: client}, nil
}
