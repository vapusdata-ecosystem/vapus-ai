package cmd

import (
	"os"

	cobra "github.com/spf13/cobra"
	plclient "github.com/vapusdata-ecosystem/vapusai/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
)

func NewClearCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pkg.ClearOps,
		Version: version,
		Short:   "This command will clear the config file of vapus cli that involves the context of the different vapusdata platform instance instances.",
		Long:    `This command will clear the config file of vapus cli that involves the context of the different vapusdata platform instance instances.`,
		Run: func(cmd *cobra.Command, args []string) {
			clearConfigDir(plclient.MasterGlobals.CfgDir)
		},
	}
	return cmd
}

func clearConfigDir(cfgDir string) {
	plclient.MasterGlobals.Logger.Debug().Msgf("clearing config from dir %v", cfgDir)
	files, err := os.ReadDir(cfgDir)
	if err != nil {
		cobra.CheckErr(err)
	}
	for _, file := range files {
		if file.IsDir() {
			plclient.MasterGlobals.Logger.Debug().Msgf("clearing config from sub dir %v", file.Name())
			clearConfigDir(file.Name())
		}
		plclient.MasterGlobals.Logger.Debug().Msgf("clearing config file %v", file.Name())
		err := os.Remove(cfgDir + "/" + file.Name())
		if err != nil {
			cobra.CheckErr(err)
		}
	}
}
