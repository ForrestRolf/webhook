package main

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"io/fs"
	"net/http"
	"os"
	"webhook/src/support"
)

//go:embed web/dist
var web embed.FS

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}
func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

func runServer(args support.Arguments) error {
	if args.LogLevel != "debug" && args.LogLevel != "trace" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(support.LoggingMiddleware())
	r.Use(cors.Default())

	www := EmbedFolder(web, "web/dist")
	r.Use(static.Serve("/", www))
	support.Setup(r, &args)

	r.NoRoute(func(c *gin.Context) {
		c.FileFromFS("index.html", www)
	})
	if err := r.Run(fmt.Sprintf("%s:%d", args.BindAddress, args.BindPort)); err != nil {
		return err
	}

	return nil
}

func main() {
	var args support.Arguments

	app := cli.NewApp()
	app.Name = "webhook"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "log, l", Value: "info",
			Usage:       "Log level [trace,debug,info,warn,error]",
			Destination: &args.LogLevel,
		},
		cli.StringFlag{
			Name: "addr, a", Value: "127.0.0.1",
			Usage:       "Bind address",
			Destination: &args.BindAddress,
		},
		cli.IntFlag{
			Name: "port, p", Value: 9000,
			Usage:       "Bind port",
			Destination: &args.BindPort,
		},
		cli.StringFlag{
			Name: "uri", Value: "",
			Usage:       "MongoDB connection uri. eg. mongodb://username:password@localhost:27017",
			Destination: &args.MongoConnectionUri,
			Required:    true,
		},
		cli.StringFlag{
			Name: "db", Value: "webhook",
			Destination: &args.Database,
		},
	}

	app.Action = func(c *cli.Context) error {
		if err := runServer(args); err != nil {
			return err
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		panic("Fatal Error")
	}
}
