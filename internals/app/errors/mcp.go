package apperr

import "errors"

var (
	ErrMcpToolNotFound                  = errors.New("error while processing mcp tool, tool not found")
	ErrMcpToolNotSupported              = errors.New("error while processing mcp tool, tool not supported")
	ErrMcpToolNotAvailable              = errors.New("error while processing mcp tool, tool not available")
	ErrMcpToolNotAuthorized             = errors.New("error while processing mcp tool, tool not authorized")
	ErrMcpToolNotConfigured             = errors.New("error while processing mcp tool, tool not configured")
	ErrMcpToolNotInitialized            = errors.New("error while processing mcp tool, tool not initialized")
	ErrMcpToolNotReady                  = errors.New("error while processing mcp tool, tool not ready")
	ErrPythonCode404                    = errors.New("error while processing python code, code not found")
	ErrPythonScriptRunFailed            = errors.New("error while running python script, script run failed")
	ErrPythonCodingActionNotFound       = errors.New("error while processing python code, action not found")
	ErrFileOperationFailed              = errors.New("error while processing file operation, operation failed")
	ErrFileOperatorToolDiscoveryFailed  = errors.New("error while discovering file operator tool, discovery failed")
	ErrFileOperatorPlanningFailed       = errors.New("error while planning file operator tool, planning failed")
	ErrCommunicationFailed              = errors.New("error while processing communication, communication failed")
	ErrCommunicationToolDiscoveryFailed = errors.New("error while discovering communication tool, discovery failed")
	ErrCommunicationPlanningFailed      = errors.New("error while planning communication tool, planning failed")
	ErrDataAnalysisPlanningFailed       = errors.New("error while planning data analysis, planning failed")
	ErrDataAnalysisCodingFailed         = errors.New("error while processing data analysis, coding failed")
	ErrGSTCalculationFailed             = errors.New("error while processing GST calculation, calculation failed")
	ErrGSTFieldMappingFailed            = errors.New("error while processing GST field mapping, mapping failed")
)
