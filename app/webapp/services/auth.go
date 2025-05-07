package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/webapp/pkgs"
	routes "github.com/vapusdata-ecosystem/vapusdata/webapp/routes"
)

type Authentication interface {
	LoginHandler(c echo.Context) error
	LoginRedirect(c echo.Context) error
	CallbackHandler(c echo.Context) error
	Logout(c echo.Context) error
}

type AuthnService struct {
	logger zerolog.Logger
}

var AuthnServiceManager Authentication

func newAuthnService(configPath string) Authentication {
	l := pkgs.GetSubDMLogger("Authentication", dmutils.GetUUID())
	routes.Login = pkgs.WebAppConfigManager.URIs.Login
	routes.LoginCallBack = pkgs.WebAppConfigManager.URIs.Callback
	routes.Logout = pkgs.WebAppConfigManager.URIs.Logout
	return &AuthnService{
		logger: l,
	}
}

func InitAuthnService(configPath string) {
	if AuthnServiceManager == nil {
		AuthnServiceManager = newAuthnService(configPath)
	}
	// Set the Auth and Callback in the AuthnServiceManager instance
	// AuthnServiceManager.SetAuths()
}

func (au *AuthnService) LoginHandler(c echo.Context) error {
	publicInfo, err := pkgs.VapusSvcInternalClientManager.PlConn.PlatformPublicInfo(c.Request().Context(), &mpb.EmptyRequest{})
	if err != nil {
		au.logger.Err(err).Msg("error while getting platform public info")
	}
	return c.Render(http.StatusOK, "login.html", map[string]any{
		"platform":      "VapusData Platform",
		"loginRedirect": routes.LoginRedirect,
		"publicInfo":    publicInfo,
		"registerAPI":   fmt.Sprintf("%s/api/v1alpha1/register", pkgs.NetworkConfigManager.GatewayURL),
		"landingPage":   routes.UIHome,
	})
}

func (au *AuthnService) LoginRedirect(c echo.Context) error {
	result, err := pkgs.VapusSvcInternalClientManager.UserConn.LoginHandler(c.Request().Context(), &mpb.EmptyRequest{})
	if err != nil {
		au.logger.Err(err).Msg("error while getting login url")
		c.Redirect(http.StatusTemporaryRedirect, pkgs.WebAppConfigManager.URIs.Login)
	}
	au.logger.Info().Msgf("HOST -- %v", c.Request().Host)
	au.logger.Info().Msgf("Login URL -- %v", result.LoginUrl)
	c.Redirect(http.StatusTemporaryRedirect, result.LoginUrl)
	return nil
}

func (au *AuthnService) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = types.ACCESS_TOKEN
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Expires = time.Unix(0, 0)
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, routes.Login)
}

func (au *AuthnService) CallbackHandler(c echo.Context) error {
	// Get the code from the query parameter
	code := c.QueryParam("code")

	au.logger.Info().Msgf("Callback received with code: %v", code)
	result, err := pkgs.VapusSvcInternalClientManager.UserConn.LoginCallback(c.Request().Context(), &pb.LoginCallBackRequest{
		Code: code,
		Host: pkgs.NetworkConfigManager.ExternalURL + routes.LoginCallBack,
	})
	if err != nil || result == nil {
		au.logger.Err(err).Msg("error login callback")
		c.Redirect(http.StatusTemporaryRedirect, pkgs.WebAppConfigManager.URIs.Login)
		return nil
	}

	idTCookie := &http.Cookie{
		Name:    pkgs.ID_TOKEN,
		Value:   result.Token.IdToken,
		Expires: time.Unix(result.Token.ValidTill, 0),
		Path:    "/",
	}
	c.SetCookie(idTCookie)
	if result.Token.AccessToken == "" {
		c.Redirect(http.StatusTemporaryRedirect, pkgs.WebAppConfigManager.URIs.Login+"?register=true")
	}
	// // Set the user profile in the cookie
	atCookie := &http.Cookie{
		Name:    types.ACCESS_TOKEN,
		Value:   result.Token.AccessToken,
		Expires: time.Unix(result.Token.ValidTill, 0),
		Path:    "/",
	}
	c.SetCookie(atCookie)
	time.Sleep(2 * time.Second)

	return c.Redirect(http.StatusTemporaryRedirect, routes.UIHome)
}
