package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	plclient "github.com/vapusdata-ecosystem/vapusai/cli/internals/platform"
	setup "github.com/vapusdata-ecosystem/vapusai/cli/internals/setup-config"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	"github.com/vapusdata-ecosystem/vapusai/core/pkgs/authn"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
)

func NewInstallerSpecGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pkg.InstallerSecretSpecGenOps,
		Short: "This command is interface to generate installer spec files for the platform installation.",
		Long:  `This command is interface to generate installer spec files for the platform installation.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := generateInstallerSpecTemplate()
			if err != nil {
				cobra.CheckErr(err)
			}
		},
	}
	return cmd
}

func generateInstallerSpecTemplate() error {
	secM := &models.DataSourceCredsParams{
		DataSourceCreds: &models.DataSourceSecrets{
			GenericCredentialModel: &models.GenericCredentialModel{
				AwsCreds:   &models.AWSCreds{},
				GcpCreds:   &models.GCPCreds{},
				AzureCreds: &models.AzureCreds{},
			},
			DB:   "",
			URL:  "",
			Port: 0,
		},
	}
	specVal, err := filetools.GenericMarshaler(&setup.VapusSecretInstallerConfig{
		SecretStore:       secM,
		BackendDataStore:  secM,
		ArtifactStore:     secM,
		BackendCacheStore: secM,
		AuthnSecrets: &authn.AuthnSecrets{
			OIDCSecrets: &authn.OIDCSecrets{},
		},
		JWTAuthnSecrets: &encryption.JWTAuthn{},
	}, "YAML")
	if err != nil {
		return err
	}
	log.Println("Generated installer spec: ", string(specVal))
	fileName := strings.ToLower("vapusdata-secrets.yaml")
	plclient.MasterGlobals.Logger.Info().Msgf("Sample installer %v spec generated with file name - %v \n", specName, fileName)
	err = os.WriteFile(fileName, specVal, 0644)
	if err != nil {
		return err
	}
	return nil
}
