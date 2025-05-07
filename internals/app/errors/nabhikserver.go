package apperr

import (
	"errors"
)

var (
	ErrInvalidDataProductConfig     = errors.New("error: Invalid Data Product Config")
	ErrInvalidDataProductFileFormat = errors.New("invalid file format requested for the Data Product Config")

	ErrInvalidManageAction = errors.New("error - invalid manage action")

	ErrInvalidUserAuthentication = errors.New("error while validating user's authentication")

	ErrDataProductNotExists     = errors.New("no data product exists")
	ErrInvalidDataProductAction = errors.New("invalid action for data product")

	ErrDataProductAccess403   = errors.New("error while querying data product, user does not have access to the data product")
	ErrDataProductQueryFailed = errors.New("error while querying data product")
	ErrInvalidQuery           = errors.New("error while querying data product, invalid query provided")

	ErrInvalidDataProductSpec   = errors.New("error while loading data product spec")
	ErrQueryProcessingFailed    = errors.New("error while processing the query,please try again or check the query or check the context of your previous messages")
	ErrInvalidQueryParams       = errors.New("error while processing the query, invalid query params provided")
	ErrInvalidROQuery           = errors.New("error while processing the query, invalid read-only query provided")
	ErrDataProductAccess404     = errors.New("error while querying data product, data product not found")
	ErrInvalidTablesInQuery     = errors.New("error while processing the query, invalid tables in the query")
	ErrCuurentlySupportOnetable = errors.New("error while processing the query, currently only one table is supported")
	ErrQueryBuildingFailed      = errors.New("error while building the query")

	ErrInvalidTablesInQuery403      = errors.New("error while processing the query, user does not have access to the tables in the query")
	ErrDataProduct404SelectManually = errors.New("error: data product not found, select manually from your list or search data product and request access for the same.")
	InvalidDataQueryInput           = errors.New("error: invalid data query input, please provide valid input with more information or pass the dataproduct parameter")

	ErrNoDataMessage = errors.New("error while processing the query, no data message found")

	ErrInvalidChartTypeRequsted  = errors.New("error while processing the query, invalid chart type requested")
	ErrChartBuildingFailedParams = errors.New("error while building chart, invalid chart params provided")
	ErrInvalidInputGiveMoreInfo  = errors.New("error while analyzing the input, be more specific with meaningful words and provide more information")

	ErrSendingEmailTryAgain          = errors.New("error while sending email, please try again")
	ErrStorageInternalError          = errors.New("error while storing data, internal error")
	ErrQueryProcessingFailedTryAgain = errors.New("error while processing the query, please try again")

	ErrDataProcessingFailed = errors.New("error while processing the data, please try again")
	ErrDataFilteringFailed  = errors.New("error while filtering the data, please try again")
	ErrDataMarshalingFailed = errors.New("error while marshalling the data, please try again")
	ErrDataFileLoadFailed   = errors.New("error while loading the data file, please try again")
	ErrStreamEnded          = errors.New("error while processing the stream, stream ended")

	ErrGettingDataProductConsumptionMetrics = errors.New("error while getting data product consumption metrics")
	ErrInvalidDataUploadAction              = errors.New("error while processing the query, invalid data upload action")
	ErrFormatConversionFailed               = errors.New("error while converting the format, please try again")

	ErrDataproductQueryAgentInitFailed    = errors.New("error while initializing data product query agent")
	ErrDataproductQueryAgentExecuteFailed = errors.New("error while executing data product query agent")
	ErrInvalidVapusAgenttRequested        = errors.New("error while processing the query, invalid fabric agent requested")
	ErrVapusAgentInitFailed               = errors.New("error while initializing fabric agent")
	ErrVapusAgentReadyFailed              = errors.New("error while checking if fabric agent is ready")
	ErrUpdatingVapusAgent                 = errors.New("error while updating fabric agent")
	ErrInvalidVapusAgent404               = errors.New("error while processing the query, invalid fabric agent 404")
	ErrInvalidVapusAgent403               = errors.New("error while processing the query, invalid fabric agent 403")
	ErrInvalidVapusAgentRequested         = errors.New("error while processing the query, invalid fabric agent requested")
	ErrCreatingVapusAgent                 = errors.New("error while creating fabric agent")
	ErrArchivingAgentFailed               = errors.New("error while archiving agent")
	ErrVapusAgent403                      = errors.New("error while processing the query, agent 403")
	ErrVapusAgentOwner403                 = errors.New("error while processing the query, agent owner 403")
	ErrInvalidAgentServiceAction          = errors.New("error while processing the query, invalid agent service action")
)
