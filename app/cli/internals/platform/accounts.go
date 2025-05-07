package plclient

import (
	"context"
	"log"
	"strings"

	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
	gwcl "github.com/vapusdata-ecosystem/vapusdata/core/app/httpcls"
	filetools "github.com/vapusdata-ecosystem/vapusdata/core/tools/files"
)

func (x *VapusCtlClient) HandleAccountAct(ctx context.Context) error {
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
	requestSpec := &pb.AccountManagerRequest{}
	err = x.protoyamlUnMarshal.Unmarshal(fileBytes, requestSpec)
	if err != nil {
		return err
	}
	fBytes, err := x.protojsonMarshal.Marshal(requestSpec)
	if err != nil {
		return err
	}
	return x.AccountManagerClient(ctx, fBytes)
}

func (x *VapusCtlClient) ListAccount(ctx context.Context) error {
	log.Println("Listing account")
	cl := gwcl.NewVapusdataServiceClient(x.GwClient, nil)
	result, err := cl.AccountGetter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	})
	if err != nil {
		return err
	}
	obj, err := x.protoyamlMarshal.Marshal(result.Output)
	if err != nil {
		return err
	}
	return pkg.ParseAndBuildYamlTable(obj)
	// maps, err := dmutils.StructToMap(result.Output)
	// pkg.LogTitles("Account Info of current loggedin instance: ", x.logger)
	// tw := pkg.NewTableWritter()
	// // tw.AppendHeader(table.Row{"Account name", "Account Id", "Data Store", "Secret Store", "Artifact Store", "Authn Method", "Authz"})
	// tw.AppendHeader(table.Row{"Heading", "Value"})
	// for k, v := range maps {
	// 	tw.AppendRow(table.Row{k, v})
	// }
	// // tw.AppendRow(table.Row{result.Output.Name, result.Output.AccountId, result.Output.BackendDataStorage.BesService, result.Output.BackendSecretStorage.BesService,
	// // 	result.Output.ArtifactStorage.BesService, result.Output.AuthnMethod, result.Output.DmAccessJwtKeys.SigningAlgorithm})
	// tw.AppendSeparator()
	// tw.Render()
	// return nil
}

func (x *VapusCtlClient) AccountManagerClient(ctx context.Context, requestSpec []byte) error {
	log.Println(string(requestSpec))
	cl := gwcl.NewVapusdataServiceClient(x.GwClient, nil)
	response, err := cl.AccountManager(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
		Body:  requestSpec,
	})
	// response, err := x.PlConn.AccountManager(ctx, requestSpec)
	if err != nil {
		return err
	}
	pkg.LogTitlesf("Account Configured Successfully with ID:", x.logger, response.Output.AccountId)
	return nil
}
