package kubernetes

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/turbot/steampipe-plugin-sdk/connection"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func GetNewClientset(ctx context.Context, connectionManager *connection.Manager) (*kubernetes.Clientset, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("GetNewClientset")

	// have we already created and cached the session?
	serviceCacheKey := "k8s" //should probably per connection/context keys...

	if cachedData, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		logger.Warn("!!!! Clientset Found in Cache !!!!")
		return cachedData.(*kubernetes.Clientset), nil
	}

	// Instantiate loader for kubeconfig file.
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Get a rest.Config from the kubeconfig file.
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(restconfig)
	if err != nil {
		return nil, err
	}

	connectionManager.Cache.Set(serviceCacheKey, clientset)
	if _, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		logger.Warn("!!!! Clientset added to cache !!!!")
	} else {
		logger.Warn("!!!! Clientset NOT Found in Cache after adding !!!!")
	}

	return clientset, err
}

func v1TimeToRFC3339(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	switch v := d.Value.(type) {
	case v1.Time:
		return v.ToUnstructured(), nil
	case *v1.Time:
		if v == nil {
			return nil, nil
		}
		return v.ToUnstructured(), nil
	default:
		return nil, fmt.Errorf("Invalid time format %T!\n", v)
	}
}

func isNotFoundError(err error) bool {
	return strings.HasSuffix(err.Error(), "not found")
}

func mergeTags(labels map[string]string, annotations map[string]string) map[string]string {
	tags := make(map[string]string)
	for k, v := range annotations {
		tags[k] = v
	}
	for k, v := range labels {
		tags[k] = v
	}
	return tags
}

//// HYDRATE FUNCTIONS

// There's probably a better way to do this than opening and parsing the config file again...
func getKubectlContext(ctx context.Context, _ *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getKubectlContext")

	// Instantiate loader for kubeconfig file.
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	rawConfig, _ := kubeconfig.RawConfig()
	return rawConfig.CurrentContext, nil

}

//// TRANSFORM FUNCTIONS

func ensureStringArray(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		// Should we return empty list instead???
		// return []string{}, nil

		return nil, fmt.Errorf("ensureStringArray - Cannot transform nil value")
	}

	switch v := d.Value.(type) {
	case []string:
		return v, nil
	case string:
		return []string{v}, nil
	default:
		str := fmt.Sprintf("%v", d.Value)
		return []string{string(str)}, nil
	}
}
