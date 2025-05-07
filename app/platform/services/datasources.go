package services

import (
	"context"
	"fmt"
	"slices"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	apppkgs "github.com/vapusdata-ecosystem/vapusdata/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	"github.com/vapusdata-ecosystem/vapusdata/core/pkgs/pbtools"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusdata/core/process"
	types "github.com/vapusdata-ecosystem/vapusdata/core/types"
	dmstores "github.com/vapusdata-ecosystem/vapusdata/platform/datastoreops"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"
	utils "github.com/vapusdata-ecosystem/vapusdata/platform/utils"
)

type DSIntAgentOpts func(*DataSourceAgent)

func WithDsAgentManagerRequest(managerRequest *pb.DataSourceManagerRequest) DSIntAgentOpts {
	return func(v *DataSourceAgent) {
		v.managerRequest = managerRequest
	}
}

func WithDsAgentGetterRequest(getterRequest *pb.DataSourceGetterRequest) DSIntAgentOpts {
	return func(v *DataSourceAgent) {
		v.getterRequest = getterRequest
	}
}

func WithDsAgentManagerAction(action string) DSIntAgentOpts {
	return func(v *DataSourceAgent) {
		v.Action = action
	}
}

type DataSourceAgent struct {
	result         *pb.DataSourceResponse
	managerRequest *pb.DataSourceManagerRequest
	getterRequest  *pb.DataSourceGetterRequest
	organization   *models.Organization
	dataSource     *models.DataSource
	*DMServices
	*processes.VapusInterfaceBase
	aiParams *models.AccountAIAttributes
}

func (x *DataSourceAgent) GetResult() *pb.DataSourceResponse {
	x.FinishAt = dmutils.GetEpochTime()
	return x.result
}

func (x *DataSourceAgent) LogAgent() {
	x.Logger.Info().Msgf("DataSourceAgent - %v action started at %v and finished at %v with status %v", x.AgentId, x.InitAt, x.FinishAt, x.Status)
}

func (x *DMServices) NewDataSourceAgent(ctx context.Context, opts ...DSIntAgentOpts) (*DataSourceAgent, error) {
	var err error
	var datasource *models.DataSource
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		x.logger.Error().Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}

	organization, err := x.DMStore.GetOrganization(ctx, vapusPlatformClaim[encryption.ClaimOrganizationKey], vapusPlatformClaim)
	if err != nil {
		return nil, dmerrors.DMError(apperr.ErrOrganization404, err)
	}

	agent := &DataSourceAgent{
		result:       &pb.DataSourceResponse{Output: &pb.DataSourceResponse_DataSourceOutput{}},
		organization: organization,
		dataSource:   datasource,
		DMServices:   x,
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			InitAt: dmutils.GetEpochTime(),
			// Ctx:       ctx,
			CtxClaim:  vapusPlatformClaim,
			AgentType: types.DATASOURCEAGENT.String(),
		},
	}
	agent.SetAgentId()
	for _, opt := range opts {
		opt(agent)
	}
	if agent.managerRequest.GetSpec() != nil && agent.managerRequest.GetSpec().GetDataSourceId() != "" {
		datasource, err = x.DMStore.GetDataSource(ctx, agent.managerRequest.GetSpec().GetDataSourceId(), vapusPlatformClaim)
		if err != nil {
			return nil, dmerrors.DMError(apperr.ErrDataSource404, err)
		}
	}
	agent.Logger = pkgs.GetSubDMLogger(types.DATASOURCEAGENT.String(), agent.AgentId)
	agent.aiParams = dmstores.GetAccountAIAttributes(ctx, x.DMStore, vapusPlatformClaim)
	return agent, nil
}

func (x *DataSourceAgent) Act(ctx context.Context, action string) error {
	x.Logger.Info().Msgf("DataSourceAgent - %v action ", x.Action)
	if action != "" {
		x.Action = action
	}
	var nctx context.Context
	switch x.Action {
	case mpb.ResourceLcActions_ADD.String():
		if x.managerRequest.GetSpec() == nil {
			return dmerrors.DMError(apperr.ErrInvalidDataSourceSpec, nil)
		}
		x.InitCreateResponse()
		x.dataSource = utils.DtNodeToObj(x.managerRequest)
		x.dataSource.Organization = x.CtxClaim[encryption.ClaimOrganizationKey]
		nctx, x.ContextCancel = pbtools.NewInCancelCtxWithAuthToken(ctx)
		go func() {
			err := x.createDataSource(nctx)
			if err != nil {
				x.Logger.Err(err).Ctx(ctx).Msgf("error while creating data source %v", x.dataSource.GetDataSourceId())
			}
		}()
		return nil
	case mpb.ResourceLcActions_UPDATE.String():
		if x.managerRequest.GetSpec() == nil {
			return dmerrors.DMError(apperr.ErrInvalidDataSourceSpec, nil)
		}
		x.dataSource = utils.DtNodeToObj(x.managerRequest)
		nctx, x.ContextCancel = pbtools.NewInCancelCtxWithAuthToken(ctx)
		go func() {
			err := x.updateDataSource(nctx)
			if err != nil {
				x.Logger.Err(err).Ctx(ctx).Msgf("error while updating data source %v", x.dataSource.GetDataSourceId())
			}
		}()

		return nil
	case mpb.ResourceLcActions_LIST.String():
		return x.listDataSources(ctx)
	case mpb.ResourceLcActions_GET.String():
		return x.describeDataSource(ctx)
	case mpb.ResourceLcActions_ARCHIVE.String():
		return x.archiveDataSource(ctx)
	default:
		x.Logger.Error().Msgf("invalid action %v", x.Action)
		return dmerrors.DMError(apperr.ErrInvalidManageAgentActions, nil) //nolint:wrapcheck
	}
}

func (x *DataSourceAgent) describeDataSource(ctx context.Context) error {
	x.Logger.Info().Msgf("describeDataSource for data source - %v", x.getterRequest.GetDataSourceId())
	ds, err := x.DMStore.GetDataSource(ctx, x.getterRequest.GetDataSourceId(), x.CtxClaim)
	if err != nil {
		x.Logger.Err(err).Ctx(ctx).Msgf("error while getting data source %v", x.dataSource.GetDataSourceId())
		return dmerrors.DMError(apperr.ErrDataSource404, err) //nolint:wrapcheck
	}
	// if ds.OwnerOrganization != x.CtxClaim[encryption.ClaimOrganizationKey] {
	// 	return dmerrors.DMError(apperr.ErrDataSource403, nil)
	// }
	// dsMD, err := x.dbStore.GetDataSourceMetaData(x.Ctx, x.datasource.VapusBase, true, x.CtxClaim)
	// if err != nil {
	// 	return []*models.DataSource{ds}, nil, err
	// }
	x.result.Output.DataSources = utils.DSObjToPb([]*models.DataSource{ds})
	return nil
}
func (x *DataSourceAgent) listDataSources(ctx context.Context) error {
	dataSourceIds := utils.GetFilterParams(x.getterRequest.GetSearchParam(), types.DataSourceSK.String())
	var filter string
	if len(dataSourceIds) > 0 {
		filter = fmt.Sprintf("organization = '%s' AND data_source_id IN (%s) AND status not in ('DELETED')", x.CtxClaim[encryption.ClaimOrganizationKey], dataSourceIds)
	} else {
		filter = fmt.Sprintf("organization = '%s' AND status not in ('DELETED')", x.CtxClaim[encryption.ClaimOrganizationKey])
	}

	dataSources, err := x.DMStore.ListDataSources(ctx, filter, x.CtxClaim)
	if err != nil {
		return err
	}
	x.result.Output.DataSources = utils.DSListObjToPb(dataSources)
	return nil
}

func (x *DataSourceAgent) updateDataSource(ctx context.Context) error {
	defer x.ContextCancel()
	if x.managerRequest.Spec.DataSourceId == "" {
		return dmerrors.DMError(apperr.ErrInvalidDataSourceId, nil)
	}
	exDataSource, err := x.DMStore.GetDataSource(ctx, x.managerRequest.Spec.DataSourceId, x.CtxClaim)
	if err != nil {
		x.Logger.Err(err).Ctx(ctx).Msgf("error while getting data source %v", x.dataSource.GetDataSourceId())
		return dmerrors.DMError(apperr.ErrDataSource404, err) //nolint:wrapcheck
	}
	if !slices.Contains(exDataSource.Editors, x.CtxClaim[encryption.ClaimUserIdKey]) || exDataSource.Status == mpb.CommonStatus_DELETED.String() {
		return dmerrors.DMError(apperr.ErrDataSource403, nil)
	}
	exDataSource.PreSaveUpdate(x.CtxClaim)
	// Uncomment when JWT is started in context of organization
	exDataSource.Status = mpb.CommonStatus_ACTIVE.String()
	for _, cd := range x.dataSource.NetParams.DsCreds {
		if cd.Credentials != nil {
			if cd.SecretName == types.EMPTYSTR {
				cd.SecretName = dmutils.GetSecretName("dataSource", exDataSource.VapusID, dmutils.GetStrEpochTime())
			}
			err = apppkgs.SaveCredentialsCreds(ctx, cd.SecretName, cd.Credentials, x.DMStore.VapusStore, x.Logger)
			if err != nil {
				return err
			}
			cd.IsAlreadyInSecretBs = true
			exDataSource.NetParams.DsCreds = append(exDataSource.NetParams.DsCreds, cd)

			cd.Credentials = nil
		}
	}
	err = x.DMStore.PutDataSource(ctx, exDataSource, x.CtxClaim)
	if err != nil {
		x.Logger.Err(err).Ctx(ctx).Msgf("error while updating data source %v", exDataSource.VapusID)
		return dmerrors.DMError(apperr.ErrCreateDataSource, err) //nolint:wrapcheck
	}
	x.result.Output.DataSources = utils.DSObjToPb([]*models.DataSource{exDataSource})
	return nil
}

func (x *DataSourceAgent) createDataSource(ctx context.Context) error {
	var err error
	defer x.ContextCancel()
	x.dataSource.SetDataSourceUuid()
	x.dataSource.PreSaveCreate(x.CtxClaim)
	// Uncomment when JWT is started in context of organization
	x.dataSource.Organization = x.CtxClaim[encryption.ClaimOrganizationKey]
	x.dataSource.Status = mpb.CommonStatus_CONFIGURING.String()
	x.dataSource.Editors = []string{x.CtxClaim[encryption.ClaimUserIdKey]}
	se, ok := types.StorageEngineMap[x.managerRequest.GetSpec().Attributes.GetServiceName()]
	if ok {
		x.dataSource.StorageEngine = se.String()
	}
	x.dataSource.ServiceProvider = x.managerRequest.GetSpec().Attributes.ServiceProvider.String()
	dsType, ok := types.DataSourceTypeMap[x.managerRequest.GetSpec().Attributes.GetServiceName()]
	if ok {
		x.dataSource.DataSourceType = dsType.String()
	} else {
		x.dataSource.DataSourceType = x.managerRequest.GetSpec().Attributes.GetServiceName().String()
	}
	err = x.DMStore.CreateDataSource(ctx, x.dataSource, x.CtxClaim)
	if err != nil {
		x.Logger.Err(err).Ctx(ctx).Msgf("error while creating data source %v", x.dataSource)
		return dmerrors.DMError(apperr.ErrCreateDataSource, err) //nolint:wrapcheck
	}
	for _, cd := range x.dataSource.NetParams.DsCreds {
		if !cd.IsAlreadyInSecretBs {
			if cd.SecretName == types.EMPTYSTR {
				cd.SecretName = dmutils.GetSecretName("dataSource", x.dataSource.VapusID, dmutils.GetStrEpochTime())
			}
			err = apppkgs.SaveCredentialsCreds(ctx, cd.SecretName, cd.Credentials, x.DMStore.VapusStore, x.Logger)
			if err != nil {
				return err
			}
			cd.Credentials = nil
			cd.IsAlreadyInSecretBs = true
		}
	}

	x.dataSource.Status = mpb.CommonStatus_SYNCING.String()
	err = x.DMStore.PutDataSource(ctx, x.dataSource, x.CtxClaim)
	if err != nil {
		x.Logger.Err(err).Ctx(ctx).Msgf("error while updating data source %v", x.dataSource)
		return dmerrors.DMError(apperr.ErrCreateDataSource, err) //nolint:wrapcheck
	}

	err = x.DMStore.PutDataSource(ctx, x.dataSource, x.CtxClaim)
	if err != nil {
		x.Logger.Err(err).Ctx(ctx).Msgf("error while updating data source %v", x.dataSource)
		return dmerrors.DMError(apperr.ErrCreateDataSource, err) //nolint:wrapcheck
	}
	x.SetCreateResponse(mpb.Resources_DATASOURCES, x.dataSource.VapusID)
	return nil
}

func (x *DataSourceAgent) archiveDataSource(ctx context.Context) error {
	if x.getterRequest.GetDataSourceId() == "" {
		return dmerrors.DMError(apperr.ErrInvalidDataSourceSpec, nil)
	}
	ds, err := x.DMStore.GetDataSource(ctx, x.getterRequest.GetDataSourceId(), x.CtxClaim)
	if err != nil {
		x.Logger.Err(err).Ctx(ctx).Msgf("error while getting data source %v", x.dataSource.GetDataSourceId())
		return dmerrors.DMError(apperr.ErrDataSource404, err) //nolint:wrapcheck
	}
	if ds.Organization != x.CtxClaim[encryption.ClaimOrganizationKey] || ds.CreatedBy != x.CtxClaim[encryption.ClaimUserIdKey] {
		return dmerrors.DMError(apperr.ErrDataSource403, nil)
	}
	ds.PreSaveDelete(x.CtxClaim)
	err = x.DMStore.PutDataSource(ctx, ds, x.CtxClaim)
	if err != nil {
		x.Logger.Err(err).Ctx(ctx).Msgf("error while archiving data source %v", ds.VapusID)
		return dmerrors.DMError(apperr.ErrArchiveDatasource, err) //nolint:wrapcheck
	}
	return nil
}

func (x *DataSourceAgent) putDataSourceSilent(ctx context.Context, status string, ds *models.DataSource, err error) {
	if err != nil {
		ds.ErrorLogs = append(ds.ErrorLogs, dmerrors.DMError(apperr.ErrSyncingDatasource, err).Error())
	}
	ds.Status = status
	x.Logger.Info().Msgf("Status of data product - %v", ds.Status)
	err = x.DMServices.DMStore.PutDataSource(ctx, ds, x.CtxClaim)
	if err != nil {
		x.Logger.Err(err).Msg("error while updating datasource in db")
	}
}
