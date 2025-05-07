package resources

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	plclient "github.com/vapusdata-ecosystem/vapusai/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
)

func NewUserCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   pkg.UserResource,
		Short: "This command is interface to interact with the platform for users resources.",
		Long:  `This command is interface to interact with the platform for users resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			plclient.MasterGlobals.VapusCtlClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(plclient.MasterGlobals.CurrentAccessToken),
				Action:      userActions(cmd.Parent().Use),
				Resource:    pkg.UserResource,
			}
			err := plclient.MasterGlobals.VapusCtlClient.HandleAction()
			if err != nil {
				cobra.CheckErr(err)
			}

			defer plclient.MasterGlobals.VapusCtlClient.Close()
		},
	}
	return cmd
}

func userActions(parentCmd string) string {
	switch parentCmd {
	case pkg.GetOps:
		return pb.UserGetterActions_LIST_USERS.String()
	case pkg.DescribeOps:
		return pb.UserGetterActions_GET_USER.String()
	default:
		return ""
	}
}

// func getuser() {
// 	err := plclient.MasterGlobals.VapusCtlClient.ListActions(pb.UserAgentOperations_GET_USER.String(), viper.GetString(currentAccessToken))
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func describeUser() {

// 	err := plclient.MasterGlobals.VapusCtlClient.DescribeActions(pb.UserAgentOperations_GET_USER.String(), viper.GetString(plclient.MasterGlobals.CurrentAccessToken), "")
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
