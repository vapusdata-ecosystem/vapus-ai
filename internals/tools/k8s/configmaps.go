package k8s

import (
	"context"

	"github.com/rs/zerolog"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sConfigmap struct {
	BinaryData map[string][]byte
	Data       map[string]string
	Namespace  string
	ConfigMap  string
	ClientSet  *kubernetes.Clientset
	Immutable  *bool
}

func (x *K8sConfigmap) GetOrCreateConfigmap(ctx context.Context, createIfNotExists bool, log zerolog.Logger) (string, error) {
	configmapName, err := x.ClientSet.CoreV1().ConfigMaps(x.Namespace).Get(ctx, x.ConfigMap, metav1.GetOptions{})
	if err == nil {
		return configmapName.Name, nil
	}
	if createIfNotExists {
		configmap := &corev1.ConfigMap{
			Immutable: x.Immutable,
			ObjectMeta: metav1.ObjectMeta{
				Name: x.ConfigMap,
			},
		}
		if x.Data != nil {
			configmap.Data = x.Data
		} else if x.BinaryData != nil {
			configmap.BinaryData = x.BinaryData
		}
		_, err = x.ClientSet.CoreV1().ConfigMaps(x.Namespace).Create(ctx, configmap, metav1.CreateOptions{})
		if err != nil {
			log.Err(err).Msgf("Error creating configmap %v", x.ConfigMap)
			return "", err
		}
		return x.ConfigMap, nil
	}
	return "", err
}

func (x *K8sConfigmap) DeleteConfigmap(ctx context.Context, log zerolog.Logger) error {
	err := x.ClientSet.CoreV1().ConfigMaps(x.Namespace).Delete(ctx, x.ConfigMap, metav1.DeleteOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting configmap %v", x.ConfigMap)
		return err
	}
	return nil
}
