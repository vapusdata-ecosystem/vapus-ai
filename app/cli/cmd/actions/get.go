package actions

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/vapusdata-ecosystem/vapusdata/cli/cmd/resources"
	plclient "github.com/vapusdata-ecosystem/vapusdata/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
)

func NewGetCmd() *cobra.Command {
	log.Println("NewGetCmd called", *plclient.MasterCommonFlags, plclient.MasterGlobals)
	getCmd := &cobra.Command{
		Use:   pkg.GetOps,
		Short: "This command will allow you to perform listing opertions based on the resources provided.",
		Long:  `This command will allow you to perform different listing operations based on the resources provided.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Get command called with resource: ", args)
			if len(args) < 1 {
				cobra.CheckErr(fmt.Errorf("no resource provided for this command, please select resource from result of this command -> " + pkg.APPNAME + " " + pkg.ExplainOps))
			}

			plclient.MasterGlobals.Logger.Info().Msgf("Get command called with resource: %s", args[0])
		},
	}
	getCmd.AddCommand(resources.NewAccountCmd(), resources.NewUserCmd(), resources.NewOrganizationCmd(),
		resources.NewDataSourceCmd(), resources.NewAIModelNodeCmd())
	return getCmd
}
