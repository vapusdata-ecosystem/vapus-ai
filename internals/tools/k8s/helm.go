package k8s

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/client-go/kubernetes"
)

type HelmInstallParams struct {
	Namespace   string
	ChartName   string
	ReleaseName string
	Version     string
	Values      map[string]any
	ChartRepo   string
	Clientset   *kubernetes.Clientset
	Logger      zerolog.Logger
	Provider    string
}

func InstallChart(params *HelmInstallParams) error {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(cli.New().RESTClientGetter(), params.Namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		params.Logger.Err(err).Msgf("Error while initializing helm action config -- %v", err)
		return dmerrors.DMError(err, nil)
	}
	repoFile := filepath.Join(cli.New().RepositoryCache, "repositories.yaml")
	// repoEntry := &repo.Entry{
	// 	Name: params.Provider,
	// 	URL:  params.ChartRepo,
	// }
	repoFileExists, err := repo.LoadFile(repoFile)
	if err != nil || !repoFileExists.Has(params.Provider) {
		params.Logger.Err(err).Msgf("Error while loading repo file -- %v", err)
	}

	// Pull and install Vault Helm chart
	return installer(params, actionConfig)
}

func installer(params *HelmInstallParams, actionConfig *action.Configuration) error {
	var err error
	install := action.NewInstall(actionConfig)
	install.Namespace = params.Namespace
	install.ReleaseName = params.ReleaseName
	install.Version = params.Version
	install.Wait = true
	install.Timeout = 300
	install.CreateNamespace = true
	chartPath, err := install.LocateChart(params.ChartName, cli.New())
	if err != nil {
		log.Fatalf("Error locating chart: %v", err)
	}
	chart, err := loader.Load(chartPath)
	if err != nil {
		log.Fatalf("Error loading chart: %v", err)
	}
	_, err = install.Run(chart, map[string]any{
		"server": map[string]any{
			"standalone": map[string]any{
				"enabled": true,
			},
			"dev": false,
		},
	})
	if err != nil {
		params.Logger.Err(err).Msgf("Error while installing helm chart -- %v", err)
		return dmerrors.DMError(err, nil)
	}
	return nil
}

// func getVaultKeys(clientset *kubernetes.Clientset) ([]string, string, error) {
// 	// Wait for the Vault secret to be created
// 	var vaultSecret *v1.Secret
// 	err := wait.PollImmediate(2*time.Second, 2*time.Minute, func() (bool, error) {
// 		secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), fmt.Sprintf("%s-keys", releaseName), metav1.GetOptions{})
// 		if err == nil {
// 			vaultSecret = secret
// 			return true, nil
// 		}
// 		return false, nil
// 	})
// 	if err != nil {
// 		return nil, "", fmt.Errorf("error waiting for Vault keys secret: %w", err)
// 	}

// 	// Retrieve unseal keys and root token from the secret
// 	unsealKeys := make([]string, 0)
// 	for _, key := range vaultSecret.Data["vault-root"] {
// 		unsealKey, err := base64.StdEncoding.DecodeString(key)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("error decoding unseal key: %w", err)
// 		}
// 		unsealKeys = append(unsealKeys, string(unsealKey))
// 	}

// 	rootToken := string(vaultSecret.Data["vault-root"])

// 	return unsealKeys, rootToken, nil
// }
