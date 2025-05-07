package policies

import (
	"github.com/casbin/casbin/v2"
	casbinmodel "github.com/casbin/casbin/v2/model"
	"github.com/rs/zerolog"
)

var DataProductAuthzModel = `
[request_definition]
r = sub, ORGANIZATION, action

[policy_definition]
p = sub, ORGANIZATION, action

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.Organization == p.Organization && r.action == p.action
`

var VapusDataPlatformAuthzModel = `
[request_definition]
# Defines the input fields for a request
r = resourceArn, action, userId, ORGANIZATION

[policy_definition]
# Defines the fields in a policy
p = resourceArn, action, userId, ORGANIZATION, effect

[policy_effect]
# Defines how to handle multiple matching policies
e = some(where (p.eft == allow))

[matchers]
# Defines the matching logic for request and policy
m = r.resourceArn == p.resourceArn && r.action == p.action && r.userId == p.userId && r.Organization == p.Organization
`

func GetDataProductAuthz(logger zerolog.Logger) (*casbin.Enforcer, error) {
	model, err := casbinmodel.NewModelFromString(DataProductAuthzModel)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating model for data product server authz")
		return nil, err
	}
	enforcer, err := casbin.NewEnforcer(model)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating enforcer for data product server authz")
		return nil, err
	}
	return enforcer, nil
}

func GetVapusResourceAuthz(logger zerolog.Logger) (*casbin.Enforcer, error) {
	model, err := casbinmodel.NewModelFromString(VapusDataPlatformAuthzModel)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating model for vapus resource authz model")
		return nil, err
	}
	enforcer, err := casbin.NewEnforcer(model)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating enforcer for vapus resource authz model")
		return nil, err
	}
	return enforcer, nil
}
