package appcl

import (
	"errors"
)

var (
	ErrAIStudioConnNotInitialized          = errors.New("AIStudio connection not initialized")
	ErrUserConnNotInitialized              = errors.New("User connection not initialized")
	ErrDataProductServerConnNotInitialized = errors.New("DataProductServer connection not initialized")
	ErrPlatConnNotInitialized              = errors.New("Platform connection not initialized")
	Erruser404                             = errors.New("User not found")
	ErrNoDataFound                         = errors.New("No data found")
	ErrPluginActionFailed                  = errors.New("Plugin action failed, this plugin is not available for your account, enable it from the plugin store under settings section")
	ErrDataProductQueryCallFailed          = errors.New("Data product query call failed")
	ErrGeneratingContent                   = errors.New("Error generating content")
	ErrInvalidPluginActionRequest          = errors.New("Invalid plugin action request")
	ErrDataProducts404                     = errors.New("Data products not found, please try again later")
	ErrDatamarketplaceConnNotInitialized   = errors.New("DataMarketplace connection not initialized")
	ErrAIUtilityConnNotInitialized          = errors.New("AIUtility connection not initialized")
)
