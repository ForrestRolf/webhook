package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"webhook/src"
	"webhook/src/model"
)

type SmtpProfile struct {
	MongoClient *mongo.Client
	Model       *model.EmailClient
	Response    *src.Response
	Logger      *logrus.Logger
}

func (p *SmtpProfile) Query(c *gin.Context) {
	profiles, err := p.Model.ListProfiles()
	if err != nil {
		p.Response.Fail(c, err.Error(), nil)
		return
	}
	p.Response.Success(c, profiles, "")
}

func (p *SmtpProfile) Store(c *gin.Context) {
	var profile *model.SmtpProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		p.Response.BadRequest(c, err.Error(), nil)
		return
	}
	if err := p.Model.AddProfile(profile); err != nil {
		p.Response.Fail(c, err.Error(), nil)
		return
	}
	p.Response.Success(c, profile, "")
}

func (p *SmtpProfile) Update(c *gin.Context) {
	id := c.Param("id")
	var profile model.SmtpProfile

	if err := c.ShouldBindJSON(&profile); err != nil {
		p.Response.BadRequest(c, err.Error(), nil)
		return
	}
	_, err := p.Model.UpdateProfile(id, profile)
	if err != nil {
		p.Response.Fail(c, err.Error(), nil)
		return
	}
	p.Response.Success(c, profile, "")
}

func (p *SmtpProfile) Detail(c *gin.Context) {
	id := c.Param("id")
	template, err := p.Model.GetProfile(id)
	if err != nil {
		p.Response.Fail(c, err.Error(), nil)
		return
	}
	p.Response.Success(c, template, "")
}

func (p *SmtpProfile) Delete(c *gin.Context) {
	id := c.Param("id")
	deleteCount, err := p.Model.DeleteProfile(id)
	if err != nil {
		p.Response.Fail(c, err.Error(), nil)
		return
	}
	p.Response.Success(c, deleteCount, "")
}
