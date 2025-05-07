package actions

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vapusdata-ecosystem/vapusdata/cli/cmd/resources"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
)

func NewDescribeCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   pkg.DescribeOps,
		Short: "This command will allow you to perform describe actions based on the resources provided.",
		Long:  `This command will allow you to perform describe actions based on the resources provided.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cobra.CheckErr(fmt.Errorf("no resource provided for this command, please select resource from result of this command -> " + pkg.APPNAME + " " + pkg.ExplainOps))
			}
		},
	}
	cmd.AddCommand(resources.NewUserCmd(), resources.NewOrganizationCmd(),
		resources.NewDataSourceCmd(), resources.NewAIModelNodeCmd())
	return cmd
}
