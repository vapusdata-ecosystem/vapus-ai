package resources

// import (
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/viper"
// 	dpb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
// 	plclient "github.com/vapusdata-ecosystem/vapusai/cli/internals/platform"
// 	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
// )

// func NewMetadataCmd() *cobra.Command {
// 	var dsId string
// 	cmd := &cobra.Command{
// 		Use:   pkg.MetaDataResource,
// 		Short: "This command is interface to interact with the platform for metadata resources.",
// 		Long:  `This command is interface to interact with the platform for metadata resources.`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			defer plclient.MasterGlobals.VapusCtlClient.Close()
// 			if plclient.MasterCommonFlags.La{
// 				plclient.MasterGlobals.VapusCtlClient.ListResourceActions("metadata")
// 				return
// 			}
// 			resAct := getMetadataAction(cmd.Parent().Use, plclient.MasterCommonFlags.Action)
// 			spinner := pkg.GetSpinner(36)
// 			spinner.Prefix = "Agent is running"
// 			spinner.Start()
// 			plclient.MasterGlobals.VapusCtlClient.ActionHandler = plclient.ActionHandlerOpts{
// 				ParentCmd:   cmd.Parent().Use,
// 				Args:        args,
// 				AccessToken: viper.GetString(plclient.MasterGlobals.CurrentAccessToken),
// 				Action:      resAct,
// 				File:        plclient.MasterCommonFlags.File,
// 				Params:      map[string]string{pkg.DatasourceKey: dsId},
// 			}

// 			err := plclient.MasterGlobals.VapusCtlClient.HandleAction()
// 			if err != nil {
// 				spinner.Stop()
// 				cobra.CheckErr(err)
// 			}

// 			spinner.Stop()

// 		},
// 	}
// 	cmd.PersistentFlags().StringVar(&dsId, "datasource", "", "Data product Id to perform the action ons")
// 	return cmd
// }

// func getMetadataAction(parentCmd string, action string) string {
// 	switch parentCmd {
// 	case pkg.GetOps:
// 		return dpb.DataSourceAgentActions_LIST_DATASOURCE.String()
// 	case pkg.DescribeOps:
// 		return dpb.DataSourceAgentActions_DESCRIBE_DATASOURCE.String()
// 	case pkg.ActOps:
// 		return action
// 	default:
// 		return pkg.ErrInvalidAction.Error()
// 	}
// }

// // func getDataSource() {
// // 	err := plclient.MasterGlobals.VapusCtlClient.ListActions(pb.DataSourceAgentActions_LIST_DATASOURCE.String(), viper.GetString(currentAccessToken))
// // 	if err != nil {
// // 		cobra.CheckErr(err)
// // 	}
// // }

// // func describeDataSource(args []string) {
// // 	if len(args) < 1 {
// // 		cobra.CheckErr("Invalid number of arguments, please provide the dataSource ID")
// // 	}
// // 	err := plclient.MasterGlobals.VapusCtlClient.DescribeActions(pb.DataSourceAgentActions_DESCRIBE_DATASOURCE.String(), viper.GetString(plclient.MasterGlobals.CurrentAccessToken), args[0])
// // 	if err != nil {
// // 		cobra.CheckErr(err)
// // 	}
// // }
