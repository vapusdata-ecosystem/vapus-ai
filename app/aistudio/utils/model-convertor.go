package utils

import (
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	models "github.com/vapusdata-ecosystem/vapusdata/core/models"
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
