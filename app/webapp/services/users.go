package services

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusdata/core/types"
	routes "github.com/vapusdata-ecosystem/vapusdata/webapp/routes"
	"github.com/vapusdata-ecosystem/vapusdata/webapp/utils"
)

func (x *WebappService) AuthOrganizationHandler(c echo.Context) error {
	at, err := utils.GetCookie(c, types.ACCESS_TOKEN)
	if err != nil || at == "" {
		return c.Redirect(http.StatusSeeOther, routes.Login)

	}
	redirectPath := ""
	val, ok := c.Request().URL.Query()["redirect"]
	if !ok {
		redirectPath = routes.UIRoute + routes.HomeGroup
	} else {
		redirectPath = val[0]
	}

	ctx := utils.GetBearerCtx(c.Request().Context(), at)
	domain := c.Param("domain")
	result, err := x.grpcClients.UserConn.AccessTokenInterface(ctx, &pb.AccessTokenInterfaceRequest{
		Organization: domain,
		Utility:      pb.AccessTokenAgentUtility_ORGANIZATION_LOGIN,
	})
	if err != nil {
		x.logger.Err(err).Msg("error while authenticating domain")
		return err
	}
	atCookie := &http.Cookie{
		Name:    types.ACCESS_TOKEN,
		Value:   result.Token.AccessToken,
		Expires: time.Unix(result.Token.ValidTill, 0),
		Path:    "/",
	}

	c.SetCookie(atCookie)
	time.Sleep(2 * time.Second)
	// Redirect to the home page
	return c.Redirect(http.StatusSeeOther, redirectPath)
}
