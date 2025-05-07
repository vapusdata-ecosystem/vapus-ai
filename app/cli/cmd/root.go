package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"

	"github.com/vapusdata-ecosystem/vapusai/cli/cmd/actions"
	plclient "github.com/vapusdata-ecosystem/vapusai/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
)

var (
	defaultCfgFolder   = ".vapusdata"
	defaultCfgFileType = "toml"
	defaultCfgFileName = "config"
	version            = "0.0.1"
)

func NewRootCmd() *cobra.Command {
	rootCmd = &cobra.Command{
		Use:     pkg.APPNAME,
		Version: version,
		Short:   "vapusctl is a cli tool tht provides an interface to interact with different services offered by VapusData Platform",
		Long:    `This cli tool will allow you to perform different operation on VapusData platform under the current context. `,
		Run: func(cmd *cobra.Command, args []string) {
			plclient.MasterGlobals.Logger.Info().Msg("Welcome to VapusData CLI")
			// cmd.Help()
		},
	}
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		setGlobalPersist()
		_, exist := ignoreConnMap[cmd.Name()]
		if !exist {
			initCurrentContextInstance()
		}
		log.Println("PersistentPreRunE called -------------> ", plclient.MasterCommonFlags, plclient.MasterGlobals)

		return nil
	}
	ss := plclient.CtlCommonFlags{}
	rootCmd.PersistentFlags().StringVar(&ss.CfgFile, "config", "", "config file (default is $HOME/.vapusdata/config.toml)")
	rootCmd.PersistentFlags().BoolVar(&ss.DebugLogFlag, "debug", true, "Enable debug/verbose logging mode")
	rootCmd.PersistentFlags().StringVarP(&ss.Action, "action", "a", "", "Action for the platform that should be executed on current resource with params in a file")
	rootCmd.PersistentFlags().StringVarP(&ss.File, "file", "f", "", "File containing the parameters for the action")
	rootCmd.PersistentFlags().StringVarP(&ss.OutputFile, "output", "o", "", "File containing the output of the action")
	rootCmd.PersistentFlags().StringVarP(&ss.SearchQuery, "query", "q", "", "Search query for the resource")
	rootCmd.PersistentFlags().BoolVarP(&ss.La, "listactions", "l", false, "List down all the actions that can be performed on the current resource")

	_ = viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	_ = viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	_ = viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("debug", ss.DebugLogFlag)
	plclient.MasterCommonFlags = &ss
	log.Println("PersistentPreRunE outcalled -------------> ", plclient.MasterCommonFlags, plclient.MasterGlobals, ss)
	rootCmd.AddCommand(NewClearCmd(),
		NewConfigCmd(),
		actions.NewExplainOps(),
		NewAuthCmd(),
		actions.NewGetCmd(),
		NewConnectCmd(),
		actions.NewDescribeCmd(),
		NewRequestSpecCmd(),
		actions.NewActCmd(),
		NewOperatorCmd())
	return rootCmd
}

func defaultConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return homeDir + "/" + defaultCfgFolder

}

func initConfig() {
	Logger.Debug().Msgf("Initiating config... at %v", plclient.MasterCommonFlags.CfgFile)
	// Read config either from plclient.MasterCommonFlags.CfgFile or from home directory!
	if plclient.MasterCommonFlags.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(plclient.MasterCommonFlags.CfgFile)
	} else {
		// Find home directory to get/create config file in there.
		plclient.MasterCommonFlags.CfgFile = filepath.Join(defaultConfigDir(), defaultCfgFileName+"."+defaultCfgFileType)

		// Search config in home directory with name ".vapusdata" (without extension).
		viper.AddConfigPath(defaultConfigDir())
		viper.SetConfigName(defaultCfgFileName)
		viper.SetConfigType(defaultCfgFileType)

		// Write the file only and gracefully handles if file already exists
		cfgErr := viper.SafeWriteConfig()
		if cfgErr != nil {
			existCfgError := viper.ConfigFileAlreadyExistsError("")
			if !errors.As(cfgErr, &existCfgError) {
				cobra.CheckErr(cfgErr)
			}
		}
		cobra.CheckErr(viper.ReadInConfig())
		viper.SetConfigFile(plclient.MasterCommonFlags.CfgFile)
	}
	// Viper read the config file on desired path
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func setGlobalPersist() {
	pkg.InitLogger(plclient.MasterCommonFlags.DebugLogFlag)
	plclient.MasterGlobals = &plclient.GlobalsPersists{
		Logger:             pkg.DmLogger,
		CfgFile:            plclient.MasterCommonFlags.CfgFile,
		DebugLogFlag:       plclient.MasterCommonFlags.DebugLogFlag,
		CfgDir:             filepath.Dir(plclient.MasterCommonFlags.CfgFile),
		Ctx:                context.Background(),
		AgentsActions:      plclient.AgentGoals,
		CurrentAccessToken: currentAccessToken,
		CurrentIdToken:     currentIdToken,
	}
}

func initCurrentContextInstance() {
	// Initialize the platform client
	currentContext := viper.GetString(currentContextKey)
	if currentContext == "" {
		return
	}
	currentContextParams := viper.GetStringMapString(currentContextParamsKey)
	if currentContextParams == nil {
		plclient.MasterGlobals.Logger.Error().Msg("No context found, please add a context")
	} else {
		if plclient.MasterGlobals.VapusCtlClient == nil {
			client, err := plclient.NewPlatFormClient(currentContextParams, plclient.MasterGlobals.Logger)
			if err != nil {
				plclient.MasterGlobals.Logger.Error().Msgf("Error connecting to current context platform instance: %v", err)
				cobra.CheckErr(pkg.ErrVapusDataPlatformNotConnected)
				return
			}
			plclient.MasterGlobals.VapusCtlClient = client
			plclient.MasterGlobals.CurrentContext = currentContextParams["name"]
			// plclient.MasterGlobals.ResourceActionMap = plclient.MasterGlobals.AgentsActions
		}
	}
}
