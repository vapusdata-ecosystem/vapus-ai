package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	plclient "github.com/vapusdata-ecosystem/vapusai/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
)

// dmCtxCmd represents the dmCtx command
func NewGetContextsCmd() *cobra.Command {

	contextCmd := &cobra.Command{
		Use:   pkg.ContextsOps,
		Short: "This command is list all the context available.",
		Long:  `This command is list all the context available.`,
		Run: func(cmd *cobra.Command, args []string) {
			listContext()
		},
	}

	return contextCmd
}

func listContext() {
	viper.SetConfigFile(plclient.MasterCommonFlags.CfgFile)
	allContexts := viper.GetStringSlice("contexts")
	currentContext := viper.GetString("current-context")
	if len(allContexts) == 0 {
		plclient.MasterGlobals.Logger.Info().Msgf("No context found. Please add a context first using '%v %v %v' command", pkg.APPNAME, pkg.ConfigResource, addContextCmd)
		os.Exit(0)
	}
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Context Name", "Current Context"})

	for _, context := range allContexts {
		if context == currentContext {
			tw.AppendRow(table.Row{context, "*"})
			tw.AppendSeparator()
			continue
		}
		tw.AppendRow(table.Row{context, ""})
		tw.AppendSeparator()
	}

	tw.Render()
}
