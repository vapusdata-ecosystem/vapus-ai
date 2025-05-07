package cmd

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-ecosystem/vapusdata/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
)

var (
	newContextName, instance, addAuthToken, namespace string
	credsConfig                                       = "credentials"
	setCurrent                                        bool
	port                                              *int32
)

func NewAddContextCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   addContextCmd,
		Short: "This command is to add a context of VapusData Platform instance",
		Long:  `This command is to add a context of VapusData Platform instance`,
		Run: func(cmd *cobra.Command, args []string) {
			addContext()
		},
	}

	cmd.PersistentFlags().StringVar(&instance, "url", "", "URL of the VapusData Platform instance")
	cmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace where the VapusData Platform instance is running in K8S")
	port = cmd.PersistentFlags().Int32("port", 0, "Port on which the VapusData Platform instance is running")
	// contextAddCmd.PersistentFlags().StringVar(&addAuthToken, "authtoken", "", "Auth token of the VapusData Platform instance")
	cmd.PersistentFlags().StringVar(&newContextName, "name", "", "Name for the context without spaces , default is '{url}-context-{unixtimestamp}'")
	cmd.PersistentFlags().BoolVar(&setCurrent, "set-current", false, "flag to set the current context while adding the context")
	return cmd
}

func addContext() {
	allContexts := viper.GetStringSlice("contexts")
	if len(allContexts) == 0 {
		viper.Set("contexts", []string{newContextName})
	} else {
		for _, context := range allContexts {
			if context == newContextName {
				plclient.MasterGlobals.Logger.Error().Msgf("Context with name %v already exists", newContextName)
				return
			}
		}
		allContexts = append(allContexts, newContextName)
		viper.Set("contexts", allContexts)
	}
	plclient.MasterGlobals.Logger.Info().Msgf("Adding context with port pointer: %v", *port)
	if !strings.Contains(instance, "http://") || !strings.Contains(instance, "https://") {
		instance = "http://" + instance
	}
	i, err := url.ParseRequestURI(instance)
	if err != nil {
		cobra.CheckErr(err)
	}
	plclient.MasterGlobals.Logger.Info().Msgf("Adding context with instance: %v", i.Scheme)
	if strings.Contains(instance, "http://") {
		instance = strings.TrimPrefix(instance, "http://")
	}
	if strings.Contains(instance, "https://") {
		instance = strings.TrimPrefix(instance, "https://")
	}
	contextParam := map[string]string{pkg.NAME: newContextName, pkg.URL: instance, pkg.PORT: strconv.Itoa(int(*port)), pkg.ADDRESS: i.Hostname(), pkg.NAMESPACE: namespace}

	plclient.MasterGlobals.Logger.Info().Msgf("Adding context with context: %v", contextParam)

	viper.Set(newContextName, contextParam)

	if setCurrent {
		viper.Set(currentContextKey, contextParam["name"])
		viper.Set(currentContextParamsKey, contextParam)
	}

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
		plclient.MasterGlobals.Logger.Info().Msg("Context added successfully!...")
	}

	cobra.CheckErr(viper.ReadInConfig())
}
