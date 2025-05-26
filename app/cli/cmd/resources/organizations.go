package resources

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	plclient "github.com/vapusdata-ecosystem/vapusai/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
)

func NewOrganizationCmd() *cobra.Command {
	var dmId string
	cmd := &cobra.Command{
		Use:   pkg.OrganizationResource,
		Short: "This command is interface to interact with the platform for domain resources.",
		Long:  `This command is interface to interact with the platform for domain resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			defer plclient.MasterGlobals.VapusCtlClient.Close()
			if plclient.MasterCommonFlags.La {
				plclient.MasterGlobals.VapusCtlClient.ListResourceActions("domains")
				return
			}
			// spinner := pkg.GetSpinner(39)
			// spinner.Start()
			// defer spinner.Stop()
			log.Println(plclient.MasterGlobals.CurrentAccessToken, "-----------------------")
			plclient.MasterGlobals.VapusCtlClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(plclient.MasterGlobals.CurrentAccessToken),
				Action:      getOrganizationAction(cmd.Parent().Use, plclient.MasterCommonFlags.Action),
				File:        plclient.MasterCommonFlags.File,
				Params:      map[string]string{pkg.OrganizationKey: dmId},
				Resource:    pkg.OrganizationResource,
			}
			err := plclient.MasterGlobals.VapusCtlClient.HandleAction()
			if err != nil {
				cobra.CheckErr(err)
			}
		},
	}
	cmd.PersistentFlags().StringVar(&dmId, "domain", "", "Data product Id to perform the action on")
	return cmd
}

func getOrganizationAction(parentCmd string, action string) string {
	switch parentCmd {
	case pkg.GetOps:
		return mpb.ResourceLcActions_LIST.String()
	case pkg.DescribeOps:
		return mpb.ResourceLcActions_GET.String()
	case pkg.ActOps:
		return action
	default:
		return pkg.ErrInvalidAction.Error()
	}
}
