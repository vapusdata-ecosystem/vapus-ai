package plclient

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
)

func (x *VapusCtlClient) ListResourceActions(resource string) {
	xx := text.FormatUpper.Apply("Actions for Resource: ")
	xx = text.Underline.Sprintf(xx)
	x.logger.Info().Msgf("\n%v", xx)
	tw := pkg.NewTableWritter()
	tw.AppendHeader(table.Row{"Actions", "Version", "Commands"})
	for _, action := range x.ResourceActionMap[resource] {
		if strings.Contains(action.(string), "INVALID") {
			continue
		}
		tw.AppendRow(table.Row{action, "v1alpha1", pkg.APPNAME + " act " + resource + " --file <Input File> "})
		tw.AppendSeparator()
	}
	tw.Render()
}

func (x *VapusCtlClient) HandleAction() error {
	if x.ActionHandler.ParentCmd == "" {
		return errors.New("invalid operations")
	}
	switch x.ActionHandler.ParentCmd {
	case pkg.GetOps:
		return x.HandleGet()
	case pkg.DescribeOps:
		return x.HandleDescription()
	case pkg.ActOps:
		return x.HandleAct()
	default:
		return pkg.ErrInvalidAction
	}
}

func (x *VapusCtlClient) HandleAct() error {
	ctx := context.Background()
	newCtx := pkg.GetBearerCtx(ctx, x.ActionHandler.AccessToken)
	switch x.ActionHandler.Resource {
	case pkg.AccountResource:
		return x.HandleAccountAct(newCtx)
	case pkg.DataSourceResource:
		return x.HandleDataSourceAct(newCtx)
	case pkg.OrganizationResource:
		return x.HandleOrganizationAct(newCtx)
	case pkg.AIModelNodeResource:
		return x.HandleAIStudioAct(newCtx)
	default:
		return pkg.ErrInvalidResource
	}
}

func (x *VapusCtlClient) HandleGet() error {
	log.Println("Getting resource..............................")
	ctx := context.Background()
	newCtx := pkg.GetBearerCtx(ctx, x.ActionHandler.AccessToken)
	switch x.ActionHandler.Resource {
	case pkg.AccountResource:
		return x.ListAccount(newCtx)
	case pkg.UserResource:
		return x.ListUser(newCtx)
	case pkg.OrganizationResource:
		return x.ListOrganizations(newCtx)
	case pkg.DataSourceResource:
		return x.ListDataSources(newCtx)
	case pkg.AIModelNodeResource:
		return x.getAIModelNodes(newCtx)
	default:
		return pkg.ErrInvalidResource
	}
}

func (x *VapusCtlClient) HandleDescription() error {
	var iden string = ""
	if x.ActionHandler.Resource != pkg.OrganizationResource {
		if len(x.ActionHandler.Args) < 1 {
			return pkg.ErrNoArgs
		}
		iden = pkg.GetDescId(x.ActionHandler.Args)
	}
	ctx := context.Background()
	newCtx := pkg.GetBearerCtx(ctx, x.ActionHandler.AccessToken)
	log.Println("Describing with resource identifier: ", iden)
	switch x.ActionHandler.Resource {
	case pkg.UserResource:
		return x.DescribeUser(newCtx)
	case pkg.OrganizationResource:
		return x.DescribeOrganizations(newCtx)
	case pkg.DataSourceResource:
		return x.DescribeDataSources(newCtx, iden)
	case pkg.AIModelNodeResource:
		return x.describeAIModelNode(newCtx, iden)
	default:
		return pkg.ErrInvalidAction
	}
}
