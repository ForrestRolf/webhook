package model

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"webhook/src/hook"
)

type LogClient struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

type Log struct {
	WebhookId   string             `json:"webhookId" bson:"webhookId"`
	WebhookName string             `json:"webhookName" bson:"webhookName"`
	Message     string             `json:"message" bson:"message"`
	Level       string             `json:"level" bson:"level"`
	Created     primitive.DateTime `json:"created" bson:"created"`
}

const (
	LogLevelWarn  string = "warn"
	LogLevelError string = "error"
	LogLevelInfo  string = "info"
	LogLevelDebug string = "debug"
	LogLevelOk    string = "OK"
)

func NewLogClient(client *mongo.Client, db string) *LogClient {
	return &LogClient{
		client:        client,
		logCollection: client.Database(db).Collection("logs"),
	}
}

func (l *LogClient) newLog(webhook *hook.Hook, msg string, level string) {
	var log = Log{
		WebhookId:   webhook.ID,
		WebhookName: webhook.Name,
		Message:     msg,
		Level:       level,
		Created:     primitive.NewDateTimeFromTime(time.Now()),
	}
	_, err := l.logCollection.InsertOne(context.TODO(), log)
	if err != nil {
		logrus.Errorf("Could not save log: %w", err)
	}
}

func (l *LogClient) AddErrorLog(webhook *hook.Hook, msg string) {
	l.newLog(webhook, msg, LogLevelError)
}

func (l *LogClient) AddWarnLog(webhook *hook.Hook, msg string) {
	l.newLog(webhook, msg, LogLevelWarn)
}

func (l *LogClient) AddInfoLog(webhook *hook.Hook, msg string) {
	l.newLog(webhook, msg, LogLevelInfo)
}

func (l *LogClient) AddDebugLog(webhook *hook.Hook, msg string) {
	l.newLog(webhook, msg, LogLevelDebug)
}

func (l *LogClient) AddLog(webhook *hook.Hook, msg string) {
	l.newLog(webhook, msg, LogLevelOk)
}

func (l *LogClient) ParseErrors(errors []error) string {
	e := ""
	for _, err := range errors {
		e += "," + err.Error()
	}
	return e
}
