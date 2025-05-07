package k8s

import (
	"context"

	"github.com/rs/zerolog"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sSecrets struct {
	Data       map[string][]byte
	StringData map[string]string
	Namespace  string
	SecretName string
	ClientSet  *kubernetes.Clientset
	SecretType corev1.SecretType
	Immutable  *bool
}

func (x *K8sSecrets) GetOrCreateScrets(ctx context.Context, createIfNotExists bool, log zerolog.Logger) (string, error) {
	secreName, err := x.ClientSet.CoreV1().Secrets(x.Namespace).Get(ctx, x.SecretName, metav1.GetOptions{})
	if err == nil {
		return secreName.Name, nil
	}
	if createIfNotExists {
		secret := &corev1.Secret{
			Type:      x.SecretType,
			Immutable: x.Immutable,
			ObjectMeta: metav1.ObjectMeta{
				Name: x.SecretName,
			},
		}
		if x.Data != nil {
			secret.Data = x.Data
		} else if x.StringData != nil {
			secret.StringData = x.StringData
		}
		_, err = x.ClientSet.CoreV1().Secrets(x.Namespace).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			log.Err(err).Msgf("Error creating secret %v", x.SecretName)
			return "", err
		}
		return x.SecretName, nil
	}
	return "", err
}

func (x *K8sSecrets) DeleteScrets(ctx context.Context, log zerolog.Logger) error {
	err := x.ClientSet.CoreV1().Secrets(x.Namespace).Delete(ctx, x.SecretName, metav1.DeleteOptions{})
	if err != nil {
		log.Err(err).Msgf("Error getting secret %v", x.SecretName)
		return err
	}
	return nil
}
