package apperr

import (
	"errors"
)

var (
	ErrPlugin404                  = errors.New("error while fetching plugin, plugin not found")
	ErrPluginPatch400             = errors.New("error while patching plugin")
	ErrPluginNotEditable          = errors.New("plugin is not editable")
	ErrPluginCreate400            = errors.New("error while configuring plugin")
	ErrPlugin403                  = errors.New("error while accessing plugin, forbidden access for the user")
	ErrPluginServiceScopeExists   = errors.New("error while configuring plugin, service scope already exists")
	ErrInvalidPluginCallerAction  = errors.New("invalid action for plugin caller")
	ErrInvalidPluginScopeParams   = errors.New("error while configuring plugin, invalid plugin scope params")
	ErrPluginPlatformScope403     = errors.New("error while configuring plugin with platform scope, forbidden access for the user. Ask platform owner to configure the plugin")
	ErrPluginOrganizationScope403 = errors.New("error while configuring plugin with Organization scope, forbidden access for the user. Ask Organization owner to configure the plugin")
	ErrUnSupportedPluginType      = errors.New("error while configuring plugin, unsupported plugin type")
	ErrPluginScope403             = errors.New("error while configuring plugin, forbidden access for the user")
	ErrInvalidPluginTypeForAction = errors.New("error while configuring plugin, invalid plugin type for the action")
	ErrFileStorePlugin404         = errors.New("error while fetching file store plugin, file store plugin not found")
	ErrFileStorePlugin400         = errors.New("error while initializing file store plugin")

	ErrPluginNotConnected         = errors.New("error while connecting plugin, plugin not connected")
	ErrInvalidPluginScope         = errors.New("error while configuring plugin, invalid plugin scope requested")
	ErrInvalidPluginObject        = errors.New("error while configuring plugin, invalid plugin object requested")
	ErrInvalidPluginCredentials   = errors.New("error while configuring plugin, invalid plugin credentials requested")
	ErrInvalidEmailerEngine       = errors.New("invalid emailer engine")
	ErrEmailerConn                = errors.New("error while connecting to emailer engine")
	ErrInvalidEmailerCredentials  = errors.New("invalid emailer credentials")
	ErrInvalidPluginService       = errors.New("error while configuring plugin, invalid plugin service requested")
	ErrInvalidPluginType          = errors.New("error while configuring plugin, invalid plugin type requested")
	ErrInvalidFileSharingType     = errors.New("error while configuring plugin, invalid file sharing type requested")
	ErrInvalidFileSharingResource = errors.New("error while configuring plugin, invalid file sharing resource requested")
	ErrInvalidCalenderService     = errors.New("invalid calendar credentials")
	ErrInvalidCalenderConn        = errors.New("error while configuring plugin, invalid plugin service requested")
)
