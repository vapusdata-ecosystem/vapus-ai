package plclient

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jedib0t/go-pretty/v6/table"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
	appconfigs "github.com/vapusdata-ecosystem/vapusai/core/app/configs"
	gwcl "github.com/vapusdata-ecosystem/vapusai/core/app/httpcls"
	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

func (x *VapusCtlClient) RetrieveLoginURL() (*pb.LoginHandlerResponse, error) {
	cl := gwcl.NewUserManagementClient(x.GwClient, nil, nil)
	return cl.LoginHandler(&gwcl.HttpRequestGeneric{
		Method: types.GET,
		Body:   []byte{},
	})
}

func (x *VapusCtlClient) RetrieveAccessToken(code, host string) (string, string, error) {
	log.Println("Retrieving access token for code: ", code)
	log.Println("Retrieving access token for url: ", host)
	req := &pb.LoginCallBackRequest{Code: code, Host: host}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return "", "", err
	}
	cl := gwcl.NewUserManagementClient(x.GwClient, nil, nil)
	result, err := cl.LoginCallback(&gwcl.HttpRequestGeneric{
		Method: types.POST,
		Body:   reqBytes,
	})
	if err != nil {
		return "", "", err
	}
	// result, err := x.UserConn.LoginCallback(context.Background(), &pb.LoginCallBackRequest{Code: code, Host: host})
	// if err != nil {
	// 	return "", "", err
	// }
	return result.Token.GetAccessToken(), result.Token.GetIdToken(), nil
}

func (x *VapusCtlClient) RetrievePlatformAccessToken(ctx context.Context, token, domain string) (string, error) {
	reqBytes, err := x.protojsonMarshal.Marshal(&pb.AccessTokenInterfaceRequest{Organization: domain, Utility: pb.AccessTokenAgentUtility_ORGANIZATION_LOGIN})
	if err != nil {
		return "", err
	}
	cl := gwcl.NewUserManagementClient(x.GwClient, nil, nil)
	result, err := cl.AccessTokenInterface(&gwcl.HttpRequestGeneric{
		Token: token,
		Body:  reqBytes,
	})
	if err != nil {
		return "", err
	}
	return result.Token.AccessToken, nil
}

func (x *VapusCtlClient) ListUser(ctx context.Context) error {
	cl := gwcl.NewUserManagementClient(x.GwClient, nil, nil)
	result, err := cl.UserGetter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	}, &pb.UserGetterRequest{})
	if err != nil {
		return err
	}
	pkg.LogTitles("Logged-in User Info: ", x.logger)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"UserId", "Name", "Added On", "Organization", "Roles"})
	for _, user := range result.Output.Users {
		tw.AppendRow(table.Row{user.UserId, user.DisplayName, user.InvitedOn, user.Roles[0].OrganizationId, user.Roles[0].Role})
	}
	tw.AppendSeparator()
	tw.Render()
	return nil
}

func (x *VapusCtlClient) DescribeUser(ctx context.Context) error {
	cl := gwcl.NewUserManagementClient(x.GwClient, nil, nil)
	result, err := cl.UserGetter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	}, &pb.UserGetterRequest{
		Action: pb.UserGetterActions_GET_USER,
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("Logged-in User Info: ", x.logger)
	x.PrintDescribe(result.Output, "user")
	return nil
}

func (x *VapusCtlClient) ListPlatformSpec() {
	fm := []interface{}{}
	for _, f := range mpb.ContentFormats_name {
		fm = append(fm, f)
	}
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Resource Name", "Spec Available", "Formats Available", "Generate command"})
	for _, spec := range mpb.Resources_name {
		if spec == mpb.Resources_INVALID_REQUEST_OBJECT.String() {
			continue
		}
		tw.AppendRow(table.Row{spec, true, pkg.NewListWritter(fm, list.StyleBulletSquare).Render(), pkg.APPNAME + " spec --name " + spec + " --generate-file=true --format yaml"})
		tw.AppendSeparator()
	}
	tw.Render()
}

func (x *VapusCtlClient) GeneratePlatformSpec(token, specName, format string, withFakeData bool) error {
	log.Println("Generating spec for ", specName)
	log.Println("Generating spec in format ", x)
	val := appconfigs.SpecMap[mpb.Resources(mpb.Resources_value[specName])]
	specVal, err := filetools.GenericMarshaler(val, format)
	if err != nil {
		return err
	}

	fileName := strings.ToLower(specName + "." + strings.ToLower(format))
	x.logger.Info().Msgf("Sample %v spec generated with file name - %v \n", specName, fileName)
	err = os.WriteFile(fileName, specVal, 0644)
	if err != nil {
		return err
	}
	return nil
}

// func (x *VapusCtlClient) GetSvcInfo(ctx context.Context) error {
// 	response, err := cl.PlatformServicesInfo(&gwcl.HttpRequestGeneric{
// 		Token: x.ActionHandler.AccessToken,
// 	}, &pb.PlatformServicesRequest{})
// 	if err != nil {
// 		return err
// 	}
// 	pkg.LogTitles("Platform Services Info:", x.logger)
// 	tw := pkg.NewTableWritter()
// 	tw.AppendHeader(table.Row{"Service Name", "Port", "Address", "Tag"})
// 	for _, item := range response.NetworkParams {
// 		tw.AppendRow(table.Row{item.SvcName, item.Port, item.SvcAddr, item.SvcTag})
// 		tw.AppendSeparator()
// 	}
// 	tw.Render()
// 	return nil
// }
