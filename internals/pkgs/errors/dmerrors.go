package dmerrors

import "errors"

var (
	// Error constants for JSON operations
	ErrJsonMarshel   = errors.New("error while marshalling JSON")
	ErrJsonUnMarshel = errors.New("error while unmarshalling JSON")
	ErrStruct2Json   = errors.New("failed to convert struct to json")

	// Error constants for viper operations
	ErrViperConfigRead = errors.New("error while reading configuration file")
	ErrViperConfigSet  = errors.New("error while setting configuration file")

	ErrUserORGANIZATION404            = errors.New("error- invalid ORGANIZATION requested, user is not attached to requested ORGANIZATION")
	ErrWriteYAMLFile                  = errors.New("error while writing to yaml file")
	ErrInvalidArgs                    = errors.New("invalid arguments provided for the command, please provide the required arguments")
	ErrNoCredentialFoundForDataSource = errors.New("no credentials found for the data source")
	ErrInvalidComplianceAgentParams   = errors.New("invalid compliance agent parameters")
	ErrInvalidComplianceAction        = errors.New("invalid compliance action")
	ErrHomeDirNotFound                = errors.New("home directory not found")
	ErrNoDataMessage                  = errors.New("no data found in file")
	ErrQueryProcessingFailed          = errors.New("error while processing query")
	ErrDataPrepFailed                 = errors.New("error while preparing data")
	ErrAgentReasoingFailed            = errors.New("error while reasoning agent")
	ErrStreamEnded                    = errors.New("error while sending stream data")
	ErrInvalidQuery                   = errors.New("invalid query provided")
	ErrUser404                        = errors.New("error- user not found")
	ErrDataProductAccess403           = errors.New("error- data product access forbidden")
	ErrDataProductAccess404           = errors.New("error- data product not found")
	ErrQueryProcessingFailedTryAgain  = errors.New("error while processing query, please try again")
	ErrFileFormatConversionFailed     = errors.New("error while converting file format")
	ErrStorageInternalError           = errors.New("error while storing data, internal error")
	ErrSendingEmailTryAgain           = errors.New("error while sending email, please try again")
	ErrDataProduct404                 = errors.New("error- data product not found")
	ErrInternalError                  = errors.New("internal error")
	ErrSSeChannelFull                 = errors.New("error- SSE channel full")
)
