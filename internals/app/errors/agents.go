package apperr

import "errors"

var (
	ErrInRecommendationByAgent               = errors.New("error while processing chart builder agent")
	EmailAgentActionFailed                   = errors.New("error while processing email agent")
	ErrCommunicationAgentFailed              = errors.New("error while processing communication agent")
	ErrCommunicationChannelNotSupported      = errors.New("error while processing communication agent, channel not supported")
	ErrNotAbleToDetermineEmailAttachmentType = errors.New("error while processing email agent, not able to determine email attachment type")
	ErrDataAnalysisSuggestionFailed          = errors.New("error while suggesting data analysis for given data")
	ErrAttachedDatasetEmpty                  = errors.New("error while processing data, attached dataset is empty")
	ErrTargetDatasetEmpty                    = errors.New("error while processing data, target dataset is empty")
	ErrInvalidToolFunctionCall2Action        = errors.New("error while processing tool function call to action")
	ErrFailedGoogleDriverOperatorAction      = errors.New("error while processing google drive operator action")
	ErrTaskmanagerAnalysisFailed             = errors.New("error while processing task manager agent")
	ErrTaskmanagerActionFailed               = errors.New("error while processing task manager agent")
	ErrPythonProjectInitFailed               = errors.New("error while initializing python project")
	ErrPythonScriptAddFailed                 = errors.New("error while adding python script")
)
