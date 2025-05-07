package actions

import (
	"fmt"

	list "github.com/jedib0t/go-pretty/v6/list"
	table "github.com/jedib0t/go-pretty/v6/table"
	text "github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	plclient "github.com/vapusdata-ecosystem/vapusdata/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
)

func NewExplainOps() *cobra.Command {

	cmd := &cobra.Command{
		Use:   pkg.ExplainOps,
		Short: "This command is to list all agents. You will have to provide the resource name as an argument to perform certain actions on that resource",
		Long: `
	You can use this command to list all the agents or get a specific resource.
			`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				cobra.CheckErr(fmt.Errorf("no resource provided for this command"))
			}
			getValidResources()
		},
	}

	return cmd
}

func getValidResources() {
	var xx string

	xx = text.FormatUpper.Apply("Goal Based Agents with list of their actions: ")
	xx = text.Underline.Sprintf(xx)
	plclient.MasterGlobals.Logger.Info().Msgf("\n%v", xx)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Agents", "Actions", "Version", "Chaining Support"})
	for resource, operations := range plclient.MasterGlobals.AgentsActions {
		tw.AppendRow(table.Row{resource, pkg.NewListWritter(operations, list.StyleMarkdown).Render(), "v1alpha1", true})
		tw.AppendSeparator()
	}
	tw.Render()
}
