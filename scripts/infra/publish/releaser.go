package helms

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

type HelmReleaser struct {
	HelmDir, HelmRegistry, ChartName, FullChartPath, OciPath string
	Chart                                                    *chart.Chart
	LogoutRegistry                                           bool
}

func NewHelmReleaser(helmDir, helmRegistry, chartName string) (*HelmReleaser, error) {
	a := &HelmReleaser{
		HelmDir:       helmDir,
		HelmRegistry:  helmRegistry,
		ChartName:     chartName,
		FullChartPath: filepath.Join(helmDir, chartName),
		OciPath:       "oci://" + helmRegistry,
	}

	chart, err := loader.LoadDir(filepath.Join(a.FullChartPath))
	if err != nil {
		logger.Info().Msgf("Failed to load chart: %v", err)
		return nil, err
	}
	a.Chart = chart
	log.Println("HELM_CHART_VERSION:", a.Chart.Metadata.Version)
	err = os.WriteFile(filepath.Join("no-push-helm-chart-version.txt"), []byte(a.Chart.Metadata.Version), 0644)
	if err != nil {
		logger.Info().Msgf("Failed to write chart version to file: %v", err)
		return nil, err
	}

	logger.Info().Msgf("Loaded chart: %v", a.Chart.Metadata.Name)
	logger.Info().Msg("Helmer releaser created")
	return a, nil
}

func (h *HelmReleaser) UploadHelmOciChart() error {
	var err error
	tempDir, err := os.MkdirTemp("", "helm-chart-")
	defer os.RemoveAll(tempDir)
	err = exec.Command("helm", "dependency", "update", h.FullChartPath).Run()
	if err != nil {
		logger.Err(err).Msgf("Failed to update the dependencies: %v", err)
		return err
	}
	logger.Info().Msg("Updated the dependencies")
	chartPackagePath, err := chartutil.Save(h.Chart, tempDir)
	if err != nil {
		logger.Err(err).Msgf("Failed to save chart: %v", err)
		return err
	}
	logger.Info().Msgf("Saved chart: %v", chartPackagePath)

	// Initialize the registry client
	// registryClient, err := registry.NewClient(
	// 	registry.ClientOptDebug(true),
	// 	registry.ClientOptWriter(os.Stdout),
	// )
	if err != nil {
		logger.Err(err).Msgf("Failed to create registry client: %v", err)
		return err
	}
	// Login to the registry
	// err = exec.Command("helm", "registry", "login", helmRegistry, "-u", helmUsername, "-p", helmPassword).Run()
	// if err != nil {
	// 	logger.Err(err).Msgf("Failed to login to registry using CLI: %v", err)
	// 	return err
	// }
	logger.Info().Msg("Logged in to registry using CLI")
	if h.LogoutRegistry {
		defer func() {
			err = exec.Command("docker", "logout", helmRegistry).Run()
			if err != nil {
				logger.Info().Msgf("Failed to logout from registry using CLI: %v", err)
			}
			logger.Info().Msg("Logged out from registry using CLI")
		}()
	}
	// Lint the chart
	// var lintoutput []byte
	// logger.Info().Msg("Linting the chart")
	// lintoutput, err = exec.Command("helm", "lint", chartPackagePath).CombinedOutput()
	// if err != nil {
	// 	logger.Info().Msgf("Helmer lint output: %v", string(lintoutput))
	// 	logger.Fatal().Msgf("Helmer lint failed: %v", err)
	// }

	// registryClient.Login(helmRegistry, registry.LoginOptBasicAuth(helmUsername, helmPassword))

	var output []byte
	output, err = exec.Command("helm", "push", chartPackagePath, h.OciPath).CombinedOutput()
	if err != nil {
		logger.Err(err).Err(err).Msgf("Failed to push chart to OCI registry using CLI: %v", err)
		return err

	}
	logger.Info().Msgf("Pushed chart to OCI registry using CLI: %v", string(output))
	digest := getDigestFromCmOp(string(output))
	logger.Info().Msgf("Pushed chart to OCI registry with digest: %v", digest)
	err = os.WriteFile(filepath.Join("no-push-helm-chart-version.txt"), []byte(digest), 0644)
	if err != nil {
		logger.Info().Msgf("Failed to write chart version to file: %v", err)
		return err
	}
	return nil
}
func (h *HelmReleaser) BumpVersion() error {
	h.Chart.Metadata.Version = bumpChartVersion(h.Chart.Metadata.Version)
	err := os.WriteFile(filepath.Join("no-push-helm-chart-version.txt"), []byte(h.Chart.Metadata.Version), 0644)
	if err != nil {
		logger.Info().Msgf("Failed to write chart version to file: %v", err)
		return err
	}
	h.Chart.Metadata.AppVersion = appVersion
	err = chartutil.SaveChartfile(filepath.Join(h.FullChartPath, "Chart.yaml"), h.Chart.Metadata)
	if err != nil {
		logger.Err(err).Msgf("Failed to save chart file: %v", err)
		return err
	}
	logger.Info().Msgf("Bumped chart version to: %v", h.Chart.Metadata.Version)
	return nil
}
func (h *HelmReleaser) UpdateValues() error {
	err := updateVapusDataValues(h.Chart, filepath.Join(h.FullChartPath, "values.yaml"))
	if err != nil {
		logger.Err(err).Msgf("Failed to update values.yaml: %v", err)
		return err
	}
	logger.Info().Msg("Updated values.yaml")
	return nil
}
