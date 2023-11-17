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
	MongoConnectionUri string
	Database           string
	BasicAuthUserName  string
	BasicAuthPassword  string
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

func Setup(router *gin.Engine, args *Arguments) *gin.RouterGroup {
	response := src.NewResponse()
	logger := GetLogger(args)

	mongoClient := connectMongo(args)
	webhookClient := model.NewWebhookClient(mongoClient, args.Database)
	logsClient := model.NewLogClient(mongoClient, args.Database)
	templateClient := model.NewTemplateClient(mongoClient, args.Database)
	emailClient := model.NewEmailClient(mongoClient, args.Database)
	smsClient := model.NewSmsClient(mongoClient, args.Database)

	w := handle.Webhook{MongoClient: mongoClient, Model: webhookClient, Response: response, Logger: logger}

	r := router.Group("/api")
	if args.BasicAuthUserName != "" && args.BasicAuthPassword != "" {
		r = router.Group("/api", gin.BasicAuth(gin.Accounts{
			args.BasicAuthUserName: args.BasicAuthPassword,
		}))
	}
	r.GET("/webhook", w.Query)
	r.POST("/webhook", w.Store)
	r.GET("/webhook/:id", w.Detail)
	r.PUT("/webhook/:id", w.Update)
	r.DELETE("/webhook/:id", w.Delete)
	r.PUT("/webhook/:id/enable", w.Enable)
	r.PUT("/webhook/:id/disable", w.Disable)
	r.POST("/webhook/:id/duplicate", w.Duplicate)
	r.POST("/import", w.Import)

	t := handle.Template{MongoClient: mongoClient, Model: templateClient, Response: response, Logger: logger}
	r.GET("/template", t.Query)
	r.POST("/template", t.Store)
	r.PUT("/template/:id", t.Update)
	r.DELETE("/template/:id", t.Delete)
	r.GET("/template/:id", t.Detail)

	p := handle.SmtpProfile{MongoClient: mongoClient, Model: emailClient, Response: response, Logger: logger}
	r.GET("/smtp/profile", p.Query)
	r.POST("/smtp/profile", p.Store)
	r.PUT("/smtp/profile/:id", p.Update)
	r.DELETE("/smtp/profile/:id", p.Delete)
	r.GET("/smtp/profile/:id", p.Detail)

	sms := handle.SmsProfile{MongoClient: mongoClient, Model: smsClient, Response: response, Logger: logger}
	r.GET("/sms/profile", sms.Query)
	r.POST("/sms/profile", sms.Store)
	r.PUT("/sms/profile/:id", sms.Update)
	r.DELETE("/sms/profile/:id", sms.Delete)
	r.GET("/sms/profile/:id", sms.Detail)

	h := handle.Hook{MongoClient: mongoClient, Model: webhookClient, Response: response, Logger: logger, LogModel: logsClient, EmailModel: emailClient, SmsModel: smsClient}
	router.Any("/hook/:id", h.HandleHook)

	l := handle.Log{MongoClient: mongoClient, Model: logsClient, Response: response}
	r.GET("/logs", l.Query)
	r.DELETE("/logs", l.Clear)

	return r
}
