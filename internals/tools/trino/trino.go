package trinocl

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/rs/zerolog"
	"github.com/trinodb/trino-go-client/trino"
	_ "github.com/trinodb/trino-go-client/trino"
	appconfigs "github.com/vapusdata-ecosystem/vapusdata/core/app/configs"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	dmlogger "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/logger"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type TrinoClient struct {
	Address string
	Authn   string //https://trino.io/docs/current/security/jwt.html
	logger  zerolog.Logger

	Catalog          string
	CatalogMountName string
	CatalogMountPath string

	CatalogSecrets          string
	CatalogSecretsMountPath string
	CatalogSecretsMountName string

	Spec           *appconfigs.TrinoDeploymentSpec
	Client         *sqlx.DB
	IsReader       bool `default:"true"`
	NeedRestartMap map[string]bool
}

type TrinoCatalog struct {
	Catalogname string
	Content     string
}

type TrinoCatalogOpsParams struct {
	Catalogs      []TrinoCatalog
	K8SClientset  *kubernetes.Clientset
	Secretname    string
	SecretContent string
}

type TrinoOpts func(*TrinoClient)

func WithAddress(address string) TrinoOpts {
	return func(tc *TrinoClient) {
		tc.Address = address
	}
}

func WithCatalog(cf string) TrinoOpts {
	return func(tc *TrinoClient) {
		tc.Catalog = cf
	}
}

func WithCatalogMountPath(cf string) TrinoOpts {
	return func(tc *TrinoClient) {
		tc.CatalogMountPath = cf
	}
}

func WithCatalogMountName(cf string) TrinoOpts {
	return func(tc *TrinoClient) {
		tc.CatalogMountName = cf
	}
}

func WithCatalogSecretsMountPath(cf string) TrinoOpts {
	return func(tc *TrinoClient) {
		tc.CatalogSecretsMountPath = cf
	}
}

func WithCatalogSecretsMountName(cf string) TrinoOpts {
	return func(tc *TrinoClient) {
		tc.CatalogSecretsMountName = cf
	}
}

func WithCatalogSecrets(cf string) TrinoOpts {
	return func(tc *TrinoClient) {
		tc.CatalogSecrets = cf
	}
}

func WithDeploymentSpec(cf *appconfigs.TrinoDeploymentSpec) TrinoOpts {
	return func(tc *TrinoClient) {
		tc.Spec = cf
	}
}

func WithReaderSpec(cf bool) TrinoOpts {
	return func(tc *TrinoClient) {
		tc.IsReader = cf
	}
}

func buildAddress(params *TrinoClient) string {
	return fmt.Sprintf("http://%s@%s:%d", DEFAULT_USERNAME, params.Spec.TrinoCordSvc, params.Spec.TrinoCordSvcPort)
}

func New(logger zerolog.Logger, opts ...TrinoOpts) *TrinoClient {
	tc := &TrinoClient{}
	for _, opt := range opts {
		opt(tc)
	}
	tc.logger = dmlogger.GetSubDMLogger(logger, "trino", "vapusdata")
	// client := trino.NewClient(trino.Config{
	// 	ServerURI: "http://trino-server:8080",
	// 	User:      "user",
	// 	Password:  "password",
	// })
	tc.Address = buildAddress(tc)
	logger.Info().Msg("Trino address: " + tc.Address)
	foobarClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       &tls.Config{
				// your config here...
			},
		},
	}
	trino.RegisterCustomClient("foobar", foobarClient)
	log.Println("Trino address: ", tc.Address)
	db, err := sqlx.Open("trino", tc.Address+"?custom_client=foobar")
	if err != nil {
		tc.logger.Fatal().Err(err).Msg("Failed to connect to Trino")
	}
	tc.Client = db
	tc.NeedRestartMap = make(map[string]bool)
	return tc
}

func (x *TrinoClient) QueryWithMapResult(query string, logger zerolog.Logger) ([]map[string]any, error) {
	results := make([]map[string]any, 0)
	query = strings.TrimSuffix(query, ";")
	x.logger.Info().Msg("QueryWithMapResult - " + query)
	rows, err := x.Client.Queryx(query)
	if err != nil {
		x.logger.Error().Err(err).Msg("Failed to query Trino")
		return results, err
	}
	log.Println("QueryWithMapResult - ", query)
	defer rows.Close()
	for rows.Next() {
		// Scan each row into a map
		row := make(map[string]any)
		err := rows.MapScan(row)
		if err != nil {
			x.logger.Error().Err(err).Msg("Failed to map scan Trino result")
			continue
		}
		results = append(results, row)
	}
	return results, nil
}

func (x *TrinoClient) SetCatalogName(name string) string {
	return name + ".properties"
}

func (x *TrinoClient) setDSSecretFile(opts *TrinoCatalogOpsParams) string {
	opts.Secretname = opts.Secretname + ".properties"
	return x.CatalogSecretsMountPath + "/" + opts.Secretname
}

func (x *TrinoClient) renderCatalogContent(content string, secFilePath string) string {
	return strings.Replace(content, "{secretsFileFullPath}", secFilePath, -1)
}

func (x *TrinoClient) AddDataSourceCatalog(ctx context.Context, agentId string, opts *TrinoCatalogOpsParams, dbParams *models.DataSourceCredsParams) error {
	if x.IsReader {
		return ErrTrinoIsReadOnly
	}
	x.NeedRestartMap[agentId] = true
	defer delete(x.NeedRestartMap, agentId)
	var err error
	x.logger = dmlogger.GetSubDMLogger(x.logger, "trino-AddDataSourceCatalog", agentId)
	secpath := x.setDSSecretFile(opts)
	// Create Trino catalog content
	for _, catalog := range opts.Catalogs {
		catalogContent := x.renderCatalogContent(catalog.Content, secpath)
		catalog.Catalogname = x.SetCatalogName(catalog.Catalogname)
		err := x.UpsertTrinoCatalog(ctx, opts, catalog.Catalogname, catalogContent, agentId)
		if err != nil {
			x.logger.Error().Err(err).Msg("Failed to UpsertTrinoCatalog in AddDataSourceCatalog")
			delete(x.NeedRestartMap, agentId)
			return err
		}
	}
	err = x.UpsertTrinoCatalogSecretFile(ctx, opts, agentId)
	if err != nil {
		delete(x.NeedRestartMap, agentId)
		x.logger.Error().Err(err).Msg("Failed to UpsertTrinoCatalogSecretFile in AddDataSourceCatalog")
		return err
	}
	// if x.NeedRestartMap[agentId] {
	// 	delete(x.NeedRestartMap, agentId)
	// 	return x.RolloutTrinoDeployment(ctx, opts)
	// } else {
	// 	delete(x.NeedRestartMap, agentId)
	// 	return nil
	// }
	return x.RolloutTrinoDeployment(ctx, opts)
}

func (x *TrinoClient) UpsertTrinoCatalog(ctx context.Context, opts *TrinoCatalogOpsParams, catalogName, catalogContent, agentId string) error {
	if x.IsReader {
		return ErrTrinoIsReadOnly
	}
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		configMap, err := opts.K8SClientset.CoreV1().ConfigMaps(x.Spec.Namespace).Get(ctx, x.Catalog, metav1.GetOptions{})
		if err != nil {
			// Create ConfigMap if it doesn't exist
			newConfigMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      x.Catalog,
					Namespace: x.Spec.Namespace,
				},
				Data: map[string]string{
					catalogName: catalogContent,
				},
			}
			_, err := opts.K8SClientset.CoreV1().ConfigMaps(x.Spec.Namespace).Create(ctx, newConfigMap, metav1.CreateOptions{})
			if err != nil {
				x.logger.Err(err).Msgf("Failed to create trino catalog: %v", err)
				return err
			}
			return err
		}
		currMD, ok := configMap.Data[catalogName]
		if !ok {
			// Update existing ConfigMap
			configMap.Data[catalogName] = catalogContent
			_, err = opts.K8SClientset.CoreV1().ConfigMaps(x.Spec.Namespace).Update(ctx, configMap, metav1.UpdateOptions{})
			if err != nil {
				x.logger.Err(err).Msgf("Failed to update trino catalog: %v", err)
				return err
			}
			x.logger.Info().Msg("Updated trino catalog")
			return nil
		} else {
			if currMD != catalogContent {
				configMap.Data[catalogName] = catalogContent
				_, err = opts.K8SClientset.CoreV1().ConfigMaps(x.Spec.Namespace).Update(ctx, configMap, metav1.UpdateOptions{})
				if err != nil {
					x.logger.Err(err).Msgf("Failed to update trino catalog: %v", err)
					return err
				}
				x.logger.Info().Msg("Updated trino catalog existing item")
				return nil
			} else {
				x.NeedRestartMap[agentId] = false
				return nil
			}
		}
	})
}

func (x *TrinoClient) UpsertTrinoCatalogSecretFile(ctx context.Context, opts *TrinoCatalogOpsParams, agentId string) error {
	if x.IsReader {
		return ErrTrinoIsReadOnly
	}
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		trinoSecrets, err := opts.K8SClientset.CoreV1().Secrets(x.Spec.Namespace).Get(ctx, x.CatalogSecrets, metav1.GetOptions{})
		if err != nil {
			// Create ConfigMap if it doesn't exist
			x.logger.Info().Msg("Creating trino secrets")
			trinoSecrets := &corev1.Secret{
				Immutable: dmutils.Bool2Ptr(false),
				Type:      corev1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name: x.CatalogSecrets,
				},
				Data: map[string][]byte{
					opts.Secretname: []byte(opts.SecretContent),
				},
			}
			_, err := opts.K8SClientset.CoreV1().Secrets(x.Spec.Namespace).Create(ctx, trinoSecrets, metav1.CreateOptions{})
			if err != nil {
				x.logger.Err(err).Msgf("Failed to create trino secrets: %v", err)
				return err
			}
			x.logger.Info().Msg("Updating trino coordinator deployment with secrets")
			deployment, err := opts.K8SClientset.AppsV1().Deployments(x.Spec.Namespace).Get(ctx, x.Spec.TrinoCordDeployment, metav1.GetOptions{})
			if err != nil {
				x.logger.Error().Err(err).Msg("Failed to get trino deployment")
				return err
			}
			deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: x.CatalogSecretsMountName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: x.CatalogSecrets,
					},
				},
			})
			for i, container := range deployment.Spec.Template.Spec.Containers {
				if container.Name == x.Spec.TrinoCordDeploymentContainer {
					deployment.Spec.Template.Spec.Containers[i].VolumeMounts = append(deployment.Spec.Template.Spec.Containers[i].VolumeMounts, corev1.VolumeMount{
						Name:      x.CatalogSecretsMountName,
						MountPath: x.CatalogSecretsMountPath,
					})
				}
			}
			x.logger.Info().Msg("Updating trino worker deployment with secrets")
			workerDeployment, err := opts.K8SClientset.AppsV1().Deployments(x.Spec.Namespace).Get(ctx, x.Spec.TrinoWorkerDeployment, metav1.GetOptions{})
			if err != nil {
				x.logger.Error().Err(err).Msg("Failed to get trino deployment")
				return err
			}
			workerDeployment.Spec.Template.Spec.Volumes = append(workerDeployment.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: x.CatalogSecretsMountName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: x.CatalogSecrets,
					},
				},
			})
			for i, container := range workerDeployment.Spec.Template.Spec.Containers {
				if container.Name == x.Spec.TrinoWorkerDeploymentContainer {
					workerDeployment.Spec.Template.Spec.Containers[i].VolumeMounts = append(workerDeployment.Spec.Template.Spec.Containers[i].VolumeMounts, corev1.VolumeMount{
						Name:      x.CatalogSecretsMountName,
						MountPath: x.CatalogSecretsMountPath,
					})
				}
			}
			// Apply the update
			_, err = opts.K8SClientset.AppsV1().Deployments(x.Spec.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
			if err != nil {
				x.logger.Error().Err(err).Msg("Failed to update trino deployment in UpsertTrinoCatalogSecretFile")
				return nil
			}
			x.logger.Info().Msg("Updated trino coordinator deployment with secrets")
			_, err = opts.K8SClientset.AppsV1().Deployments(x.Spec.Namespace).Update(ctx, workerDeployment, metav1.UpdateOptions{})
			if err != nil {
				x.logger.Error().Err(err).Msg("Failed to update trino workerDeployment in UpsertTrinoCatalogSecretFile")
				return err
			}
			x.logger.Info().Msg("Updated trino worker deployment with secrets")
			return nil
		}
		exSecret, ok := trinoSecrets.Data[opts.Secretname]
		if !ok {
			trinoSecrets.Data[opts.Secretname] = []byte(opts.SecretContent)
			_, err = opts.K8SClientset.CoreV1().Secrets(x.Spec.Namespace).Update(ctx, trinoSecrets, metav1.UpdateOptions{})
			return err
		} else {
			if string(exSecret) != opts.SecretContent {
				trinoSecrets.Data[opts.Secretname] = []byte(opts.SecretContent)
				_, err = opts.K8SClientset.CoreV1().Secrets(x.Spec.Namespace).Update(ctx, trinoSecrets, metav1.UpdateOptions{})
				if err != nil {
					x.logger.Error().Err(err).Msg("Failed to update trino secrets")
					return err
				}
				x.logger.Info().Msg("Updated trino secrets")
				return err
			} else {
				x.NeedRestartMap[agentId] = false
				return nil
			}
		}
	})
}

func (x *TrinoClient) RolloutTrinoDeployment(ctx context.Context, opts *TrinoCatalogOpsParams) error {
	if x.IsReader {
		return ErrTrinoIsReadOnly
	}
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		deployment, err := opts.K8SClientset.AppsV1().Deployments(x.Spec.Namespace).Get(ctx, x.Spec.TrinoCordDeployment, metav1.GetOptions{})
		if err != nil {
			x.logger.Error().Err(err).Msg("Failed to get trino deployment")
			return err
		}
		if deployment.Spec.Template.Annotations == nil {
			deployment.Spec.Template.Annotations = make(map[string]string)
		}
		deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

		// Apply the update
		_, err = opts.K8SClientset.AppsV1().Deployments(x.Spec.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
		if err != nil {
			x.logger.Error().Err(err).Msg("Failed to update trino deployment in RolloutTrinoDeployment")
			return err
		}
		return nil
	})
	if err != nil {
		x.logger.Error().Err(err).Msg("Failed to update trino deployment in RolloutTrinoDeployment")
	}
	time.Sleep(4 * time.Second)
	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		workerDeployment, err := opts.K8SClientset.AppsV1().Deployments(x.Spec.Namespace).Get(ctx, x.Spec.TrinoWorkerDeployment, metav1.GetOptions{})
		if err != nil {
			x.logger.Error().Err(err).Msg("Failed to get trino workerDeployment")
			return err
		}
		if workerDeployment.Spec.Template.Annotations == nil {
			workerDeployment.Spec.Template.Annotations = make(map[string]string)
		}
		workerDeployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

		// Apply the update
		_, err = opts.K8SClientset.AppsV1().Deployments(x.Spec.Namespace).Update(ctx, workerDeployment, metav1.UpdateOptions{})
		if err != nil {
			x.logger.Error().Err(err).Msg("Failed to update trino workerDeployment in RolloutTrinoDeployment")
			return err
		}
		return nil
	})
	if err != nil {
		x.logger.Error().Err(err).Msg("Failed to update trino workerDeployment in RolloutTrinoDeployment")
	}

	// Wait for pods to restart
	return nil
}
