package resources

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-ecosystem/vapusdata/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
)

func NewAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pkg.AccountResource,
		Short: "This command is interface to interact with the platform for account resources.",
		Long:  `This command is interface to interact with the platform for account resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			if plclient.MasterCommonFlags.La {
				plclient.MasterGlobals.VapusCtlClient.ListResourceActions("account")
				return
			}
			resAct := getAccountAction(cmd.Parent().Use, plclient.MasterCommonFlags.Action)
			// spinner := pkg.GetSpinner(36)
			// spinner.Start()
			plclient.MasterGlobals.VapusCtlClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(plclient.MasterGlobals.CurrentAccessToken),
				Action:      resAct,
				File:        plclient.MasterCommonFlags.File,
				Resource:    pkg.AccountResource,
			}
			log.Println("Action: ", resAct)
			err := plclient.MasterGlobals.VapusCtlClient.HandleAction()
			// spinner.Stop()
			if err != nil {
				cobra.CheckErr(err)
			}
			defer plclient.MasterGlobals.VapusCtlClient.Close()

		},
	}
	return cmd
}

func getAccountAction(parentCmd string, action string) string {
	switch parentCmd {
	case pkg.ActOps:
		return action
	default:
		return ""
	}
}

// func accountActions(parentCmd string) {
// 	switch parentCmd {
// 	case pkg.GetOps:
// 		getAccount()
// 	case pkg.DescribeOps:
// 		describeAccount()
// 	default:
// 		cobra.CheckErr("Invalid action")
// 	}
// }

// func getAccount() {
// 	err := plclient.MasterGlobals.VapusCtlClient.ListActions(pb.AccountAgentActions_LIST_ACCOUNT.String(), viper.GetString(currentAccessToken))
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func describeAccount() {
// 	err := plclient.MasterGlobals.VapusCtlClient.DescribeActions(pb.AccountAgentActions_LIST_ACCOUNT.String(), viper.GetString(plclient.MasterGlobals.CurrentAccessToken), "")
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
