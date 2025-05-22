package cmd

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"sync"

	cobra "github.com/spf13/cobra"
	plclient "github.com/vapusdata-ecosystem/vapusai/cli/internals/platform"
	setup "github.com/vapusdata-ecosystem/vapusai/cli/internals/setup-config"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/databases"
	secretstore "github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/secrets-stores"
	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
)

var secretsFile, valuesFile, tlsCert, tlsKey string

func NewInstallerSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pkg.SetupCmd,
		Version: version,
		Short:   "This command will setup the config file of vapusoperator that holds the configuration of the vapusoperator.",
		Long:    `This command will setup the config file of vapusoperator that holds the configuration of the vapusoperator.`,
		Run: func(cmd *cobra.Command, args []string) {
			plclient.MasterGlobals.Logger.Info().Msg("Setting up the config file for installation.")
			addStores()
		},
	}
	cmd.PersistentFlags().StringVar(&secretsFile, "secrets", "", "Secrets file containing the secret values of the vapusdata platform")
	cmd.PersistentFlags().StringVar(&valuesFile, "values", "", "Values file containing the configuration of the vapusdata platform")
	cmd.PersistentFlags().StringVar(&tlsCert, "tlscert", "", "TLS certificate file")
	cmd.PersistentFlags().StringVar(&tlsKey, "tlskey", "", "TLS key file")
	return cmd
}
func addStores() {
	ctx := context.Background()
	var err error
	var installerValueBytes []byte
	plclient.MasterGlobals.Logger.Info().Msg("Reading the config file for installation")
	secretBytes, err := os.ReadFile(secretsFile)
	if err != nil {
		cobra.CheckErr(err)
	}
	result := &setup.VapusSecretInstallerConfig{}
	secretResult := &setup.VapusSecretsMap{}
	err = filetools.GenericUnMarshaler(secretBytes, result, filetools.GetConfFileType(secretsFile))
	if err != nil {
		cobra.CheckErr(err)
	}
	log.Println(result.SecretStore.DataSourceEngine, result.SecretStore.DataSourceSvcProvider, "========================")
	log.Println(result.BackendDataStore.DataSourceEngine, result.BackendDataStore.DataSourceSvcProvider, "========================")
	installerValue := &setup.VapusInstallerConfig{}
	if valuesFile != "" {
		installerValueBytes, err = os.ReadFile(valuesFile)
		if err != nil {
			cobra.CheckErr(err)
		}

		err = filetools.GenericUnMarshaler(installerValueBytes, installerValue, filetools.GetConfFileType(secretsFile))
		if err != nil {
			cobra.CheckErr(err)
		}
	} else {
		valuesFile = "vapus-installer.yaml"
	}
	plclient.MasterGlobals.Logger.Info().Msg("Secrets file loaded successfully")
	log.Println(result.SecretStore.DataSourceEngine, result.SecretStore.DataSourceSvcProvider, "========================")
	secretClient, err := secretstore.New(ctx, secretstore.WithDataSourceCredsParams(result.SecretStore), secretstore.WithLogger(plclient.MasterGlobals.Logger))
	if err != nil {
		plclient.MasterGlobals.Logger.Error().Msgf("Error in creating secret store client: %v", err)
		cobra.CheckErr(err)
	}
	plclient.MasterGlobals.Logger.Info().Msg("Secret store client created")
	log.Println(result.CreateDatabase, "========================")
	if result.CreateDatabase {
		plclient.MasterGlobals.Logger.Info().Msgf("Creating database %s", result.BackendDataStore.DataSourceCreds.DB)
		vapusDB := result.BackendDataStore.DataSourceCreds.DB
		result.BackendDataStore.DataSourceCreds.DB = "postgres"
		dbcl, err := databases.New(ctx, databases.WithInApp(true), databases.WithDataSourceCredsParams(result.BackendDataStore), databases.WithLogger(plclient.MasterGlobals.Logger))
		if err != nil {
			cobra.CheckErr(err)
		}
		defer dbcl.Close()
		result.BackendDataStore.DataSourceCreds.DB = vapusDB
		query := "CREATE DATABASE " + result.BackendDataStore.DataSourceCreds.DB
		err = dbcl.RunDDLs(ctx, &query)
		if err != nil {
			plclient.MasterGlobals.Logger.Error().Msgf("Error in creating database: %v", err)
		}
	} else {
		plclient.MasterGlobals.Logger.Info().Msgf("Database %s already exists", result.BackendDataStore.DataSourceCreds.DB)
	}
	plclient.MasterGlobals.Logger.Info().Msgf("Database %s created successfully", result.BackendDataStore.DataSourceCreds.DB)
	plclient.MasterGlobals.Logger.Info().Msgf("Backend data store client created, creating database")
	plclient.MasterGlobals.Logger.Info().Msgf("Secret store client created")
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		secretResult.AuthnSecrets.Secret = plclient.GetSecretName("")
		err = secretClient.WriteSecret(ctx, result.AuthnSecrets, secretResult.AuthnSecrets.Secret)
		if err != nil {
			cobra.CheckErr(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		secretResult.BackendDataStore.Secret = plclient.GetSecretName("")
		err = secretClient.WriteSecret(ctx, result.BackendDataStore, secretResult.BackendDataStore.Secret)
		if err != nil {
			cobra.CheckErr(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		secretResult.FileStore.Secret = plclient.GetSecretName("")
		err = secretClient.WriteSecret(ctx, result.FileStore, secretResult.FileStore.Secret)
		if err != nil {
			cobra.CheckErr(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		secretResult.BackendCacheStore.Secret = plclient.GetSecretName("")
		err = secretClient.WriteSecret(ctx, result.BackendCacheStore, secretResult.BackendCacheStore.Secret)
		if err != nil {
			cobra.CheckErr(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		secretResult.JWTAuthnSecrets.Secret = plclient.GetSecretName("")
		err = secretClient.WriteSecret(ctx, result.JWTAuthnSecrets, secretResult.JWTAuthnSecrets.Secret)
		if err != nil {
			cobra.CheckErr(err)
		}
	}()
	wg.Wait()

	plclient.MasterGlobals.Logger.Info().Msgf("Secrets added successfully, mapped in the config file")

	tlsKeyBytes, err := os.ReadFile(tlsKey)
	if err != nil {
		cobra.CheckErr(err)
	}
	tlsCertBytes, err := os.ReadFile(tlsCert)
	if err != nil {
		cobra.CheckErr(err)
	}
	decodetlsKey := base64.StdEncoding.EncodeToString(tlsKeyBytes)
	if err != nil {
		cobra.CheckErr(err)
	}
	decodetlsCert := base64.StdEncoding.EncodeToString(tlsCertBytes)
	if err != nil {
		cobra.CheckErr(err)
	}
	installerValue.TLSCert.Key = decodetlsKey
	installerValue.TLSCert.Cert = decodetlsCert

	installerValue.SecretStore = result.SecretStore
	installerValue.Secrets = secretResult
	installerValueBytes, err = filetools.GenericMarshaler(installerValue, filetools.GetConfFileType(valuesFile))
	if err != nil {
		cobra.CheckErr(err)
	}
	err = os.WriteFile(valuesFile, installerValueBytes, 0644)
	if err != nil {
		cobra.CheckErr(err)
	}
	// Write the file only and gracefully handles if file already exists
}
