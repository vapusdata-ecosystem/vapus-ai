package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-ecosystem/vapusai/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
)

var (
	format, specName                      string
	generateFile, listSpecs, withFakeData bool
)

func NewRequestSpecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pkg.SpecsOps,
		Short: "This command is interface to list or generate sample request files for the platform interfaces.",
		Long:  `This command is interface to list or generate sample request files for the platform interfaces.`,
		Run: func(cmd *cobra.Command, args []string) {
			specOps()
		},
	}
	cmd.PersistentFlags().StringVar(&specName, "name", "", "Name of the resource spec to generate")
	cmd.PersistentFlags().BoolVar(&generateFile, "generate-file", false, "Flag to generate the request spec in requested file format and creates a file")
	cmd.Flags().BoolVar(&listSpecs, "ls", false, "Flag that will list all name of all the available specs")
	// cmd.PersistentFlags().BoolVar(&withFakeData, "with-fake", false, "Flag that will populate the generated spec with fake data")
	cmd.PersistentFlags().StringVar(&format, "format", "", "format of the file to generate the request spec")
	return cmd
}

func specOps() {
	if listSpecs {
		plclient.MasterGlobals.VapusCtlClient.ListPlatformSpec()
		return
	}
	if specName == "" {
		cobra.CheckErr("spec name is required")
	}
	err := plclient.MasterGlobals.VapusCtlClient.GeneratePlatformSpec(viper.GetString(currentAccessToken), specName, format, withFakeData)
	if err != nil {
		cobra.CheckErr(err)
	}
}
