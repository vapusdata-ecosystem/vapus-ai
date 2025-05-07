package services

import (
	"slices"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/vapusdata-ecosystem/vapusai/core/pkgs/pbtools"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
	"github.com/vapusdata-ecosystem/vapusai/webapp/models"
	routes "github.com/vapusdata-ecosystem/vapusai/webapp/routes"
	utils "github.com/vapusdata-ecosystem/vapusai/webapp/utils"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

var OrganizationMap map[string]string

func (x *WebappService) BaseGroupClientCalls(c echo.Context, obj *models.GlobalContexts) error {
	errCh := make(chan error, 3)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		currentOrganization, err := x.grpcClients.GetCurrentOrganization(c)

		if err != nil {
			errCh <- err
			return
		}
		obj.CurrentOrganization = currentOrganization
		errCh <- nil
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result, err := x.grpcClients.GetAccountInfo(c)
		if err != nil {
			errCh <- err
			return
		}
		obj.Account = result.Output
		errCh <- nil
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		res, domainMap, err := x.grpcClients.GetUserInfo(c, "")

		if err != nil {
			errCh <- err
			return
		}
		obj.UserInfo = res
		obj.OrganizationMap = domainMap
		errCh <- nil
	}()

	wg.Wait()
	close(errCh)
	if err := <-errCh; err != nil {
		// If any goroutine returns an error, log it and return the error
		x.logger.Err(err).Msg("error while getting data from server")
		return err
	}
	for _, ud := range obj.UserInfo.Roles {
		if ud.OrganizationId == obj.CurrentOrganization.OrganizationId {
			if slices.Contains(ud.Role, mpb.UserRoles_ORG_OWNER.String()) {
				obj.IsOrganizationOwner = true
			}
		}
	}
	obj.CurrentUrl = c.Request().URL.EscapedPath()
	// obj.OrganizationMap = x.LoadOrganizationMap(c)
	return nil
}

func (x *WebappService) getSectionGlobals(c echo.Context, nav, currentpage string) (*models.GlobalContexts, error) {
	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		BottomMenuMap:  routes.BottonMenuList,
		LoginUrl:       routes.Login,
		AccessTokenKey: types.ACCESS_TOKEN,
		CurrentNav:     nav,
		CurrentSideBar: currentpage,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) getDataMarketplaceSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {
	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		BottomMenuMap:  routes.BottonMenuList,
		CurrentSideBar: currentpage,
		CurrentNav:     routes.DataMarketplaceNav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: types.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) getHomeSectionGlobals(c echo.Context) (*models.GlobalContexts, error) {
	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		BottomMenuMap:  routes.BottonMenuList,
		CurrentNav:     routes.HomeNav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: types.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) getNabhikSectionGlobals(c echo.Context) (*models.GlobalContexts, error) {
	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		BottomMenuMap:  routes.BottonMenuList,
		CurrentNav:     routes.NabhikAINav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: types.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) getInsightsSectionGlobals(c echo.Context) (*models.GlobalContexts, error) {
	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		BottomMenuMap:  routes.BottonMenuList,
		CurrentNav:     routes.NabhikAINav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: types.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) getDatamanagerSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {

	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		BottomMenuMap:  routes.BottonMenuList,
		SidebarMap:     routes.DatamanagerSideList,
		CurrentSideBar: currentpage,
		CurrentNav:     routes.MyOrganizationNav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: types.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// func (x *WebappService) getExploreSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {

// 	obj := &models.GlobalContexts{
// 		NavMenuMap:     routes.NavMenuList,
// 		SidebarMap:     routes.ExploreSideList,
// 		CurrentSideBar: currentpage,
// 		CurrentNav:     routes.ExploreNav.String(),
// 		LoginUrl:       routes.Login,
// 		AccessTokenKey: pkgs.ACCESS_TOKEN,
// 	}
// 	err := x.BaseGroupClientCalls(c, obj)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return obj, nil
// }

func (x *WebappService) getAiStudioSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {

	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		BottomMenuMap:  routes.BottonMenuList,
		SidebarMap:     routes.ManageAINavSideList,
		CurrentSideBar: currentpage,
		CurrentNav:     routes.ManageAINav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: types.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) getStudioSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {

	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		BottomMenuMap:  routes.BottonMenuList,
		CurrentSideBar: currentpage,
		CurrentNav:     routes.VapusStudioNav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: types.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) getSettingsSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {
	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		BottomMenuMap:  routes.BottonMenuList,
		CurrentSideBar: currentpage,
		CurrentNav:     routes.SettingsNav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: types.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) getDeveloperSectionGlobals(c echo.Context, currentpage string) (*models.GlobalContexts, error) {
	obj := &models.GlobalContexts{
		NavMenuMap:     routes.NavMenuList,
		BottomMenuMap:  routes.BottonMenuList,
		CurrentSideBar: currentpage,
		CurrentNav:     routes.DevelopersNav.String(),
		LoginUrl:       routes.Login,
		AccessTokenKey: types.ACCESS_TOKEN,
	}
	err := x.BaseGroupClientCalls(c, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (x *WebappService) responseHandler(c echo.Context, status int, template string, data map[string]any) error {
	return c.Render(status, template, data)
}

func GetUrlParams(c echo.Context, resourceName string) (string, error) {
	_, ok := types.UrlResouceIdMap[resourceName]
	if !ok {
		return "", utils.ErrInvalidURLResourceName
	}
	return c.Param(resourceName), nil
}

func GetProtoYamlString(obj protoreflect.ProtoMessage) string {
	bytess, err := pbtools.ProtoYamlMarshal(obj)
	if err != nil {
		return ""
	}
	return string(bytess)
}

func GetUserCurrentOrganizationRole(user *mpb.User, currOrganization string) []string {
	for _, domain := range user.Roles {
		if domain.OrganizationId == currOrganization {
			return domain.Role
		}
	}
	return []string{}
}

func BuildBacklistingUrl(url string) string {
	urlOps := strings.Split(url, "/")
	if len(urlOps) > 1 {
		return strings.Join(urlOps[:len(urlOps)-1], "/")
	}
	return url
}
