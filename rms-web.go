package main

import (
	"embed"
	"fmt"
	"github.com/RacoonMediaServer/rms-web/internal/cctv"
	"github.com/RacoonMediaServer/rms-web/internal/journal"
	"github.com/RacoonMediaServer/rms-web/internal/multimedia"
	"github.com/RacoonMediaServer/rms-web/internal/services"
	"github.com/RacoonMediaServer/rms-web/internal/settings"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/RacoonMediaServer/rms-web/internal/config"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var Version = "v0.0.0"

const serviceName = "rms-web"

//go:embed web
var webFS embed.FS

//go:embed templates
var templatesFS embed.FS

func main() {
	logger.Infof("%s %s", serviceName, Version)
	defer logger.Info("DONE.")

	useDebug := false

	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version(Version),
		micro.Flags(
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"debug"},
				Usage:       "debug log level",
				Value:       false,
				Destination: &useDebug,
			},
		),
	)

	service.Init(
		micro.Action(func(context *cli.Context) error {
			configFile := fmt.Sprintf("/etc/rms/%s.json", serviceName)
			if context.IsSet("config") {
				configFile = context.String("config")
			}
			return config.Load(configFile)
		}),
	)

	if useDebug {
		_ = logger.Init(logger.WithLevel(logger.DebugLevel))
	}

	f := servicemgr.NewServiceFactory(service)

	root := template.New("root").Funcs(ui.Functions)
	templates := template.Must(root.ParseFS(templatesFS, "templates/*.tmpl"))
	cfg := config.Config()

	web := gin.Default()
	web.SetHTMLTemplate(templates)
	web.StaticFS("/css", http.FS(wrapFS(webFS, "web/css")))
	web.StaticFS("/img", http.FS(wrapFS(webFS, "web/img")))
	web.StaticFS("/js", http.FS(wrapFS(webFS, "web/js")))

	web.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/journal")
	})

	web.NoRoute(func(ctx *gin.Context) {
		ui.DisplayError(ctx, http.StatusNotFound, "Упс! Страницы не существует")
	})

	journalService := journal.New(f)
	journalService.Register(web.Group("/journal"))

	settingsService := settings.New(f)
	settingsService.Register(web.Group("/settings"))

	mediaService := multimedia.New(f)
	mediaService.Register(web.Group("/multimedia"))

	if cfg.Cctv.Enabled {
		cctvService := cctv.New(f)
		cctvService.Register(web.Group("/cctv"))
	}

	services.Register(web.Group("/services"))

	if err := web.Run(fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port)); err != nil {
		logger.Fatalf("Run web server failed: %s", err)
	}
}

func wrapFS(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}

	return sub
}
