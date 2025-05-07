package plclient

import (
	"context"
	"log"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
	gwcl "github.com/vapusdata-ecosystem/vapusdata/core/app/httpcls"
	filetools "github.com/vapusdata-ecosystem/vapusdata/core/tools/files"
)

func (x *VapusCtlClient) HandleOrganizationAct(ctx context.Context) error {
	var fileBytes []byte
	var err error
	if x.ActionHandler.File != "" {
		x.protoyamlUnMarshal.Path = x.ActionHandler.File
		fileBytes, err = filetools.ReadFile(x.ActionHandler.File)
		if err != nil {
			return err
		}
		x.inputFormat = strings.ToUpper(filetools.GetConfFileType(x.ActionHandler.File))
	}
	requestSpec := &pb.OrganizationManagerRequest{}
	err = x.protoyamlUnMarshal.Unmarshal(fileBytes, requestSpec)
	if err != nil {
		return err
	}
	fBytes, err := x.protojsonMarshal.Marshal(requestSpec)
	if err != nil {
		return err
	}
	log.Println("Request Spec: ", string(fBytes))
	return nil
}

func (x *VapusCtlClient) ConfigureOrganization(ctx context.Context, requestSpec []byte) error {
	cl := gwcl.NewOrganizationClient(x.GwClient, nil)
	response, err := cl.Create(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
		Body:  requestSpec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Organization Configured Successfully with ID:", x.logger, response.Output.Organizations[0].OrganizationId)
	return nil
}

func (x *VapusCtlClient) UpgradeOrganizationArtifacts(ctx context.Context, requestSpec []byte) error {
	cl := gwcl.NewOrganizationClient(x.GwClient, nil)
	response, err := cl.UpgradeOS(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
		Body:  requestSpec,
	}, &pb.OrganizationGetterRequest{})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Artifacts successfully upgraded for domain :", x.logger, response.Output.Organizations[0].OrganizationId)
	return nil
}

func (x *VapusCtlClient) PatchOrganization(ctx context.Context, requestSpec []byte) error {
	cl := gwcl.NewOrganizationClient(x.GwClient, nil)
	response, err := cl.Update(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
		Body:  requestSpec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Organization deployment infra is successfully upgraded for domain :", x.logger, response.Output.Organizations[0].OrganizationId)
	return nil
}

func (x *VapusCtlClient) ListOrganizations(ctx context.Context) error {
	log.Println("x.ActionHandler.AccessToken: =========>>>>>>>>>>>>>>>>>", x.ActionHandler.AccessToken)
	cl := gwcl.NewOrganizationClient(x.GwClient, nil)
	result, err := cl.Getter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	}, &pb.OrganizationGetterRequest{})
	if err != nil {
		return err
	}
	// result, err := vapushttpcl.DataProductManager(x.ActionHandler.AccessToken, reqBytes)
	// if err != nil {
	// 	return err
	// }
	pkg.LogTitles("List of Organization Info: ", x.logger)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Name", "Id", "Type", "Total users", "status"})
	for _, dm := range result.Output.Organizations {
		tw.AppendRow(table.Row{dm.Name, dm.OrganizationId, dm.OrganizationType.String(), len(dm.GetUsers()), dm.Status})
		tw.AppendSeparator()
	}
	tw.Render()
	return nil
}

func (x *VapusCtlClient) DescribeOrganizations(ctx context.Context) error {
	cl := gwcl.NewOrganizationClient(x.GwClient, nil)
	result, err := cl.Getter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	}, &pb.OrganizationGetterRequest{})
	if err != nil {
		return err
	}
	pkg.LogTitles("Logged-In user's Organization Info: ", x.logger)
	x.PrintDescribe(result.Output.Organizations[0], "domain")
	return nil
}
