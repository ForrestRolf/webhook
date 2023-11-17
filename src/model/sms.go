package model

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SmsClient struct {
	client            *mongo.Client
	profileCollection *mongo.Collection
}

type SmsProfile struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Provider string             `json:"provider,omitempty" bson:"provider"`
	Name     string             `json:"name,omitempty" bson:"name"`
	AK       string             `json:"ak,omitempty" bson:"ak"`
	SK       string             `json:"sk,omitempty" bson:"sk"`
	From     string             `json:"from,omitempty" bson:"from"`
}

func NewSmsClient(client *mongo.Client, db string) *SmsClient {
	return &SmsClient{
		client:            client,
		profileCollection: client.Database(db).Collection("sms_profiles"),
	}
}

func (sms *SmsClient) AddProfile(profile *SmsProfile) error {
	_, err := sms.profileCollection.InsertOne(context.TODO(), profile)
	if err != nil {
		return err
	}
	return nil
}

func (sms *SmsClient) ListProfiles() ([]SmsProfile, error) {
	profiles := make([]SmsProfile, 0)
	opts := options.Find().SetSort(bson.D{{"_id", -1}})
	cur, err := sms.profileCollection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	if err = cur.All(context.TODO(), &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

func (sms *SmsClient) GetProfile(id string) (SmsProfile, error) {
	var profile SmsProfile
	objectId, _ := primitive.ObjectIDFromHex(id)
	res := sms.profileCollection.FindOne(context.TODO(), bson.M{"_id": objectId})
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

func (sms *SmsClient) DeleteProfile(id string) (int64, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	res, err := sms.profileCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (sms *SmsClient) UpdateProfile(id string, profile SmsProfile) (int, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	res, err := sms.profileCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.D{{
		"$set", bson.D{
			{"provider", profile.Provider},
			{"ak", profile.AK},
			{"sk", profile.SK},
			{"from", profile.From},
			{"name", profile.Name},
		},
	}})
	if err != nil {
		return 0, err
	}
	return int(res.ModifiedCount), nil
}
