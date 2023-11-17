package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"webhook/src"
	"webhook/src/model"
)

type SmsProfile struct {
	MongoClient *mongo.Client
	Model       *model.SmsClient
	Response    *src.Response
	Logger      *logrus.Logger
}

func (sms *SmsProfile) Query(c *gin.Context) {
	profiles, err := sms.Model.ListProfiles()
	if err != nil {
		sms.Response.Fail(c, err.Error(), nil)
		return
	}
	sms.Response.Success(c, profiles, "")
}

func (sms *SmsProfile) Store(c *gin.Context) {
	var profile *model.SmsProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		sms.Response.BadRequest(c, err.Error(), nil)
		return
	}
	if err := sms.Model.AddProfile(profile); err != nil {
		sms.Response.Fail(c, err.Error(), nil)
		return
	}
	sms.Response.Success(c, profile, "")
}

func (sms *SmsProfile) Update(c *gin.Context) {
	id := c.Param("id")
	var profile model.SmsProfile

	if err := c.ShouldBindJSON(&profile); err != nil {
		sms.Response.BadRequest(c, err.Error(), nil)
		return
	}
	_, err := sms.Model.UpdateProfile(id, profile)
	if err != nil {
		sms.Response.Fail(c, err.Error(), nil)
		return
	}
	sms.Response.Success(c, profile, "")
}

func (sms *SmsProfile) Detail(c *gin.Context) {
	id := c.Param("id")
	template, err := sms.Model.GetProfile(id)
	if err != nil {
		sms.Response.Fail(c, err.Error(), nil)
		return
	}
	sms.Response.Success(c, template, "")
}

func (sms *SmsProfile) Delete(c *gin.Context) {
	id := c.Param("id")
	deleteCount, err := sms.Model.DeleteProfile(id)
	if err != nil {
		sms.Response.Fail(c, err.Error(), nil)
		return
	}
	sms.Response.Success(c, deleteCount, "")
}
