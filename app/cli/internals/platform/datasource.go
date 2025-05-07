package plclient

import (
	"context"
	"log"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
	gwcl "github.com/vapusdata-ecosystem/vapusai/core/app/httpcls"
	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
)

func (x *VapusCtlClient) HandleDataSourceAct(ctx context.Context) error {
	var fileBytes []byte
	var err error
	if x.ActionHandler.File != "" {
		x.protoyamlUnMarshal.Path = x.ActionHandler.File
		fileBytes, err = filetools.ReadFile(x.ActionHandler.File)
		if err != nil {
			return err
		}
		// fileBytes, err = yaml.YAMLToJSON(bbytes)
		// if err != nil {
		// 	return err
		// }
		x.inputFormat = strings.ToUpper(filetools.GetConfFileType(x.ActionHandler.File))
	}
	// if len(x.ActionHandler.Identifier) < 1 && x.ActionHandler.File == "" {
	// 	return pkg.ErrNoArgs
	// }
	requestSpec := &pb.DataSourceManagerRequest{}
	err = x.protoyamlUnMarshal.Unmarshal(fileBytes, requestSpec)
	if err != nil {
		return err
	}
	if requestSpec.Spec == nil {
		return pkg.ErrInvalidRequestSpec
	}
	fBytes, err := x.protojsonMarshal.Marshal(requestSpec)
	if err != nil {
		return err
	}
	log.Println("RequestSpec: ", string(fBytes))
	return nil
}

func (x *VapusCtlClient) ListDataSources(ctx context.Context) error {
	cl := gwcl.NewDataSourceClient(x.GwClient, nil)
	result, err := cl.Getter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	}, &pb.DataSourceGetterRequest{})

	if err != nil {
		return err
	}
	pkg.LogTitles("List of dataSources registered in current domain: ", x.logger)
	if len(result.Output.GetDataSources()) == 0 {
		pkg.LogTitles("\nNo Data Sources found", x.logger)
	} else {
		tw := pkg.NewTableWritter()
		tw.AppendHeader(table.Row{"Name", "Id", "Data Stores", "Storage Engine", "Service Name", "status", "Organization"})
		for _, source := range result.Output.GetDataSources() {
			tw.AppendRow(table.Row{source.Name, source.DataSourceId, strings.Join(source.NetParams.Databases, ","), source.Attributes.StorageEngine, source.Attributes.ServiceName, source.Status, source.ResourceBase.Organization})
			tw.AppendSeparator()
		}
		tw.Render()
	}
	return nil
}

func (x *VapusCtlClient) DescribeDataSources(ctx context.Context, iden string) error {
	cl := gwcl.NewDataSourceClient(x.GwClient, nil)
	result, err := cl.Getter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	}, &pb.DataSourceGetterRequest{
		DataSourceId: iden,
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("Data Source Info: ", x.logger)
	x.PrintDescribe(result.Output.DataSources[0], "datasource")
	return nil
}

func (x *VapusCtlClient) ConfigureDataSource(ctx context.Context, requestSpec []byte) error {
	cl := gwcl.NewDataSourceClient(x.GwClient, nil)
	response, err := cl.Create(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
		Body:  requestSpec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Data Source Configured Successfully with ID:", x.logger, response.Output.DataSources[0].DataSourceId)
	return nil
}
