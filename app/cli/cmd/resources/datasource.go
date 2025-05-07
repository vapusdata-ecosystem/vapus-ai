package resources

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	plclient "github.com/vapusdata-ecosystem/vapusdata/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
)

func NewDataSourceCmd() *cobra.Command {
	var dsId string
	cmd := &cobra.Command{
		Use:   pkg.DataSourceResource,
		Short: "This command is interface to interact with the platform for dataSource resources.",
		Long:  `This command is interface to interact with the platform for dataSource resources.`,
		Run: func(cmd *cobra.Command, args []string) {
			defer plclient.MasterGlobals.VapusCtlClient.Close()
			if plclient.MasterCommonFlags.La {
				plclient.MasterGlobals.VapusCtlClient.ListResourceActions("datasources")
				return
			}
			// spinner := pkg.GetSpinner(36)
			// spinner.Start()
			plclient.MasterGlobals.VapusCtlClient.ActionHandler = plclient.ActionHandlerOpts{
				ParentCmd:   cmd.Parent().Use,
				Args:        args,
				AccessToken: viper.GetString(plclient.MasterGlobals.CurrentAccessToken),
				Action:      getDatasourceAction(cmd.Parent().Use, plclient.MasterCommonFlags.Action),
				File:        plclient.MasterCommonFlags.File,
				Params:      map[string]string{pkg.DatasourceKey: dsId},
				Resource:    pkg.DataSourceResource,
			}
			// defer spinner.Stop()
			err := plclient.MasterGlobals.VapusCtlClient.HandleAction()
			if err != nil {
				cobra.CheckErr(err)
			}
			log.Println("Action completed successfully")
		},
	}
	cmd.PersistentFlags().StringVar(&dsId, "datasource", "", "Data source Id to perform the action ons")
	return cmd
}

func getDatasourceAction(parentCmd string, action string) string {
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

// func getDataSource() {
// 	err := plclient.MasterGlobals.VapusCtlClient.ListActions(pb.DataSourceAgentActions_LIST_DATASOURCE.String(), viper.GetString(currentAccessToken))
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }

// func describeDataSource(args []string) {
// 	if len(args) < 1 {
// 		cobra.CheckErr("Invalid number of arguments, please provide the dataSource ID")
// 	}
// 	err := plclient.MasterGlobals.VapusCtlClient.DescribeActions(pb.DataSourceAgentActions_DESCRIBE_DATASOURCE.String(), viper.GetString(plclient.MasterGlobals.CurrentAccessToken), args[0])
// 	if err != nil {
// 		cobra.CheckErr(err)
// 	}
// }
