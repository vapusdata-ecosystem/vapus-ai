package clients

import (
	"github.com/labstack/echo/v4"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dpb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
)

func (s *GrpcClient) GetCurrentOrganization(eCtx echo.Context) (*mpb.Organization, error) {
	result, err := s.OrganizationConn.Get(s.SetAuthCtx(eCtx), &dpb.OrganizationGetterRequest{})
	if err != nil || len(result.Output.GetOrganizations()) == 0 {
		s.logger.Err(err).Msg("error while getting domain current logged in info")
		return nil, err
	}
	return result.Output.GetOrganizations()[0], nil
}

func (s *GrpcClient) GetOrganizations(eCtx echo.Context) []*mpb.Organization {
	result, err := s.OrganizationConn.List(s.SetAuthCtx(eCtx), &dpb.OrganizationGetterRequest{})
	if err != nil {
		s.logger.Err(err).Msg("error while getting list of domains")
		return []*mpb.Organization{}
	}
	return result.Output.GetOrganizations()
}

func (x *GrpcClient) DescribeDataSources(eCtx echo.Context, dsId string) *mpb.DataSource {
	result, err := x.DatasourceConn.Get(x.SetAuthCtx(eCtx), &dpb.DataSourceGetterRequest{
		DataSourceId: dsId,
	})

	if err != nil || result.Output == nil || len(result.Output.GetDataSources()) != 1 {
		x.logger.Err(err).Msg("error while getting DataSource details")
		return nil
	}
	return result.Output.GetDataSources()[0]
}

func (x *GrpcClient) ListDataSources(eCtx echo.Context) []*mpb.DataSource {
	result, err := x.DatasourceConn.List(x.SetAuthCtx(eCtx), &dpb.DataSourceGetterRequest{})

	if err != nil || result.Output == nil {
		x.logger.Err(err).Msg("error while getting DataSource list")
		return nil
	}
	return result.Output.GetDataSources()
}
