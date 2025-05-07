package helms

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
)

var (
	logger zerolog.Logger
	helmRegistry, platformSvcDigest, platformSvcTag,
	vapusctlSvcDigest, vapusctlSvcTag,
	dataworkerSvcDigest, dataworkerSvcTag,
	webappDigest, webappTag,
	aistudioDigest, aistudioTag,
	aigatewayTag, aigatewayDigest,
	nabrunnersDigest, nabrunnersTag,
	nabhikserverSvcDigest, nabhikserverSvcTag,
	vapusDcSvcDigest, vapusDcSvcTag, appVersion string
	upload, bumpVersion, updateValues, registryLogout bool
)

func initFlags() {
	flag.StringVar(&helmRegistry, "helm-registry", "", "URL of registry")
	flag.StringVar(&platformSvcDigest, "platform-svc-digest", "", "Platform service digest")
	flag.StringVar(&platformSvcTag, "platform-svc-tag", "", "Platform service tag")
	flag.StringVar(&vapusctlSvcDigest, "vapusctl-svc-digest", "", "vapusctl service digest")
	flag.StringVar(&vapusctlSvcTag, "vapusctl-svc-tag", "", "vapusctl service tag")
	flag.StringVar(&dataworkerSvcDigest, "dataworker-svc-digest", "", "dataworker service digest")
	flag.StringVar(&dataworkerSvcTag, "dataworker-svc-tag", "", "dataworker service tag")
	flag.StringVar(&webappDigest, "webapp-svc-digest", "", "webapp service digest")
	flag.StringVar(&webappTag, "webapp-svc-tag", "", "webapp service tag")
	flag.StringVar(&vapusDcSvcDigest, "vapus-dc-svc-digest", "", "vapus-dc service digest")
	flag.StringVar(&vapusDcSvcTag, "vapus-dc-svc-tag", "", "vapus-dc service tag")
	flag.StringVar(&nabhikserverSvcDigest, "nabhikserver-svc-digest", "", "nabhikserver service digest")
	flag.StringVar(&nabhikserverSvcTag, "nabhikserver-svc-tag", "", "nabhikserver service tag")
	flag.StringVar(&aistudioDigest, "aistudio-svc-digest", "", "aistudio service digest")
	flag.StringVar(&aistudioTag, "aistudio-svc-tag", "", "aistudio service tag")
	flag.StringVar(&nabrunnersDigest, "nabrunners-svc-digest", "", "nabrunners service digest")
	flag.StringVar(&nabrunnersTag, "nabrunners-svc-tag", "", "nabrunners service tag")
	flag.StringVar(&aigatewayDigest, "aigateway-svc-digest", "", "aigateway service digest")
	flag.StringVar(&aigatewayTag, "aigateway-svc-tag", "", "aigateway service tag")
	flag.StringVar(&appVersion, "appVersion", "", "App version of the chart")
	flag.BoolVar(&upload, "upload", true, "Flag to control the upload of helm chart")
	flag.BoolVar(&bumpVersion, "bump-version", true, "Flag to bump the version of helm chart")
	flag.BoolVar(&updateValues, "update-values", false, "Flag to update the values.yaml")
	flag.BoolVar(&registryLogout, "registry-logout", true, "Flag to update the values.yaml")
	flag.Parse()
}

type ArtifactValues struct {
	Tag    string `yaml:"tag"`
	Digest string `yaml:"digest"`
}

func HelmChartOps() string {
	logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	initFlags()
	log.Println(helmRegistry, platformSvcDigest, platformSvcTag, vapusctlSvcDigest, vapusctlSvcTag, dataworkerSvcDigest, dataworkerSvcTag, webappDigest, webappTag, nabhikserverSvcDigest, nabhikserverSvcTag, vapusDcSvcDigest, vapusDcSvcTag, aistudioDigest, aistudioTag, appVersion, upload, bumpVersion, updateValues)
	logger.Info().Msg("Starting helm chart operations with flags")
	helmReleaser, err := NewHelmReleaser("../../deployments", helmRegistry, "vapusdata-platform")
	if err != nil {
		logger.Info().Msgf("Failed to create helm releaser: %v", err)
		return ""
	}
	helmReleaser.LogoutRegistry = registryLogout
	tempDir, err := os.MkdirTemp("", "helm-chart-")
	logger.Info().Msgf("helmRegistry: %v", helmRegistry)
	if err != nil {
		logger.Info().Msgf("Error creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Load the chart from the local directory
	if !upload {
		logger.Info().Msg("Bumping version of helm chart")
		err := helmReleaser.BumpVersion()
		if err != nil {
			return ""
		}
	}

	if updateValues {
		logger.Info().Msg("Updating values.yaml")
		err := helmReleaser.UpdateValues()
		if err != nil {
			return ""
		}
	}
	if upload {
		logger.Info().Msg("Uploading helm chart")
		err = helmReleaser.UploadHelmOciChart()
		if err != nil {
			return ""
		}
	} else {
		logger.Info().Msg("Skipping upload of helm chart")
		return helmReleaser.Chart.Metadata.Version
	}
	return ""
}

func getDigestFromCmOp(text string) string {
	a := strings.Split(text, "Digest:")
	return strings.Trim(strings.Split(a[1], "\n")[0], " ")
}

func SetHelmChartVersion(current, bumptType string) string {
	if current == "" {
		v, err := semver.Make("0.0.1")
		if err != nil {
			log.Println("error while creating semver", err)
			return ""
		}
		return v.String()

	}
	v, _ := semver.Make(current)

	switch strings.ToLower(bumptType) {
	case "major":
		v.Major++
		v.Minor = 0
		v.Patch = 0
		return v.String()
	case "minor":
		v.Minor++
		v.Patch = 0
		return v.String()
	case "patch":
		v.Patch++
		return v.String()
	default:
		return current
	}
}

func bumpChartVersion(current string) string {
	return SetHelmChartVersion(current, "PATCH")
}

func updateSvcArtifacts(tag, digest string, values map[string]any) map[string]any {
	log.Println("Updating service artifacts -------> ", tag, digest)
	svcArtifacts := values["artifacts"].(map[string]any)
	if tag != "" {
		svcArtifacts["tag"] = tag
	}
	if digest != "" {
		svcArtifacts["digest"] = digest
	}
	values["artifacts"] = svcArtifacts
	return values
}

func updateVapusDataValues(chart *chart.Chart, file string) error {
	log.Println("Updating values.yaml", file)
	bytes, err := os.ReadFile(file)
	if err != nil {
		logger.Info().Msgf("Failed to read file: %v", err)
		return err
	}
	result, err := chartutil.ReadValues(bytes)
	if err != nil {
		logger.Info().Msgf("Failed to read values: %v", err)
		return err
	}
	values := result.AsMap()
	log.Println("Values before are -----> ", values)
	if pls, ok := values["platform"].(map[string]any); ok {
		values["platform"] = updateSvcArtifacts(platformSvcTag, platformSvcDigest, pls)
	}
	if pls, ok := values["aistudio"].(map[string]any); ok {
		values["aistudio"] = updateSvcArtifacts(aistudioTag, aistudioDigest, pls)
	}
	if pls, ok := values["webapp"].(map[string]any); ok {
		values["webapp"] = updateSvcArtifacts(webappTag, webappDigest, pls)
	}
	if pls, ok := values["nabhikserver"].(map[string]any); ok {
		values["nabhikserver"] = updateSvcArtifacts(nabhikserverSvcTag, nabhikserverSvcDigest, pls)
	}
	if pls, ok := values["nabrunners"].(map[string]any); ok {
		values["nabrunners"] = updateSvcArtifacts(nabrunnersTag, nabrunnersDigest, pls)
	}
	if pls, ok := values["vapusaigateway"].(map[string]any); ok {
		values["vapusaigateway"] = updateSvcArtifacts(aigatewayTag, aigatewayDigest, pls)
	}
	logger.Info().Msg("Svc Values are updated")
	if vdas, ok := values[VAPUSDATA_ARTIFACTS].(map[string]any); ok {

		logger.Info().Msg("Updating vapusdata artifacts")
		if dataArtifacts, ok := vdas["dataworker"].(map[string]any); ok {
			if dataworkerSvcDigest != "" {
				dataArtifacts["digest"] = dataworkerSvcDigest
			}
			if dataworkerSvcTag != "" {
				dataArtifacts["tag"] = dataworkerSvcTag
			}
			vdas["dataworker"] = dataArtifacts
		}

		if vdcArtifacts, ok := vdas["vdc"].(map[string]any); ok {
			if nabhikserverSvcDigest != "" {
				vdcArtifacts["digest"] = vapusDcSvcDigest
			}
			if nabhikserverSvcTag != "" {
				vdcArtifacts["tag"] = vapusDcSvcTag
			}
			vdas["vdc"] = vdcArtifacts
		}
		values[VAPUSDATA_ARTIFACTS] = vdas
	}

	logger.Info().Msg("Vapusdata artifacts are updated")
	bytes, err = yaml.Marshal(values)
	if err != nil {
		logger.Err(err).Msgf("Failed to marshal values: %v", err)
		return err
	}
	logger.Info().Msg("Values are marshalled")

	err = os.WriteFile(file, bytes, 0644)
	if err != nil {
		logger.Err(err).Msgf("Failed to write values: %v", err)
	}
	logger.Info().Msg("Values are updated in values.yaml")
	return err
}
