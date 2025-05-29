package cmd

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/rs/zerolog"
	cobra "github.com/spf13/cobra"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	plclient "github.com/vapusdata-ecosystem/vapusai/cli/internals/platform"
	setup "github.com/vapusdata-ecosystem/vapusai/cli/internals/setup-config"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/databases"
	secretstore "github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/secrets-stores"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	"github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
	vtls "github.com/vapusdata-ecosystem/vapusai/core/tools/tls"
)

var secretsFile, valuesFile string

func NewInstallerSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pkg.SetupCmd,
		Version: version,
		Short:   "This command will setup the config file of vapusoperator that holds the configuration of the vapusoperator.",
		Long:    `This command will setup the config file of vapusoperator that holds the configuration of the vapusoperator.`,
		Run: func(cmd *cobra.Command, args []string) {
			plclient.MasterGlobals.Logger.Info().Msg("Setting up the config file for installation.")
			configurator := NewSetupConfigurator()
			configurator.configure()
		},
	}
	cmd.PersistentFlags().StringVar(&secretsFile, "secrets", "", "Secrets file containing the secret values of the vapusdata platform")
	cmd.PersistentFlags().StringVar(&valuesFile, "values", "", "Values file containing the configuration of the vapusdata platform")
	return cmd
}

type SetupConfigurator struct {
	inputConfig  *setup.VapusSecretInstallerConfig
	outputConfig *setup.VapusInstallerConfig
	secretsmap   *setup.VapusSecretsMap
	logger       zerolog.Logger
}

func NewSetupConfigurator() *SetupConfigurator {
	plclient.MasterGlobals.Logger.Info().Msg("Reading the config file for installation")
	secretBytes, err := os.ReadFile(secretsFile)
	if err != nil {
		cobra.CheckErr(err)
	}
	inputConfig := &setup.VapusSecretInstallerConfig{}
	err = filetools.GenericUnMarshaler(secretBytes, inputConfig, filetools.GetConfFileType(secretsFile))
	if err != nil {
		cobra.CheckErr(err)
	}
	log.Println(string(secretBytes), "========================")
	if inputConfig.SecretStore == nil {
		inputConfig.SecretStore = &models.DataSourceCredsParams{}
	}
	return &SetupConfigurator{
		inputConfig:  inputConfig,
		outputConfig: &setup.VapusInstallerConfig{},
		secretsmap:   &setup.VapusSecretsMap{},
		logger:       plclient.MasterGlobals.Logger,
	}
}

func (x *SetupConfigurator) configure() {
	ctx := context.Background()
	var err error
	var installerValueBytes []byte
	plclient.MasterGlobals.Logger.Info().Msg("Reading the config file for installation")

	if valuesFile != "" {
		installerValueBytes, err = os.ReadFile(valuesFile)
		if err != nil {
			cobra.CheckErr(err)
		}

		err = filetools.GenericUnMarshaler(installerValueBytes, x.outputConfig, filetools.GetConfFileType(secretsFile))
		if err != nil {
			cobra.CheckErr(err)
		}
	} else {
		valuesFile = "vapus-installer.yaml"
	}
	plclient.MasterGlobals.Logger.Info().Msg("Secrets file loaded successfully")
	log.Println(x.inputConfig.SecretStore.DataSourceEngine, x.inputConfig.SecretStore.DataSourceSvcProvider, "========================")
	secretClient, err := secretstore.New(ctx, secretstore.WithDataSourceCredsParams(x.inputConfig.SecretStore), secretstore.WithLogger(plclient.MasterGlobals.Logger))
	if err != nil {
		plclient.MasterGlobals.Logger.Error().Msgf("Error in creating secret store client: %v", err)
		cobra.CheckErr(err)
	}
	plclient.MasterGlobals.Logger.Info().Msg("Secret store client created")
	log.Println(x.inputConfig.CreateDatabase, "========================")
	if x.inputConfig.CreateDatabase {
		plclient.MasterGlobals.Logger.Info().Msgf("Creating database %s", x.inputConfig.BackendDataStore.DataSourceCreds.DB)
		vapusDB := x.inputConfig.BackendDataStore.DataSourceCreds.DB
		x.inputConfig.BackendDataStore.DataSourceCreds.DB = "postgres"
		dbcl, err := databases.New(ctx, databases.WithInApp(true), databases.WithDataSourceCredsParams(x.inputConfig.BackendDataStore), databases.WithLogger(plclient.MasterGlobals.Logger))
		if err != nil {
			cobra.CheckErr(err)
		}
		defer dbcl.Close()
		x.inputConfig.BackendDataStore.DataSourceCreds.DB = vapusDB
		query := "CREATE DATABASE " + x.inputConfig.BackendDataStore.DataSourceCreds.DB
		err = dbcl.RunDDLs(ctx, &query)
		if err != nil {
			plclient.MasterGlobals.Logger.Error().Msgf("Error in creating database: %v", err)
		}
	} else {
		plclient.MasterGlobals.Logger.Info().Msgf("Database %s already exists", x.inputConfig.BackendDataStore.DataSourceCreds.DB)
	}
	plclient.MasterGlobals.Logger.Info().Msgf("Database %s created successfully", x.inputConfig.BackendDataStore.DataSourceCreds.DB)
	plclient.MasterGlobals.Logger.Info().Msgf("Backend data store client created, creating database")
	plclient.MasterGlobals.Logger.Info().Msgf("Secret store client created")
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		x.secretsmap.AuthnSecrets.Secret = plclient.GetSecretName("")
		err = secretClient.WriteSecret(ctx, x.inputConfig.AuthnSecrets, x.secretsmap.AuthnSecrets.Secret)
		if err != nil {
			cobra.CheckErr(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		x.secretsmap.BackendDataStore.Secret = plclient.GetSecretName("")
		err = secretClient.WriteSecret(ctx, x.inputConfig.BackendDataStore, x.secretsmap.BackendDataStore.Secret)
		if err != nil {
			cobra.CheckErr(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		x.secretsmap.FileStore.Secret = plclient.GetSecretName("")
		err = secretClient.WriteSecret(ctx, x.inputConfig.FileStore, x.secretsmap.FileStore.Secret)
		if err != nil {
			cobra.CheckErr(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		x.secretsmap.BackendCacheStore.Secret = plclient.GetSecretName("")
		err = secretClient.WriteSecret(ctx, x.inputConfig.BackendCacheStore, x.secretsmap.BackendCacheStore.Secret)
		if err != nil {
			cobra.CheckErr(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		x.manageJWTAuthn()
		x.secretsmap.JWTAuthnSecrets.Secret = plclient.GetSecretName("")
		err = secretClient.WriteSecret(ctx, x.inputConfig.JWTAuthnSecrets, x.secretsmap.JWTAuthnSecrets.Secret)
		if err != nil {
			cobra.CheckErr(err)
		}
	}()
	wg.Wait()

	plclient.MasterGlobals.Logger.Info().Msgf("Secrets added successfully, mapped in the config file")
	x.manageTls()
	x.outputConfig.SecretStore = x.inputConfig.SecretStore
	x.outputConfig.Secrets = x.secretsmap
	installerValueBytes, err = filetools.GenericMarshaler(x.outputConfig, filetools.GetConfFileType(valuesFile))
	if err != nil {
		cobra.CheckErr(err)
	}
	err = os.WriteFile(valuesFile, installerValueBytes, 0644)
	if err != nil {
		cobra.CheckErr(err)
	}
	// Write the file only and gracefully handles if file already exists
}

func (x *SetupConfigurator) manageTls() {
	var tlsCertBytes, tlsKeyBytes []byte
	if x.inputConfig.TLSCert.AutoGenerate {
		x.logger.Info().Msg("Generating TLS certificate, selecting algorithm and bitsize")
		algo, err := pkg.EncryptionAlgorithm.Run()
		if err != nil {
			x.logger.Error().Msgf("Error in reading the algorithm: %v", err)
			cobra.CheckErr(err)
		}
		bitSizeS, err := pkg.TlsBitsize.Run()
		if err != nil {
			x.logger.Error().Msgf("Error in reading the bitsize: %v", err)
			cobra.CheckErr(err)
		}
		log.Println(bitSizeS, "========================")
		bitSize, err := strconv.Atoi(bitSizeS)
		if err != nil {
			x.logger.Error().Msgf("Error in converting the bitsize to int: %v, so switching to default", err)
			bitSize = 2048
		}
		log.Println(bitSize, "========================")
		certGenerator, err := vtls.NewTLSOperator(vtls.WithTLSOperatorParams(&vtls.TLSOperatorOpts{
			Algo:   mpb.EncryptionAlgo(mpb.EncryptionAlgo_value[algo]),
			Logger: plclient.MasterGlobals.Logger,
		}))
		log.Println(certGenerator, "========================")
		if err != nil {
			x.logger.Error().Msgf("Error in creating the cert generator: %v", err)
			cobra.CheckErr(err)
		}
		tlsCert, err := certGenerator.GenerateTlsPvtKey(&vtls.TLSCreateParams{
			Template: x509.Certificate{
				Subject: pkix.Name{},
			},
			BitSize: bitSize,
		})
		if err != nil {
			x.logger.Error().Msgf("Error in generating the cert: %v", err)
			cobra.CheckErr(err)
		}
		x.outputConfig.TLSCert = &setup.TLSCert{
			AutoGenerate: true,
		}
		x.outputConfig.TLSCert.Key = tlsCert.KeyPem
		x.outputConfig.TLSCert.Cert = tlsCert.CertPem
	} else if x.inputConfig.TLSCert.CertFile != "" && x.inputConfig.TLSCert.KeyFile != "" {
		tlsKeyBytes, err = os.ReadFile(x.inputConfig.TLSCert.KeyFile)
		if err != nil {
			x.logger.Error().Msgf("Error in reading the key file: %v", err)
			cobra.CheckErr(err)
		}
		tlsCertBytes, err = os.ReadFile(x.inputConfig.TLSCert.CertFile)
		if err != nil {
			x.logger.Error().Msgf("Error in reading the cert file: %v", err)
			cobra.CheckErr(err)
		}
		decodetlsKey := base64.StdEncoding.EncodeToString(tlsKeyBytes)
		if err != nil {
			x.logger.Error().Msgf("Error in decoding the key: %v", err)
			cobra.CheckErr(err)
		}
		decodetlsCert := base64.StdEncoding.EncodeToString(tlsCertBytes)
		if err != nil {
			x.logger.Error().Msgf("Error in decoding the cert: %v", err)
			cobra.CheckErr(err)
		}
		x.outputConfig.TLSCert = &setup.TLSCert{
			AutoGenerate: false,
		}
		x.outputConfig.TLSCert.Key = decodetlsKey
		x.outputConfig.TLSCert.Cert = decodetlsCert
	} else {
		x.outputConfig.TLSCert.Key = x.inputConfig.TLSCert.Key
		x.outputConfig.TLSCert.Cert = x.inputConfig.TLSCert.Cert
	}
}

func (x *SetupConfigurator) manageJWTAuthn() {
	if x.inputConfig.JWTAuthnSecrets == nil || x.inputConfig.JWTAuthnSecrets.PrivateJWTKey == "" || x.inputConfig.JWTAuthnSecrets.PublicJWTKey == "" {
		algo, err := pkg.EncryptionAlgorithm.Run()
		if err != nil {
			x.logger.Error().Msgf("Error in reading the algorithm: %v", err)
			cobra.CheckErr(err)
		}
		bitSizeS, err := pkg.EncryptionAlgorithmBitSize.Run()
		if err != nil {
			x.logger.Error().Msgf("Error in reading the bitsize: %v", err)
			cobra.CheckErr(err)
		}
		bitSize, err := strconv.Atoi(bitSizeS)
		if err != nil {
			x.logger.Error().Msgf("Error in converting the bitsize to int: %v, so switching to default", err)
		}
		encrypter, err := encryption.NewVapusDataJwtAuthn(&encryption.JWTAuthn{
			SigningAlgorithm: algo,
			Bitsize:          bitSize,
		})
		if err != nil {
			x.logger.Error().Msgf("Error in creating the encrypter: %v", err)
			cobra.CheckErr(err)
		}
		pbKey, pvKey, err := encrypter.GenerateKeys(bitSize)
		if err != nil {
			x.logger.Error().Msgf("Error in generating the keys: %v", err)
			cobra.CheckErr(err)
		}
		x.inputConfig.JWTAuthnSecrets.PrivateJWTKey = pvKey
		x.inputConfig.JWTAuthnSecrets.PublicJWTKey = pbKey
		x.inputConfig.JWTAuthnSecrets.SigningAlgorithm = algo
		x.inputConfig.JWTAuthnSecrets.Bitsize = bitSize
	}
}
