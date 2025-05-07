package actions

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vapusdata-ecosystem/vapusai/cli/cmd/resources"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
)

func NewAddUsersCmd() *cobra.Command {

	getCmd := &cobra.Command{
		Use:   pkg.AddUserOps,
		Short: "This command will allow you to add users to the system",
		Long:  `This command will allow you to add users to the system`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cobra.CheckErr(fmt.Errorf("no resource provided for this command, please select resource from result of this command -> " + pkg.APPNAME + " " + pkg.ExplainOps))
			}
		},
	}
	getCmd.AddCommand(resources.NewAccountCmd(), resources.NewUserCmd(), resources.NewOrganizationCmd())
	return getCmd
}
