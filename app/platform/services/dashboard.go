package services

import (
	"context"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

func (dms *DMServices) DashboardSvc(ctx context.Context) (*mpb.OrganizationDashboard, error) {

	resp := &mpb.OrganizationDashboard{}

	return resp, nil
}
