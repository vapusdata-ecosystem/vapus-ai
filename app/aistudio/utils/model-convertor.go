package utils

import (
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
)

func AIMNPb2Obj(obj *mpb.AIModelNode) *models.AIModelNode {
	return (&models.AIModelNode{}).ConvertFromPb(obj)
}

func AIMNArrM2Pb(obj []*models.AIModelNode) []*mpb.AIModelNode {
	var res []*mpb.AIModelNode
	for _, v := range obj {
		res = append(res, v.ConvertToPb())
	}
	return res
}

func AIMNList2Pb(obj []*models.AIModelNode) []*mpb.AIModelNode {
	var res []*mpb.AIModelNode
	for _, v := range obj {
		res = append(res, v.ConvertToListingPb())
	}
	return res
}

func AIMNListPb2Obj(obj []*mpb.AIModelNode) []*models.AIModelNode {
	var res []*models.AIModelNode
	for _, v := range obj {
		res = append(res, AIMNPb2Obj(v))
	}
	return res
}

func AIPRArrObj2Pb(obj []*models.AIPrompt) []*mpb.AIPrompt {
	var res []*mpb.AIPrompt
	for _, v := range obj {
		res = append(res, v.ConvertToPb())
	}
	return res
}

func AIPRListObj2Pb(obj []*models.AIPrompt) []*mpb.AIPrompt {
	var res []*mpb.AIPrompt
	for _, v := range obj {
		res = append(res, v.ConvertToListingPb())
	}
	return res
}

func AIPROPb2Obj(obj *mpb.AIPrompt) *models.AIPrompt {
	return (&models.AIPrompt{}).ConvertFromPb(obj)
}

func AIPRPb2Obj(obj []*mpb.AIPrompt) []*models.AIPrompt {
	var res []*models.AIPrompt
	for _, v := range obj {
		res = append(res, AIPROPb2Obj(v))
	}
	return res
}

func AIGDPb2Obj(obj []*mpb.AIGuardrails) []*models.AIGuardrails {
	var res []*models.AIGuardrails
	for _, v := range obj {
		res = append(res, (&models.AIGuardrails{}).ConvertFromPb(v))
	}
	return res
}

func AIGDArrObjToPb(obj []*models.AIGuardrails) []*mpb.AIGuardrails {
	var res []*mpb.AIGuardrails
	for _, v := range obj {
		res = append(res, v.ConvertToPb())
	}
	return res
}

func AIGDListObjToPb(obj []*models.AIGuardrails) []*mpb.AIGuardrails {
	var res []*mpb.AIGuardrails
	for _, v := range obj {
		res = append(res, v.ConvertToListingPb())
	}
	return res
}

func AISCLToPbChat(obj []*models.AIStudioChat) []*pb.AIStudioChat {
	result := make([]*pb.AIStudioChat, 0)
	for _, v := range obj {
		result = append(result, v.ConvertToListingPb())
	}
	return result
}

func AISCOToPbChat(obj *models.AIStudioChat) *pb.AIStudioChat {
	return obj.ConvertToPb()
}

func FbALToPbAgent(agents []*models.VapusAgents) []*mpb.VapusAgent {
	var result []*mpb.VapusAgent
	for _, c := range agents {
		result = append(result, c.ConvertToPb())
	}
	return result
}

// Utils func to convert proto account creation request to local obj
func AccPboObj(request *pb.AccountManagerRequest) *models.Account {
	return &models.Account{}

}

// Utils func to convert proto organization  creation request to local obj
func DmNodeToObj(request *pb.OrganizationManagerRequest) *models.Organization {
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
func DtNodeToObj(request *pb.DataSourceManagerRequest) *models.DataSource {
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
