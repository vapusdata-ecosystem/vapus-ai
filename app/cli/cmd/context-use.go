package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-ecosystem/vapusdata/cli/internals/platform"
)

var (
	currentContext       string
	currentContextParams map[string]string
)

const (
	currentContextKey       = "current-context"
	currentContextParamsKey = "current-context-params"
)

// dmCtxCmd represents the dmCtx command
func NewUseContextCmd() *cobra.Command {
	contextUseCmd := &cobra.Command{
		Use:   useContextCmd,
		Short: "This command will set the context to current VapusData platform instance.",
		Long:  `This command will set the context to current VapusData platform instance.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 || args[0] == "" || len(args) > 1 {
				plclient.MasterGlobals.Logger.Info().Msg("Please provide a context name to use the specific context")
				cobra.CheckErr(fmt.Errorf("no context provided for this command"))
			}
			currentContext = args[0]
			useContext()
		},
	}

	return contextUseCmd
}

func useContext() {
	plclient.MasterGlobals.Logger.Info().Msgf("Setting the current context to - %v", currentContext)
	viper.Set(currentContextKey, currentContext)

	currentContextParams = viper.GetStringMapString(currentContext)

	viper.Set(currentContextParamsKey, currentContextParams)

	viper.SetConfigType(defaultCfgFileType)
	viper.AddConfigPath(plclient.MasterCommonFlags.CfgFile)

	// Write the file only and gracefully handles if file already exists
	cfgErr := viper.SafeWriteConfig()
	if cfgErr != nil {
		cfgErr = viper.WriteConfig()
		if cfgErr != nil {
			cobra.CheckErr(cfgErr)
			plclient.MasterGlobals.Logger.Error().Msgf("Error addting the context: %v", cfgErr)
		}
	} else {
		plclient.MasterGlobals.Logger.Info().Msgf("Current context is set successfully to - %v", currentContext)
	}
	cobra.CheckErr(viper.ReadInConfig())

}
