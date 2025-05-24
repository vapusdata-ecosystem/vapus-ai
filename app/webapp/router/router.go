package router

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	"github.com/vapusdata-ecosystem/vapusai/webapp/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/webapp/routes"
	"github.com/vapusdata-ecosystem/vapusai/webapp/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"limitWords":                limitWords,
		"epochConverter":            EpochConverter,
		"inSlice":                   InSlice,
		"toJSON":                    toJSON,
		"stringCheck":               stringCheck,
		"limitletters":              limitletters,
		"sliceContains":             SliceContains,
		"escapeHTML":                escapeHTML,
		"addRand":                   addRand,
		"randBool":                  randBool,
		"marshalToYaml":             marshalToYaml,
		"epochConverterFull":        EpochConverterFull,
		"epochConverterFullSeconds": EpochConverterFullSeconds,
		"strContains":               strContains,
		"protoToJSON":               protoToJSON,
		"sliceLen":                  sliceLen[any],
		"slugToTitle":               slugToTitle,
		"epochConverterTextDate":    EpochConverterTextDate,
		"escapeJSON":                escapeJSON,
		"joinSlice":                 joinSlice[string],
		"joinElements":              joinElements[any],
		"enumoTitle":                enumoTitle,
		"strTitle":                  strTitle,
		"intCheck":                  intCheck,
		"getSlicelen":               getSlicelen[any],
		"limitlettersWD":            limitlettersWD,
		"floatTruncate":             floatTruncate[any],
		"strUpper":                  strUpper,
		"typeOf":                    typeOf,
		"parseJSON":                 parseJSON,
	}
}

func loadTemplatesRecursively(root string, funcMap template.FuncMap) (*template.Template, error) {
	tmpl := template.New("").Funcs(funcMap)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err := tmpl.ParseFiles(path)
			if err != nil {
				log.Printf("Error parsing template %s: %v", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func GetNewRouter() *echo.Echo {
	var err error
	app := echo.New()
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST},
	}))
	app.Static("/static", "static")
	templates := template.Must(template.New("").Funcs(getFuncMap()).ParseGlob("templates/*.html"))
	filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			template.Must(templates.ParseFiles(path))
		}
		return nil
	})
	renderer := &TemplateRenderer{
		templates: templates,
	}

	app.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	app.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/ui/")
	})
	services.InitWebappSvc()
	app.Renderer = renderer

	gwmux := runtime.NewServeMux(
		runtime.WithOutgoingHeaderMatcher(runtime.DefaultHeaderMatcher),
		runtime.WithIncomingHeaderMatcher(runtime.DefaultHeaderMatcher),
		runtime.WithMarshalerOption("application/protobuf", &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
				UseEnumNumbers:  false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(100*1024*1024),
			grpc.MaxCallSendMsgSize(100*1024*1024),
		),
	}
	err = pb.RegisterVapusdataServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering VapusdataService handler from endpoint")
	}

	err = pb.RegisterDatasourceServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering DataProduct service handler from endpoint")
	}

	err = pb.RegisterPluginServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering PluginService handler from endpoint")
	}
	err = pb.RegisterOrganizationServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering OrganizationService handler from endpoint")
	}
	err = pb.RegisterUserManagementServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering UserManagementService handler from endpoint")
	}

	err = pb.RegisterUtilityServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering UtilityService handler from endpoint")
	}

	err = pb.RegisterAIGuardrailsHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIAgents handler from endpoint")
	}
	err = pb.RegisterAgentServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering Vapus Agents Service from endpoint")
	}
	err = pb.RegisterAgentStudioHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering Vapus Agents Studio from endpoint")
	}
	err = pb.RegisterAIStudioHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIStudio handler from endpoint")
	}
	err = pb.RegisterAIModelsHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIModels handler from endpoint")
	}
	err = pb.RegisterAIPromptsHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIPrompts handler from endpoint")
	}
	err = pb.RegisterAIGuardrailsHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering AIGuardrails handler from endpoint")
	}

	err = pb.RegisterSecretServiceHandlerFromEndpoint(context.Background(), gwmux, pkgs.VapusSvcInternalClientManager.AIStudioDns, opts)
	if err != nil {
		pkgs.DmLogger.Fatal().Err(err).Msg("error while registering SecretService handler from endpoint")
	}

	app.Any("/api/v1alpha1/*", echo.WrapHandler(gwmux))
	uigrp := app.Group(routes.UIRoute)
	homeRouter(uigrp)
	authnRouter(app)
	domainRouter(uigrp)
	CommonRouter(uigrp)
	studioRouters(uigrp)
	manageAIRoutes(uigrp)
	settingsRouter(uigrp)
	developersRouters(uigrp)
	insightsRouters(uigrp)
	// nabhikTasksRouters(uigrp)
	app.HTTPErrorHandler = func(err error, c echo.Context) {
		var (
			code    = http.StatusInternalServerError
			message interface{}
		)

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message
		} else {
			message = http.StatusText(code)
		}

		// If 404, render the custom 404 template
		if code == http.StatusNotFound {
			if err := c.Render(http.StatusNotFound, "404.html", nil); err != nil {
				c.Logger().Error(err)
			}
			return
		}
		if code == http.StatusInternalServerError {
			if err := c.Render(http.StatusInternalServerError, "400.html", nil); err != nil {
				c.Logger().Error(err)
			}
			return
		}
		if code == http.StatusUnauthorized {
			if err := c.Render(http.StatusUnauthorized, "403.html", nil); err != nil {
				c.Logger().Error(err)
			}
			return
		}

		// For other errors, send JSON response or default message
		if !c.Response().Committed {
			c.JSON(code, map[string]interface{}{
				"code":    code,
				"message": message,
			})
		}
	}

	return app
}
