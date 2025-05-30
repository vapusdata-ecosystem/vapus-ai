package actions

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vapusdata-ecosystem/vapusai/cli/cmd/resources"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
)

func NewValidateCmd() *cobra.Command {

	getCmd := &cobra.Command{
		Use:   pkg.ValidateOps,
		Short: "This command will allow you to perform listing opertions based on the resources provided.",
		Long:  `This command will allow you to perform different listing operations based on the resources provided.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cobra.CheckErr(fmt.Errorf("no resource provided for this command, please select resource from result of this command -> " + pkg.APPNAME + " " + pkg.ExplainOps))
			}
		},
	}
	getCmd.AddCommand(
		resources.NewDataSourceCmd(),
	)
	return getCmd
}
