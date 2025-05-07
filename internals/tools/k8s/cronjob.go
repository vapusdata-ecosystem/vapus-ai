package k8s

import "k8s.io/client-go/kubernetes"

type K8SCronjob struct {
	Namespace string
	ConfigMap string
	ClientSet *kubernetes.Clientset
	Immutable *bool
}
