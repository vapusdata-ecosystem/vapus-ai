package k8s

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"unicode"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	gcptp "github.com/vapusdata-ecosystem/vapusai/core/thirdparty/gcp"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type NewK8SInstance struct {
	clientSet *kubernetes.Clientset
}

func OrganizationK8SConfig(ctx context.Context, k8sInfra *models.K8SInfraParams, logger zerolog.Logger) (*api.Config, error) {
	if k8sInfra == nil {
		return nil, ErrK8SInfraParamsNil
	}
	if k8sInfra.GetKubeConfig() != "" {
		kubeConfig := &api.Config{}
		switch k8sInfra.KubeConfigFormat {
		case mpb.ContentFormats_YAML.String():
			err := yaml.Unmarshal([]byte(k8sInfra.GetKubeConfig()), kubeConfig)
			if err != nil {
				logger.Err(err).Msgf("Error while unmarshalling kubeconfig -- %v", err)
			}
			return kubeConfig, nil
		case mpb.ContentFormats_JSON.String():
			err := json.Unmarshal([]byte(k8sInfra.GetKubeConfig()), kubeConfig)
			if err != nil {
				logger.Err(err).Msgf("Error while unmarshalling kubeconfig -- %v", err)
			}
		}
	}
	switch k8sInfra.InfraService {
	case mpb.InfraService_GKE.String():
		svKey, _ := dmutils.Base64Detectors(k8sInfra.GetCredentials().GcpCreds.ServiceAccountKey)
		gkeK8sConfig, err := gcptp.GetGkeKubeConfig(ctx, &gcptp.GcpConfig{
			ProjectID:         k8sInfra.GetCredentials().GcpCreds.ProjectId,
			Region:            k8sInfra.GetCredentials().GcpCreds.Region,
			Zone:              k8sInfra.GetCredentials().GcpCreds.Zone,
			ServiceAccountKey: []byte(svKey),
		}, k8sInfra.Name, logger)
		if err != nil {
			logger.Err(err).Msgf("Error while getting GKE kubeconfig -- %v", err)
			return nil, err
		}
		return gkeK8sConfig, nil
	default:
		return nil, ErrUnsupportedServiceProvider
	}
}

func GetK8sClusteAPI(logger zerolog.Logger, kubeConfig *api.Config, ORGANIZATION string) (*rest.Config, error) {
	var err error
	var config *rest.Config
	if kubeConfig != nil {
		// Build the Kubernetes client configuration
		kubeconfigPath := filepath.Join(os.TempDir(), ORGANIZATION, "kubeconfig")
		if err := clientcmd.WriteToFile(*kubeConfig, kubeconfigPath); err != nil {
			logger.Err(err).Msg("error writing kubeconfig to file")
		}
		defer dmutils.DeleteFile(kubeconfigPath)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			logger.Err(err).Msg("error connecting to kubernetes from config filepaths")
			return nil, err
		}
		return config, nil
	}
	config, err = rest.InClusterConfig()
	if err != nil {
		logger.Err(err).Msg("error connecting to kubernetes from in cluster config")
		kubeconfig := os.Getenv("KUBECONFIG")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logger.Err(err).Msg("error connecting to kubernetes from config filepaths")
			return nil, err
		} else {
			return config, nil
		}
	}
	return config, nil
}

func Getk8sName(name, resourceName, resourceType string) string {
	if name != "" {
		return K8SNameValidation(name)
		// guuid.New().String()
	}
	log.Println("name is empty")
	if len(dmutils.SlugifyBase(resourceName)) < 15 {
		name = dmutils.SlugifyBase(resourceName)
	} else {
		name = dmutils.SlugifyBase(resourceName)[:15]
	}
	name = name + "-" + dmutils.GetUUID()[0:5]
	log.Println("name is", name)
	return K8SNameValidation(name)
}

func K8SNameValidation(name string) string {
	finalName := ""
	if len(name) > 62 {
		finalName = name[:58]
	} else {
		finalName = name
	}
	if !unicode.IsLetter(rune(finalName[0])) && !unicode.IsDigit(rune(finalName[0])) {
		finalName = "v" + finalName
	}
	if !unicode.IsLetter(rune(finalName[len(finalName)-1])) && !unicode.IsDigit(rune(finalName[len(finalName)-1])) {
		finalName = finalName + "v"
	}
	return finalName
}

func GetK8SClientSet(cf *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(cf)
}

func GetHostK8SClientSet(log zerolog.Logger) (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Err(err).Msg("error connecting to kubernetes from in cluster config")
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func GetK8SClientSetFromConfig(log zerolog.Logger, kubeConfig *api.Config, ORGANIZATION string) (*kubernetes.Clientset, error) {
	cf, err := GetK8sClusteAPI(log, kubeConfig, ORGANIZATION)
	if err != nil {
		log.Err(err).Msg("error while getting InClusterConfig kubeconfig")
		return nil, err
	}
	return kubernetes.NewForConfig(cf)
}

func GetLocalK8SClientSet(kubeconfig string) (*kubernetes.Clientset, *rest.Config, error) {
	if kubeconfig == "" {
		dirname, err := os.UserHomeDir()
		if err != nil {
			return nil, nil, dmerrors.ErrHomeDirNotFound
		}
		kubeconfig = filepath.Join(dirname, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, nil, ErrKubeConfigNotFound
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, ErrK8SCouldNotConnect
	}

	return clientset, config, nil
}
