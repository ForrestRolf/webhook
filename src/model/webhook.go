package model

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"webhook/src/hook"
)

type Rules struct {
	And   *hook.AndRule   `json:"and,omitempty"`
	Or    *hook.OrRule    `json:"or,omitempty"`
	Match *hook.MatchRule `json:"match,omitempty"`
}

type Webhook struct {
	Id                    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name                  string             `json:"name,omitempty"`
	Description           string             `json:"description,omitempty"`
	Triggers              *hook.Rules        `json:"triggers,omitempty"`
	Actions               *[]hook.Action     `json:"actions,omitempty"`
	Enabled               bool               `json:"enabled,omitempty"`
	PassArgumentsToAction []hook.Argument    `json:"passArgumentsToAction,omitempty" bson:"passArgumentsToAction"`
	RunCount              int                `json:"runCount" bson:"runCount"`
	LastRunAt             primitive.DateTime `json:"lastRunAt,omitempty" bson:"lastRunAt,omitempty"`
	CallCount             int                `json:"callCount" bson:"callCount"`
	AuthToken             string             `json:"authToken,omitempty" bson:"authToken,omitempty"`
	Debug                 bool               `json:"debug,omitempty" bson:"debug"`
	SaveRequest           []string           `json:"saveRequest,omitempty" bson:"saveRequest"`
}

type WebhookClient struct {
	client            *mongo.Client
	webhookCollection *mongo.Collection
}

func NewWebhookClient(client *mongo.Client, db string) *WebhookClient {
	return &WebhookClient{
		client:            client,
		webhookCollection: client.Database(db).Collection("webhooks"),
	}
}

func (c *WebhookClient) AddWebhook(webhook *Webhook) error {
	webhook.Enabled = true
	webhook.RunCount = 0
	webhook.CallCount = 0
	_, err := c.webhookCollection.InsertOne(context.TODO(), webhook)
	if err != nil {
		return err
	}
	return nil
}

func (c *WebhookClient) ListWebhooks() ([]Webhook, error) {
	webhooks := make([]Webhook, 0)
	opts := options.Find().SetSort(bson.D{{"_id", -1}})
	cur, err := c.webhookCollection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	if err = cur.All(context.TODO(), &webhooks); err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (c *WebhookClient) GetWebhook(id string) (Webhook, error) {
	var webhook Webhook
	objectId, _ := primitive.ObjectIDFromHex(id)
	res := c.webhookCollection.FindOne(context.TODO(), bson.M{"_id": objectId})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return webhook, nil
		}
		return webhook, res.Err()
	}
	if err := res.Decode(&webhook); err != nil {
		return webhook, err
	}
	return webhook, nil
}

func (c *WebhookClient) UpdateWebhook(id string, webhook Webhook) (int, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	res, err := c.webhookCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.D{{
		"$set", bson.D{
			{"name", webhook.Name},
			{"description", webhook.Description},
			{"triggers", webhook.Triggers},
			{"actions", webhook.Actions},
			{"passArgumentsToAction", webhook.PassArgumentsToAction},
			{"authToken", webhook.AuthToken},
			{"saveRequest", webhook.SaveRequest},
			{"debug", webhook.Debug},
		},
	}})
	if err != nil {
		return 0, err
	}
	return int(res.ModifiedCount), nil
}

func (c *WebhookClient) DeleteWebhook(id string) (int64, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	res, err := c.webhookCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (c *WebhookClient) IncreaseCount(id string, field string) (int, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	res, err := c.webhookCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.D{
		{
			"$inc", bson.D{
				{field, 1},
			},
		},
		{
			"$set", bson.D{
				{"lastRunAt", primitive.NewDateTimeFromTime(time.Now())},
			},
		},
	})
	if err != nil {
		return 0, err
	}
	return int(res.ModifiedCount), nil
}

func (c *WebhookClient) SetEnabled(id string, enable bool) (int, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	res, err := c.webhookCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.D{{
		"$set", bson.D{
			{"enabled", enable},
		},
	}})
	if err != nil {
		return 0, err
	}
	return int(res.ModifiedCount), nil
}
