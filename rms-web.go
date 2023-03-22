package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/RacoonMediaServer/rms-web/internal/config"
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

	_ = servicemgr.NewServiceFactory(service)

	root := template.New("root")
	templates := template.Must(root.ParseFS(templatesFS, "templates/*.tmpl"))

	web := gin.Default()
	web.SetHTMLTemplate(templates)
	web.StaticFS("/css", http.FS(wrapFS(webFS, "web/css")))

	web.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "main.tmpl", nil)
	})

	web.Run(":8080")
}

func wrapFS(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}

	return sub
}
