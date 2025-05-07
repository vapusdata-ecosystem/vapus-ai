package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
)

func NewOperatorCmd() *cobra.Command {

	operatorCmd := &cobra.Command{
		Use:   pkg.OperatorOps,
		Short: "This command will allow you to manage vapusdata application operations on kubernetes cluster.",
		Long:  `This command will allow you to manage vapusdata application operations on kubernetes cluster.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cobra.CheckErr(fmt.Errorf("no action provided for this command, please select resource from result of this command -> " + pkg.OperatorOps + " --help"))
			}
		},
	}
	operatorCmd.AddCommand(NewInstallerCmd(), NewInstallerSpecGenCmd(), NewInstallerSetupCmd())
	return operatorCmd
}
