package plclient

import (
	"context"
	"log"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	aipb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
	gwcl "github.com/vapusdata-ecosystem/vapusai/core/app/httpcls"
	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
)

func (x *VapusCtlClient) HandleAIStudioAct(ctx context.Context) error {
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
	requestSpec := &aipb.AIModelNodeManagerRequest{}
	err = x.protoyamlUnMarshal.Unmarshal(fileBytes, requestSpec)
	if err != nil {
		log.Println(err)
		return err
	}
	if requestSpec.Spec == nil {
		return pkg.ErrInvalidRequestSpec
	}
	fBytes, err := x.protojsonMarshal.Marshal(requestSpec)
	if err != nil {
		return err
	}
	return x.AIModelNodeManager(ctx, fBytes)
}

func (x *VapusCtlClient) getAIModelNodes(ctx context.Context) error {
	cl := gwcl.NewAIModelsClient(x.GwClient, nil)
	result, err := cl.Getter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	}, &aipb.AIModelNodeGetterRequest{})
	if err != nil {
		return err
	}
	if len(result.GetOutput().GetAiModelNodes()) < 1 {
		pkg.LogTitles("\nNo AI Model Nodes found", x.logger)
		return nil
	}
	pkg.LogTitles("List of AI Model Nodes: ", x.logger)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Node ID", "Name", "Service provider", "Scope", "Organization", "Generative Models Count", "Embedding Models Count"})
	for _, node := range result.GetOutput().GetAiModelNodes() {
		tw.AppendRow(table.Row{node.GetModelNodeId(), node.GetName(), node.GetAttributes().ServiceProvider,
			node.GetAttributes().Scope, node.GetAttributes().ApprovedOrganizations, len(node.GetAttributes().GenerativeModels),
			len(node.GetAttributes().EmbeddingModels)})
	}
	tw.AppendSeparator()
	tw.Render()
	return nil
}

func (x *VapusCtlClient) describeAIModelNode(ctx context.Context, modelsNodeId string) error {
	log.Println("Describe AI Model Node")
	cl := gwcl.NewAIModelsClient(x.GwClient, nil)
	result, err := cl.Getter(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
	}, &aipb.AIModelNodeGetterRequest{
		AiModelNodeId: modelsNodeId,
	})
	if err != nil {
		return err
	}
	pkg.LogTitles("Details of AI Model Node: ", x.logger)
	x.PrintDescribe(result.Output.AiModelNodes[0], "AI Model Node")
	return nil
	// obj, err := x.protoyamlMarshal.Marshal(result.Output.AiModelNodes[0])
	// if err != nil {
	// 	return err
	// }
	// return pkg.ParseAndBuildYamlTable(obj)
}

func (x *VapusCtlClient) AIModelNodeManager(ctx context.Context, requestSpec []byte) error {
	cl := gwcl.NewAIModelsClient(x.GwClient, nil)
	response, err := cl.Create(&gwcl.HttpRequestGeneric{
		Token: x.ActionHandler.AccessToken,
		Body:  requestSpec,
	})
	if err != nil {
		return err
	}
	pkg.LogTitlesf("AI Model Node Configured Successfully with ID:", x.logger, response.Output.AiModelNodes[0].ModelNodeId)
	return nil
}
