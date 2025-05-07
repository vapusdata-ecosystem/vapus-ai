package utils

import (
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dpb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
)

// Utils func to convert proto account creation request to local obj
func AccPboObj(request *pb.AccountManagerRequest) *models.Account {
	return &models.Account{}

}

// Utils func to convert proto organization  creation request to local obj
func DmNodeToObj(request *dpb.OrganizationManagerRequest) *models.Organization {
	return (&models.Organization{}).ConvertFromPb(request.GetSpec())
}

// Utils func to convert local organization  object to proto object
func DmArrToPb(objList []*models.Organization) []*mpb.Organization {
	if objList == nil {
		return nil
	}
	var pbList []*mpb.Organization
	for _, obj := range objList {
		pbList = append(pbList, obj.ConvertToPb())
	}
	return pbList
}

func DmListToPb(objList []*models.Organization) []*mpb.Organization {
	if objList == nil {
		return nil
	}
	var pbList []*mpb.Organization
	for _, obj := range objList {
		pbList = append(pbList, obj.ConvertToListingPb())
	}
	return pbList
}

// Utils func to convert proto datasource  creation request to local obj
func DtNodeToObj(request *dpb.DataSourceManagerRequest) *models.DataSource {
	return (&models.DataSource{}).ConvertFromPb(request.GetSpec())
}

func DSObjToPb(objList []*models.DataSource) []*mpb.DataSource {
	if objList == nil {
		return nil
	}
	var pbList []*mpb.DataSource
	for _, obj := range objList {
		pbList = append(pbList, obj.ConvertToPb())
	}
	return pbList
}

func DSListObjToPb(objList []*models.DataSource) []*mpb.DataSource {
	if objList == nil {
		return nil
	}
	var pbList []*mpb.DataSource
	for _, obj := range objList {
		pbList = append(pbList, obj.ConvertToListingPb())
	}
	return pbList
}

func Dtn2Obj(request *mpb.DataSource) *models.DataSource {
	return (&models.DataSource{}).ConvertFromPb(request)
}

func DmUArToPb(users []*models.Users, organization string) []*mpb.User {
	if users == nil {
		return nil
	}
	var pbList []*mpb.User
	for _, obj := range users {
		pbList = append(pbList, obj.ConvertToPb(organization))
	}
	return pbList
}

func DmUListingToPb(users []*models.Users, organization string) []*mpb.User {
	if users == nil {
		return nil
	}
	var pbList []*mpb.User
	for _, obj := range users {
		pbList = append(pbList, obj.ConvertToListingPb(organization))
	}
	return pbList
}

func DmUPbToObj(users *mpb.User) *models.Users {
	if users == nil {
		return nil
	}
	return (&models.Users{}).ConvertFromPb(users)
}

func DPPLObj2Pb(dws []*models.Plugin) []*mpb.Plugin {
	if dws == nil {
		return nil
	}
	var dpList []*mpb.Plugin
	for _, dp := range dws {
		dpList = append(dpList, dp.ConvertToPb())
	}
	return dpList
}

func SMPb2Obj(dws *mpb.SecretStore) *models.SecretStore {
	if dws == nil {
		return nil
	}
	return (&models.SecretStore{}).ConvertFromPb(dws)
}

func SMObj2Pb(dws *models.SecretStore) *mpb.SecretStore {
	if dws == nil {
		return nil
	}
	return dws.ConvertToPb()
}
