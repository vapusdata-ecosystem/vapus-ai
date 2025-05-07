package authn

import (
	"strings"

	"github.com/rs/zerolog"
)

type AuthnSub string

var (
	GOOGLE   AuthnSub = "google"
	LINKEDIN AuthnSub = "linkedin"
	GITHUB   AuthnSub = "github"
)

func (a AuthnSub) String() string {
	return string(a)
}

func GetAuthnConnectionSub(claim map[string]any, logger zerolog.Logger) (AuthnSub, error) {
	sub, ok := claim["sub"]
	if !ok {
		logger.Error().Msgf("error: sub not found in authenticated idtoken")
		return "", ErrInvalidUserAuthenticationSub
	}
	subj, ok := sub.(string)
	if !ok {
		logger.Error().Msgf("error: sub not found in authenticated idtoken")
		return "", ErrInvalidUserAuthenticationSub
	}
	subs := strings.Split(subj, "|")
	if len(subs) != 2 {
		logger.Error().Msgf("error: sub is not in the correct format")
		return "", ErrInvalidUserAuthenticationSub
	}
	if strings.Contains(subs[0], "google") {
		return GOOGLE, nil
	} else if strings.Contains(subs[0], "linkedin") {
		return LINKEDIN, nil
	} else if strings.Contains(subs[0], "github") {
		return GITHUB, nil
	}
	logger.Error().Msgf("error: sub doesn't match any of the supported authn connections")
	return "", ErrInvalidUserAuthenticationSub
}

func GetEmail(claim map[string]any) (string, error) {
	email, ok := claim["email"]
	if !ok {
		return "", ErrClaimEmailNotFound
	}
	return email.(string), nil
}

func (a AuthnSub) GetFirstName(claim map[string]any) (string, error) {
	switch a {
	case GOOGLE:
		val, ok := claim["given_name"]
		if !ok {
			return "", ErrClaimFirstNameNotFound
		}
		return val.(string), nil
	case GITHUB:
		val, ok := claim["name"]
		if !ok {
			return "", ErrClaimFirstNameNotFound
		}
		return val.(string), nil
	default:
		return "", ErrClaimFirstNameNotFound
	}
}

func (a AuthnSub) GetDisplayName(claim map[string]any) (string, error) {
	switch a {
	case GOOGLE:
		val, ok := claim["name"]
		if !ok {
			return "", ErrClaimDisplayNameNotFound
		}
		return val.(string), nil
	case GITHUB:
		val, ok := claim["name"]
		if !ok {
			return "", ErrClaimDisplayNameNotFound
		}
		return val.(string), nil
	default:
		return "", ErrClaimDisplayNameNotFound
	}
}

func (a AuthnSub) GetLastName(claim map[string]any) (string, error) {
	switch a {
	case GOOGLE:
		val, ok := claim["family_name"]
		if !ok {
			return "", ErrClaimLastNameNotFound
		}
		return val.(string), nil
	case GITHUB:
		return "", nil
	default:
		return "", ErrClaimLastNameNotFound
	}
}

func (a AuthnSub) GetProfilePic(claim map[string]any) (string, error) {
	switch a {
	case GOOGLE:
		val, ok := claim["picture"]
		if !ok {
			return "", ErrClaimLastNameNotFound
		}
		return val.(string), nil
	case GITHUB:
		val, ok := claim["picture"]
		if !ok {
			return "", ErrClaimLastNameNotFound
		}
		return val.(string), nil
	default:
		return "", ErrClaimLastNameNotFound
	}
}
