package kubernetes

import (
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// convertUnstructuredDataToType converts the file content into a concrete resource type based on the kind of the file content.
func convertUnstructuredDataToType(obj *unstructured.Unstructured) (any, error) {
	switch obj.GetKind() {
	case "ClusterRoleBinding":
		targetObj := &rbacv1.ClusterRoleBinding{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "CustomResourceDefinition":
		targetObj := &apiextensionsv1.CustomResourceDefinition{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "ClusterRole":
		targetObj := &rbacv1.ClusterRoleBinding{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "ConfigMap":
		targetObj := &corev1.ConfigMap{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "DaemonSet":
		targetObj := &appsv1.DaemonSet{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "Deployment":
		targetObj := &appsv1.Deployment{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "EndpointSlice":
		targetObj := &discoveryv1.EndpointSlice{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "Endpoints":
		targetObj := &corev1.Endpoints{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "Event":
		targetObj := &corev1.Event{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "HorizontalPodAutoscaler":
		targetObj := &autoscalingv1.HorizontalPodAutoscaler{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "Ingress":
		targetObj := &networkingv1.Ingress{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "Job":
		targetObj := &batchv1.Job{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "LimitRange":
		targetObj := &corev1.LimitRange{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "Namespace":
		targetObj := &corev1.Namespace{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "NetworkPolicy":
		targetObj := &networkingv1.NetworkPolicy{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "Node":
		targetObj := &corev1.Node{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "PersistentVolumeClaim":
		targetObj := &corev1.PersistentVolumeClaim{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "PersistentVolume":
		targetObj := &corev1.PersistentVolume{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "PodDisruptionBudget":
		targetObj := &policyv1.PodDisruptionBudget{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "PodSecurityPolicy":
		targetObj := &policyv1beta1.PodSecurityPolicy{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "Pod":
		targetObj := &corev1.Pod{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "ReplicaSet":
		targetObj := &appsv1.ReplicaSet{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "ReplicationController":
		targetObj := &corev1.ReplicationController{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "ResourceQuota":
		targetObj := &corev1.ResourceQuota{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "RoleBinding":
		targetObj := &rbacv1.RoleBinding{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "Role":
		targetObj := &rbacv1.Role{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "Secret":
		targetObj := &corev1.Secret{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "ServiceAccount":
		targetObj := &corev1.ServiceAccount{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "Service":
		targetObj := &corev1.Service{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "StatefulSet":
		targetObj := &appsv1.StatefulSet{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	case "StorageClass":
		targetObj := &storagev1.StorageClass{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	default:
		targetObj := &unstructured.Unstructured{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &targetObj)
		if err != nil {
			return nil, err
		}
		return targetObj, nil
	}
}
