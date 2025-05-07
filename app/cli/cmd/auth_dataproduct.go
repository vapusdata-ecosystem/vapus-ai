package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-ecosystem/vapusai/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
)

var (
	loginDataProduct string
)

// authCmd represents the auth command
func NewDataProductAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pkg.DataProductResource,
		Short: "Login to the VapusData platform instance using Authenticator",
		Long:  `This command is used to login to the VapusData platform`,
		Run: func(cmd *cobra.Command, args []string) {
			generateDataProductAccessToken(args)
		},
	}
	cmd.Flags().StringVar(&loginDataProduct, "dataproduct", "", "uses provided data product context for logging in")
	return cmd
}

func generateDataProductAccessToken(args []string) {
	var err error
	if loginDataProduct == "" {
		cobra.CheckErr(pkg.ErrMissingDataProductLogin)
	}
	// accessToken := viper.GetString(currentAccessToken)
	// newDPAccessToken, err := plclient.MasterGlobals.VapusCtlClient.RetrieveDataProductAccessToken(context.Background(), accessToken, loginDataProduct)
	// if err != nil {
	// 	plclient.MasterGlobals.Logger.Error().Err(err).Msg("failed to retrieve platform access token")
	// 	cobra.CheckErr(err)
	// }
	newDPAccessToken := viper.GetString(currentAccessToken)
	viper.Set(currentProductAccessToken, newDPAccessToken)
	err = viper.WriteConfig()
	if err != nil {
		cobra.CheckErr(err)
	}
	plclient.MasterGlobals.Logger.Info().Msgf("successfully logged in to data product - %v.\n Your access token is %v ,Please copy it or save it somewhere safe", loginDataProduct, newDPAccessToken)
}
