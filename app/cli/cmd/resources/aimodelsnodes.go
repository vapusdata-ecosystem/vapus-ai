package resources

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-ecosystem/vapusdata/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
)

func NewAIModelNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pkg.AIModelNodeResource,
		Short: "This command is interface to interact with the ai studio resources.",
		Long:  `This command is interface to interact with the ai studio resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			if plclient.MasterCommonFlags.La {
				plclient.MasterGlobals.VapusCtlClient.ListResourceActions(pkg.AIModelNodeResource)
				return
			}
			plclient.MasterGlobals.VapusCtlClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(plclient.MasterGlobals.CurrentAccessToken),
				Action:      plclient.MasterCommonFlags.Action,
				File:        plclient.MasterCommonFlags.File,
				Resource:    pkg.AIModelNodeResource,
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
