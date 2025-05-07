package policies

import (
	"github.com/casbin/casbin/v2"
	"github.com/rs/zerolog"
)

type VapusPlatformAuthz struct {
	Enforcer *casbin.Enforcer
	logger   zerolog.Logger
}

// In use
func NewVapusPlatformAuthzEnforcer(logger zerolog.Logger) (*VapusPlatformAuthz, error) {
	enforcer, err := GetVapusResourceAuthz(logger)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating enforcer")
		return nil, err
	}
	return &VapusPlatformAuthz{
		Enforcer: enforcer,
	}, nil
}

func (x *VapusPlatformAuthz) ValidatePolicy(resource string, ctxClaim map[string]string) bool {
	ok, err := x.Enforcer.Enforce(resource)
	if err != nil {
		x.logger.Error().Err(err).Msgf("Error enforcing policy for resource %v", resource)
		return false
	}
	return ok
}
