package model

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TemplateClient struct {
	client             *mongo.Client
	templateCollection *mongo.Collection
}

type Template struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title"`
	Description string             `json:"description,omitempty" bson:"description"`
	Language    string             `json:"language,omitempty" bson:"language"`
	Content     string             `json:"content,omitempty" bson:"content"`
	RefCount    int                `json:"refCount,omitempty" bson:"refCount"`
	Created     primitive.DateTime `json:"created,omitempty" bson:"created"`
	Updated     primitive.DateTime `json:"updated,omitempty" bson:"updated"`
}

func NewTemplateClient(client *mongo.Client, db string) *TemplateClient {
	return &TemplateClient{
		client:             client,
		templateCollection: client.Database(db).Collection("templates"),
	}
}

func (t *TemplateClient) AddTemplate(template *Template) error {
	template.Created = primitive.NewDateTimeFromTime(time.Now())
	template.Updated = primitive.NewDateTimeFromTime(time.Now())
	template.RefCount = 0
	_, err := t.templateCollection.InsertOne(context.TODO(), template)
	if err != nil {
		return err
	}
	return nil
}

func (t *TemplateClient) ListTemplates() ([]Template, error) {
	templates := make([]Template, 0)
	opts := options.Find().SetSort(bson.D{{"_id", -1}})
	cur, err := t.templateCollection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	if err = cur.All(context.TODO(), &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (t *TemplateClient) GetTemplate(id string) (Template, error) {
	var template Template
	objectId, _ := primitive.ObjectIDFromHex(id)
	res := t.templateCollection.FindOne(context.TODO(), bson.M{"_id": objectId})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return template, nil
		}
		return template, res.Err()
	}
	if err := res.Decode(&template); err != nil {
		return template, err
	}
	return template, nil
}

func (t *TemplateClient) UpdateTemplate(id string, template Template) (int, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	res, err := t.templateCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.D{{
		"$set", bson.D{
			{"title", template.Title},
			{"description", template.Description},
			{"content", template.Content},
			{"updated", primitive.NewDateTimeFromTime(time.Now())},
		},
	}})
	if err != nil {
		return 0, err
	}
	return int(res.ModifiedCount), nil
}

func (t *TemplateClient) DeleteTemplate(id string) (int64, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	res, err := t.templateCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
