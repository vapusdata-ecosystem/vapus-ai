package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-ecosystem/vapusai/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
)

var (
	loginOrganization string
)

// authCmd represents the auth command
func NewOrganizationAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pkg.OrganizationResource,
		Short: "Login to the VapusData platform instance using Authenticator",
		Long:  `This command is used to login to the VapusData platform`,
		Run: func(cmd *cobra.Command, args []string) {
			generateOrganizationAccessToken(args)
		},
	}

	cmd.Flags().StringVar(&loginOrganization, "domain", "", "uses provided domain context for logging in")
	return cmd
}

func generateOrganizationAccessToken(args []string) {
	var err error
	if loginOrganization == "" {
		plclient.MasterGlobals.Logger.Info().Msg("no domain provided for login, system will login to default domain")
	}
	accessToken := viper.GetString(currentAccessToken)
	newAccessToken, err := plclient.MasterGlobals.VapusCtlClient.RetrievePlatformAccessToken(context.Background(), accessToken, loginOrganization)
	if err != nil {
		plclient.MasterGlobals.Logger.Error().Err(err).Msg("failed to retrieve platform access token")
		cobra.CheckErr(err)
	}

	viper.Set(currentAccessToken, newAccessToken)
	err = viper.WriteConfig()
	if err != nil {
		plclient.MasterGlobals.Logger.Error().Err(err).Msg("failed to write new access token to config")
		cobra.CheckErr(err)
	}
	plclient.MasterGlobals.Logger.Info().Msgf("successfully logged in to domain - %v", loginOrganization)
	defer plclient.MasterGlobals.VapusCtlClient.Close()
}
