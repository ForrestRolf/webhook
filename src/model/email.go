package model

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EmailClient struct {
	client            *mongo.Client
	profileCollection *mongo.Collection
}

type SmtpProfile struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name"`
	Host     string             `json:"host,omitempty" bson:"host"`
	Port     int                `json:"port,omitempty" bson:"port"`
	Tls      bool               `json:"tls,omitempty" bson:"tls"`
	Username string             `json:"username,omitempty" bson:"username"`
	Password string             `json:"password,omitempty" bson:"password"`
	Sender   string             `json:"sender,omitempty" bson:"sender"`
}

func NewEmailClient(client *mongo.Client, db string) *EmailClient {
	return &EmailClient{
		client:            client,
		profileCollection: client.Database(db).Collection("smtp_profiles"),
	}
}

func (e *EmailClient) AddProfile(profile *SmtpProfile) error {
	_, err := e.profileCollection.InsertOne(context.TODO(), profile)
	if err != nil {
		return err
	}
	return nil
}

func (e *EmailClient) ListProfiles() ([]SmtpProfile, error) {
	profiles := make([]SmtpProfile, 0)
	opts := options.Find().SetSort(bson.D{{"_id", -1}})
	cur, err := e.profileCollection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	if err = cur.All(context.TODO(), &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

func (e *EmailClient) GetProfile(id string) (SmtpProfile, error) {
	var profile SmtpProfile
	objectId, _ := primitive.ObjectIDFromHex(id)
	res := e.profileCollection.FindOne(context.TODO(), bson.M{"_id": objectId})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return profile, nil
		}
		return profile, res.Err()
	}
	if err := res.Decode(&profile); err != nil {
		return profile, err
	}
	return profile, nil
}

func (e *EmailClient) DeleteProfile(id string) (int64, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	res, err := e.profileCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (e *EmailClient) UpdateProfile(id string, profile SmtpProfile) (int, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	res, err := e.profileCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.D{{
		"$set", bson.D{
			{"name", profile.Name},
			{"host", profile.Host},
			{"port", profile.Port},
			{"username", profile.Username},
			{"password", profile.Password},
			{"sender", profile.Sender},
			{"tls", profile.Tls},
		},
	}})
	if err != nil {
		return 0, err
	}
	return int(res.ModifiedCount), nil
}
