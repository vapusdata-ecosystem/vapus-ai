package services

import (
	"context"
	"encoding/json"
	"slices"

	"github.com/databricks/databricks-sql-go/logger"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/aistudio/utils"
	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	apppkgss "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	models "github.com/vapusdata-ecosystem/vapusai/core/models"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusai/core/process"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

type AIModelIntAgent struct {
	*processes.VapusInterfaceBase
	managerRequest      *pb.AIModelNodeManagerRequest
	getterRequest       *pb.AIModelNodeGetterRequest
	usageInsightRequest *pb.AIModelNodeInsightsRequest
	result              *pb.AIModelNodeResponse
	insightResult       *pb.AIModelNodeInsightsResponse
	dmStore             *aidmstore.AIStudioDMStore
}

type AIModelIntAgentOpts func(*AIModelIntAgent)

func WithModelManagerRequest(managerRequest *pb.AIModelNodeManagerRequest) AIModelIntAgentOpts {
	return func(v *AIModelIntAgent) {
		v.managerRequest = managerRequest
	}
}

func WithModelGetterRequest(getterRequest *pb.AIModelNodeGetterRequest) AIModelIntAgentOpts {
	return func(v *AIModelIntAgent) {
		v.getterRequest = getterRequest
	}
}

func WithModelManagerAction(action mpb.ResourceLcActions) AIModelIntAgentOpts {
	return func(v *AIModelIntAgent) {
		v.Action = action.String()
	}
}
func WithModelNodeInsightsRequest(usageInsightRequest *pb.AIModelNodeInsightsRequest) AIModelIntAgentOpts {
	return func(v *AIModelIntAgent) {
		v.usageInsightRequest = usageInsightRequest
	}
}

func (s *AIStudioServices) NewAIModelIntAgent(ctx context.Context, opts ...AIModelIntAgentOpts) (*AIModelIntAgent, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.Logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, dmerrors.DMError(encryption.ErrInvalidJWTClaims, nil)
	}
	nodeAgent := &AIModelIntAgent{
		result:  &pb.AIModelNodeResponse{Output: &pb.AIModelNodeResponse_AIModelNodeResponse{}},
		dmStore: s.DMStore,
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			CtxClaim: vapusPlatformClaim,
			// Ctx:      ctx,
			InitAt: dmutils.GetEpochTime(),
		},
		insightResult: &pb.AIModelNodeInsightsResponse{},
	}
	nodeAgent.SetAgentId()
	for _, opt := range opts {
		opt(nodeAgent)
	}
	nodeAgent.Logger = pkgs.GetSubDMLogger(types.AISTUDIONODE.String(), nodeAgent.AgentId)
	return nodeAgent, nil
}

func (v *AIModelIntAgent) GetAgentId() string {
	return v.AgentId
}

func (v *AIModelIntAgent) GetResult() *pb.AIModelNodeResponse {
	return v.result
}

func (v *AIModelIntAgent) GetInsightResult() *pb.AIModelNodeInsightsResponse {
	return v.insightResult
}
func (v *AIModelIntAgent) Act(ctx context.Context) error {
	switch v.GetAction() {
	case mpb.ResourceLcActions_CREATE.String():
		return v.configureAIModelNode(ctx)
	case mpb.ResourceLcActions_UPDATE.String():
		return v.updateAIModelNode(ctx)
	case mpb.ResourceLcActions_GET.String():
		return v.describeAIModelNode(ctx)
	case mpb.ResourceLcActions_LIST.String():
		return v.listAIModelNodes(ctx)
	case mpb.ResourceLcActions_SYNC.String():
		return v.syncAIModelNode(ctx)
	case mpb.ResourceLcActions_ARCHIVE.String():
		return v.archiveAIModelNode(ctx)
	case mpb.ResourceLcActions_INSIGHTS.String():
		return v.listUsageInsights(ctx)
	default:
		v.Logger.Error().Msg("invalid action for AIModelIntAgent")
		return dmerrors.DMError(apperr.ErrInvalidAction, nil)
	}
}

func (v *AIModelIntAgent) updateCachePool(modelNode *models.AIModelNode, add, update bool) {
	if modelNode.SecurityGuardrails != nil {
		for _, guardrail := range modelNode.SecurityGuardrails.Guardrails {
			pkgs.GuardrailPoolManager.AddGuardrails(modelNode.VapusID, guardrail)
		}
	}
	_, _ = pkgs.AIModelNodeConnectionPoolManager.GetorSetConnection(modelNode, add, update)
	_, _ = pkgs.AIModelNodeConnectionPoolManager.GetorSetNodeObject(modelNode.VapusID, modelNode, true)
	return
}

func (v *AIModelIntAgent) configureAIModelNode(ctx context.Context) error {
	node := utils.AIMNPb2Obj(v.managerRequest.GetSpec())
	if node == nil {
		v.Logger.Error().Msg("error while unmarshalling ai model node from managerRequest")
		return dmerrors.DMError(apperr.ErrInvalidAIModelNodeRequestSpec, nil)
	}
	var err error
	node.SetAINodeId()
	if len(node.Editors) == 0 {
		node.Editors = []string{v.CtxClaim[encryption.ClaimUserIdKey]}
	} else {
		node.Editors = append(node.Editors, v.CtxClaim[encryption.ClaimUserIdKey])
	}
	node.PreSaveCreate(v.CtxClaim)
	node.Status = mpb.CommonStatus_ACTIVE.String()
	if node.Scope == mpb.ResourceScope_ORGANIZATION_SCOPE.String() {
		if len(node.ApprovedOrganizations) == 0 {
			node.ApprovedOrganizations = []string{v.CtxClaim[encryption.ClaimOrganizationKey]}
		} else {
			node.ApprovedOrganizations = append(node.ApprovedOrganizations, v.CtxClaim[encryption.ClaimOrganizationKey])
		}
	}
	if node.ApprovedOrganizations == nil {
		node.ApprovedOrganizations = []string{v.CtxClaim[encryption.ClaimOrganizationKey]}
	} else {
		node.ApprovedOrganizations = append(node.ApprovedOrganizations, v.CtxClaim[encryption.ClaimOrganizationKey])
	}
	node.Organization = v.CtxClaim[encryption.ClaimOrganizationKey]
	newCtx := context.Background()

	if node.NetworkParams.SecretName == "" {
		secName := dmutils.GetSecretName("aistudio", node.VapusID, "aiModelNode")
		// Secret Name fetched
		err = apppkgss.SaveCredentialsCreds(ctx, secName, node.NetworkParams.Credentials, v.dmStore.VapusStore, v.Logger)
		node.NetworkParams.SecretName = secName
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while setting ai model node credentials")
			return dmerrors.DMError(apperr.ErrSettingAIModelNodeCredentials, err)
		}
	} else {
		secrets := &models.GenericCredentialModel{}
		secretStr, err := v.dmStore.VapusStore.SecretStore.ReadSecret(ctx, node.NetworkParams.SecretName)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while fetching the secrets")
			return dmerrors.DMError(apperr.ErrVapusSecret404, err)
		}
		err = json.Unmarshal([]byte(dmutils.AnyToStr(secretStr)), secrets)
		if err != nil {
			logger.Err(err).Ctx(ctx).Msg("error while unmarshelling creds from secret store.")
			return dmerrors.DMError(apperr.ErrDataSourceCredsSecretGet, err)
		}

		node.NetworkParams.Credentials = secrets
	}

	if node.DiscoverModels {
		err = crawlAIModels(newCtx, node, v.Logger) // nolint:errcheck,gosec //
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while crawling ai models")
			return dmerrors.DMError(apperr.ErrCrawlingAIModels, err)
		}
	}

	node.NetworkParams.Credentials = nil

	err = v.dmStore.ConfigureGetAIModelNode(ctx, node, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while saving ai model node metadata in datastore")
		return dmerrors.DMError(apperr.ErrUpdatingAIModelNode, err)
	}
	v.updateCachePool(node, true, false)
	v.SetCreateResponse(mpb.Resources_AIMODELS, node.VapusID)
	return nil
}

func (v *AIModelIntAgent) syncAIModelNode(ctx context.Context) error {
	if v.getterRequest == nil {
		v.Logger.Error().Msg("error while unmarshalling ai model nodes from managerRequest")
		return dmerrors.DMError(apperr.ErrInvalidAIModelNodeRequestSpec, nil)
	}
	// fetching Model Details
	nodeObj, err := v.dmStore.GetAIModelNode(ctx, v.getterRequest.GetAiModelNodeId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model node from datastore")
		return dmerrors.DMError(apperr.ErrAIModelNode404, err)
	}
	// Checking if the users has the editing rights
	if !slices.Contains(nodeObj.Editors, v.CtxClaim[encryption.ClaimUserIdKey]) {
		v.Logger.Error().Msg("error while validating user access")
		return dmerrors.DMError(apperr.ErrAIModelNode403, nil)
	}

	nodeObj.PreSaveUpdate(v.CtxClaim[encryption.ClaimUserIdKey])

	// To get the Network params through Secret Store
	netParams, err := v.dmStore.GetAIModelNodeNetworkParams(ctx, nodeObj)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model node network params")
		return dmerrors.DMError(apperr.ErrGetAIModelNetParams, err)
	}
	nodeObj.NetworkParams = netParams

	err = crawlAIModels(ctx, nodeObj, v.Logger) // nolint:errcheck,gosec //
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while crawling ai models")
		return dmerrors.DMError(apperr.ErrCrawlingAIModels, err)
	}
	nodeObj.NetworkParams.Credentials = nil
	err = v.dmStore.PutAIModelNode(ctx, nodeObj)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while saving ai model node metadata in datastore")
		return dmerrors.DMError(apperr.ErrSavingAIModelNode, err)
	}
	v.result.Output.AiModelNodes = utils.AIMNArrM2Pb([]*models.AIModelNode{nodeObj})
	return nil
}

func (v *AIModelIntAgent) updateAIModelNode(ctx context.Context) error {

	if v.managerRequest.GetSpec() == nil {
		v.Logger.Error().Msg("error while unmarshalling ai model nodes from updateAIModelNode")
		return dmerrors.DMError(apperr.ErrInvalidAIModelNodeRequestSpec, nil)
	}
	node := v.managerRequest.GetSpec()
	nodeObj, err := v.dmStore.GetAIModelNode(ctx, node.ModelNodeId, v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model node from datastore")
		return dmerrors.DMError(apperr.ErrAIModelNode404, err)
	}
	if nodeObj.Status != mpb.CommonStatus_ACTIVE.String() || nodeObj.DeletedAt != 0 {
		v.Logger.Error().Msg("ai model node is not active")
		return dmerrors.DMError(apperr.ErrResourceNotActive, nil)
	}
	if !slices.Contains(nodeObj.Editors, v.CtxClaim[encryption.ClaimUserIdKey]) {
		v.Logger.Error().Msg("error while validating user access")
		return dmerrors.DMError(apperr.ErrAIModelNode403, nil)
	}
	if nodeObj.Organization != v.CtxClaim[encryption.ClaimOrganizationKey] {
		v.Logger.Error().Msg("error while validating domain access")
		return dmerrors.DMError(apperr.ErrAIModelNode403, nil)
	}
	mConvertor := func(m []*mpb.AIModelBase) []*models.AIModelBase {
		r := make([]*models.AIModelBase, 0)
		if m != nil {
			for _, gm := range m {
				r = append(r, (&models.AIModelBase{}).ConvertFromPb(gm))
			}
		}
		return r
	}
	nodeObj.PreSaveUpdate(v.CtxClaim[encryption.ClaimUserIdKey])
	if nodeObj.Scope == mpb.ResourceScope_ORGANIZATION_SCOPE.String() {

		nodeObj.ApprovedOrganizations = node.Attributes.GetApprovedOrganizations()
	}
	nodeObj.SecurityGuardrails = (&models.SecurityGuardrails{}).ConvertFromPb(node.GetSecurityGuardrails())
	nodeObj.GenerativeModels = mConvertor(node.Attributes.GetGenerativeModels())
	nodeObj.EmbeddingModels = mConvertor(node.Attributes.GetEmbeddingModels())
	if node.Attributes.GetNetworkParams().Credentials != nil {
		secname := node.Attributes.GetNetworkParams().SecretName
		if secname == "" || secname == nodeObj.NetworkParams.SecretName {
			secname = dmutils.GetSecretName("aistudio", "", dmutils.GetStrEpochTime())
		}
		creds := (&models.GenericCredentialModel{}).ConvertFromPb(node.Attributes.GetNetworkParams().Credentials)
		err = apppkgss.SaveCredentialsCreds(ctx, secname, creds, v.dmStore.VapusStore, v.Logger)
		if err != nil {
			v.Logger.Error().Err(err).Msg("error while updating ai model node credentials")
			return dmerrors.DMError(apperr.ErrSettingAIModelNodeCredentials, err)
		}
		nodeObj.NetworkParams.SecretName = secname
	} else if node.Attributes.GetNetworkParams().SecretName == "" {
		nodeObj.NetworkParams.SecretName = ""
	}
	nodeObj.NetworkParams.Credentials = nil
	err = v.dmStore.PutAIModelNode(ctx, nodeObj)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while updating ai model node metadata in datastore")
		return dmerrors.DMError(apperr.ErrUpdatingAIModelNode, err)
	}
	v.updateCachePool(nodeObj, false, true)
	v.result.Output.AiModelNodes = utils.AIMNArrM2Pb([]*models.AIModelNode{nodeObj})
	return nil
}

func (v *AIModelIntAgent) describeAIModelNode(ctx context.Context) error {
	result, err := v.dmStore.GetAIModelNode(ctx, v.getterRequest.GetAiModelNodeId(), v.CtxClaim)
	if err != nil || result == nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model node describe from datastore")
		return dmerrors.DMError(err, nil)
	}

	v.result.Output.AiModelNodes = utils.AIMNArrM2Pb([]*models.AIModelNode{result})
	return nil
}

func (v *AIModelIntAgent) archiveAIModelNode(ctx context.Context) error {
	if v.getterRequest.AiModelNodeId == "" {
		v.Logger.Error().Msg("error: invalid ai prompt id")
		return dmerrors.DMError(apperr.ErrAIPrompt404, nil)
	}
	// TODO: Track the failed nodes and return the failed nodes in the response of DMresponse.
	result, err := v.dmStore.GetAIModelNode(ctx, v.getterRequest.GetAiModelNodeId(), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model node describe from datastore")
		return dmerrors.DMError(err, nil)
	}
	result.Status = mpb.CommonStatus_DELETED.String()
	result.PreSaveDelete(v.CtxClaim)
	err = v.dmStore.PutAIModelNode(ctx, result)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while archiving ai model node")
		return dmerrors.DMError(err, nil)
	}
	pkgs.AIModelNodeConnectionPoolManager.RemoveConnection(result)
	pkgs.AIModelNodeConnectionPoolManager.RemoveNodeObject(result.VapusID)
	v.result.Output.AiModelNodes = []*mpb.AIModelNode{}
	return nil
}

func (v *AIModelIntAgent) listAIModelNodes(ctx context.Context) error {
	// condition := fmt.Sprintf("(status = 'ACTIVE' AND deleted_at IS NULL AND scope = 'PLATFORM_SCOPE' OR (scope = 'DOMAIN_SCOPE' AND '%s' = ANY(approved_domains))) OR domain='%s' ORDER BY created_at DESC", v.CtxClaim[encryption.ClaimOrganizationKey], v.CtxClaim[encryption.ClaimOrganizationKey])
	result, err := v.dmStore.ListAIModelNodes(ctx, apppkgss.ListResourceWithGovernance(v.CtxClaim), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model node list from datastore")
		return dmerrors.DMError(err, nil)
	}
	v.result.Output.AiModelNodes = utils.AIMNList2Pb(result)
	return nil
}

// Insights
func (v *AIModelIntAgent) listUsageInsights(ctx context.Context) error {
	// condition := fmt.Sprintf("(status = 'ACTIVE' AND deleted_at IS NULL AND scope = 'PLATFORM_SCOPE' OR (scope = 'DOMAIN_SCOPE' AND '%s' = ANY(approved_domains))) OR domain='%s' ORDER BY created_at DESC", v.CtxClaim[encryption.ClaimOrganizationKey], v.CtxClaim[encryption.ClaimOrganizationKey])
	result, err := v.dmStore.GetAIModelNodeId(ctx, apppkgss.ListResourceWithGovernance(v.CtxClaim), v.CtxClaim)
	if err != nil {
		v.Logger.Error().Err(err).Msg("error while getting ai model node list from datastore")
		return dmerrors.DMError(err, nil)
	}

	clientNodeId := v.usageInsightRequest.GetAiModelNodeId()
	clientModelName := v.usageInsightRequest.GetModel()

	if len(clientNodeId) > 0 {
		set := make(map[string]bool)
		for _, val := range clientNodeId {
			set[val] = true
		}
		modelNodeObservabilities := []*mpb.ModelNodeObservability{}
		for _, val := range result {
			if set[*val] {
				modelNodeObservability, err := v.dmStore.GetAIModelUsageInsights(ctx, apppkgss.ListResourceWithGovernance(v.CtxClaim), val, clientModelName, v.CtxClaim)
				if err != nil {
					v.Logger.Error().Err(err).Msg("error while getting ai model node list from datastore")
					return dmerrors.DMError(err, nil)
				}
				if modelNodeObservability != nil {
					modelNodeObservabilities = append(modelNodeObservabilities, modelNodeObservability)
				}
			}
		}
		v.insightResult.ModelNodeObservability = modelNodeObservabilities
		return nil

	} else {
		modelNodeObservabilities := []*mpb.ModelNodeObservability{}
		for _, modelNode := range result {
			modelNodeObservability, err := v.dmStore.GetAIModelUsageInsights(ctx, apppkgss.ListResourceWithGovernance(v.CtxClaim), modelNode, clientModelName, v.CtxClaim)
			if err != nil {
				v.Logger.Error().Err(err).Msg("error while getting ai model node list from datastore")
				return dmerrors.DMError(err, nil)
			}
			if modelNodeObservability != nil {
				modelNodeObservabilities = append(modelNodeObservabilities, modelNodeObservability)
			}
		}
		v.insightResult.ModelNodeObservability = modelNodeObservabilities
		return nil
	}
}
