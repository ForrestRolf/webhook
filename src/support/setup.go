package support

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"webhook/src"
	"webhook/src/handle"
	"webhook/src/model"
)

type Arguments struct {
	LogLevel           string
	BindAddress        string
	BindPort           int
	StaticContents     string
	MongoConnectionUri string
	Database           string
}

func connectMongo(args *Arguments) *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(args.MongoConnectionUri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	return client
}

func Setup(router *gin.Engine, args *Arguments) {
	response := src.NewResponse()
	logger := GetLogger(args)

	mongoClient := connectMongo(args)
	webhookClient := model.NewWebhookClient(mongoClient, args.Database)
	logsClient := model.NewLogClient(mongoClient, args.Database)

	w := handle.Webhook{MongoClient: mongoClient, Model: webhookClient, Response: response, Logger: logger}
	router.GET("/webhook", w.Query)
	router.POST("/webhook", w.Store)
	router.GET("/webhook/:id", w.Detail)
	router.PUT("/webhook/:id", w.Update)
	router.DELETE("/webhook/:id", w.Delete)
	router.PUT("/webhook/:id/enable", w.Enable)
	router.PUT("/webhook/:id/disable", w.Disable)
	router.POST("/webhook/:id/duplicate", w.Duplicate)
	router.POST("/import", w.Import)

	h := handle.Hook{MongoClient: mongoClient, Model: webhookClient, Response: response, Logger: logger, LogModel: logsClient}
	router.Any("/hook/:id", h.HandleHook)

	l := handle.Log{MongoClient: mongoClient, Model: logsClient, Response: response}
	router.GET("/logs", l.Query)
}
