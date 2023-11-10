package model

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"time"
	"webhook/src/hook"
)

type LogClient struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

type WebhookLogClient struct {
	client *LogClient
	hook   *hook.Hook
}

type ActionLogClient struct {
	id     string
	driver string
	client *LogClient
	hook   *hook.Hook
}

type Log struct {
	Type         string             `json:"type,omitempty" bson:"type"`
	LogId        string             `json:"logId,omitempty" bson:"logId"`
	WebhookId    string             `json:"webhookId" bson:"webhookId"`
	WebhookName  string             `json:"webhookName" bson:"webhookName"`
	Message      string             `json:"message" bson:"message"`
	Level        string             `json:"level" bson:"level"`
	Created      primitive.DateTime `json:"created" bson:"created"`
	ActionDriver string             `json:"actionDriver,omitempty" bson:"actionDriver"`
}

const (
	LogLevelWarn  string = "warn"
	LogLevelError string = "error"
	LogLevelInfo  string = "info"
	LogLevelDebug string = "debug"
	LogLevelOk    string = "OK"
)
const (
	LogTypeWebhook string = "webhook"
	LogTypeAction  string = "action"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func NewLogClient(client *mongo.Client, db string) *LogClient {
	return &LogClient{
		client:        client,
		logCollection: client.Database(db).Collection("logs"),
	}
}

func NewWebhookLogClient(client *LogClient, hook *hook.Hook) *WebhookLogClient {
	return &WebhookLogClient{
		client: client,
		hook:   hook,
	}
}

func NewActionLogClient(driver string, logId string, client *LogClient, hook *hook.Hook) *ActionLogClient {
	return &ActionLogClient{
		id:     logId,
		driver: driver,
		client: client,
		hook:   hook,
	}
}

func (l *LogClient) NewLog(log Log) {
	_, err := l.logCollection.InsertOne(context.TODO(), log)
	if err != nil {
		logrus.Errorf("Could not save log: %s", err.Error())
	}
}
func (l *LogClient) NewWebhookLog(webhook *hook.Hook, msg string, level string) {
	var log = Log{
		Type:        LogTypeWebhook,
		LogId:       webhook.Name,
		WebhookId:   webhook.ID,
		WebhookName: webhook.Name,
		Message:     msg,
		Level:       level,
		Created:     primitive.NewDateTimeFromTime(time.Now()),
	}
	l.NewLog(log)
}
func (l *LogClient) NewActionLog(driver string, logId string, webhook *hook.Hook, msg string, level string) {
	var log = Log{
		Type:         LogTypeAction,
		LogId:        "a-" + logId,
		WebhookId:    webhook.ID,
		WebhookName:  webhook.Name,
		Message:      msg,
		Level:        level,
		Created:      primitive.NewDateTimeFromTime(time.Now()),
		ActionDriver: driver,
	}
	l.NewLog(log)
}

func (l *LogClient) ParseErrors(errors []error) string {
	e := ""
	for _, err := range errors {
		e += "," + err.Error()
	}
	return e
}

func (l *LogClient) GenerateLogId(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (l *LogClient) QueryLogs(id string, limit int64) ([]Log, error) {
	logs := make([]Log, 0)

	opts := options.Find().SetSort(bson.D{{"_id", -1}}).SetLimit(limit)
	filter := bson.D{}
	if id != "" {
		filter = bson.D{{"webhookId", id}}
	}
	cur, err := l.logCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return logs, err
	}
	if err = cur.All(context.TODO(), &logs); err != nil {
		return logs, err
	}
	return logs, nil
}

func (l *LogClient) ClearLogs(days int) (int64, error) {
	t := time.Now().Add(time.Duration(-days*24) * time.Hour)
	filter := bson.D{{"created", bson.D{{"$lt", primitive.NewDateTimeFromTime(t)}}}}
	results, err := l.logCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return results.DeletedCount, nil
}

func (w *WebhookLogClient) AddErrorLog(msg string) {
	w.client.NewWebhookLog(w.hook, msg, LogLevelError)
}
func (w *WebhookLogClient) AddWarnLog(msg string) {
	w.client.NewWebhookLog(w.hook, msg, LogLevelWarn)
}
func (w *WebhookLogClient) AddInfoLog(msg string) {
	w.client.NewWebhookLog(w.hook, msg, LogLevelInfo)
}
func (w *WebhookLogClient) AddDebugLog(msg string) {
	if w.hook.Debug {
		w.client.NewWebhookLog(w.hook, msg, LogLevelDebug)
	}
}
func (w *WebhookLogClient) AddLog(msg string) {
	w.client.NewWebhookLog(w.hook, msg, LogLevelOk)
}

func (l *ActionLogClient) AddLog(msg string) {
	l.client.NewActionLog(l.driver, l.id, l.hook, msg, LogLevelOk)
}
func (l *ActionLogClient) AddErrorLog(msg string) {
	l.client.NewActionLog(l.driver, l.id, l.hook, msg, LogLevelError)
}
func (l *ActionLogClient) AddWarnLog(msg string) {
	l.client.NewActionLog(l.driver, l.id, l.hook, msg, LogLevelWarn)
}
func (l *ActionLogClient) AddInfoLog(msg string) {
	l.client.NewActionLog(l.driver, l.id, l.hook, msg, LogLevelInfo)
}
func (l *ActionLogClient) AddDebugLog(msg string) {
	if l.hook.Debug {
		l.client.NewActionLog(l.driver, l.id, l.hook, msg, LogLevelDebug)
	}
}
