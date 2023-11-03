package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"os"
	"webhook/src/support"
)

func runServer(args support.Arguments) error {
	if args.LogLevel != "debug" && args.LogLevel != "trace" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(support.LoggingMiddleware())
	r.Use(cors.Default())
	r.Use(static.Serve("/", static.LocalFile(args.StaticContents, false)))
	support.Setup(r, &args)

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
			Name: "static, s", Value: "./static/",
			Usage:       "Static contents path",
			Destination: &args.StaticContents,
		},
		cli.StringFlag{
			Name: "uri", Value: "",
			Usage:       "MongoDB connection uri. eg. mongodb://username:password@localhost:27017/webhook",
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
