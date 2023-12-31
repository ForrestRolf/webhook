package handle

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"webhook/src"
	"webhook/src/model"
)

type Webhook struct {
	MongoClient *mongo.Client
	Model       *model.WebhookClient
	Response    *src.Response
	Logger      *logrus.Logger
}

func (w *Webhook) Query(c *gin.Context) {
	orderBy := c.Query("orderBy")
	webhooks, err := w.Model.ListWebhooks(orderBy)
	if err != nil {
		w.Response.Fail(c, err.Error(), nil)
		return
	}
	w.Response.Success(c, webhooks, "")
}

func (w *Webhook) Store(c *gin.Context) {
	var webhook *model.Webhook

	if err := c.ShouldBindJSON(&webhook); err != nil {
		w.Response.BadRequest(c, err.Error(), nil)
		return
	}
	if err := w.Model.AddWebhook(webhook); err != nil {
		w.Response.Fail(c, err.Error(), nil)
		return
	}
	w.Response.Success(c, webhook, "")
}

func (w *Webhook) Update(c *gin.Context) {
	id := c.Param("id")
	var webhook model.Webhook

	if err := c.ShouldBindJSON(&webhook); err != nil {
		w.Response.BadRequest(c, err.Error(), nil)
		return
	}
	_, err := w.Model.UpdateWebhook(id, webhook)
	if err != nil {
		w.Response.Fail(c, err.Error(), nil)
		return
	}
	w.Response.Success(c, webhook, "")
}

func (w *Webhook) Detail(c *gin.Context) {
	id := c.Param("id")
	webhook, err := w.Model.GetWebhook(id)
	if err != nil {
		w.Response.Fail(c, err.Error(), nil)
		return
	}
	w.Response.Success(c, webhook, "")
}

func (w *Webhook) Delete(c *gin.Context) {
	id := c.Param("id")
	deletedCount, err := w.Model.DeleteWebhook(id)
	if err != nil {
		w.Response.Fail(c, err.Error(), nil)
		return
	}
	w.Response.Success(c, deletedCount, "")
}

func (w *Webhook) Enable(c *gin.Context) {
	id := c.Param("id")
	modifiedCount, err := w.Model.SetEnabled(id, true)
	if err != nil {
		w.Response.Fail(c, err.Error(), nil)
		return
	}
	if modifiedCount == 0 {
		w.Response.NotFound(c, "")
		return
	}
	w.Response.Success(c, nil, "")
}

func (w *Webhook) Disable(c *gin.Context) {
	id := c.Param("id")
	modifiedCount, err := w.Model.SetEnabled(id, false)
	if err != nil {
		w.Response.Fail(c, err.Error(), nil)
		return
	}
	if modifiedCount == 0 {
		w.Response.NotFound(c, "")
		return
	}
	w.Response.Success(c, nil, "")
}

func (w *Webhook) Duplicate(c *gin.Context) {
	id := c.Param("id")
	webhook, err := w.Model.GetWebhook(id)
	if err != nil {
		w.Response.Fail(c, err.Error(), nil)
		return
	}
	webhook.Name = webhook.Name + " Copy"
	webhook.Id = primitive.NewObjectID()

	if err := w.Model.AddWebhook(&webhook); err != nil {
		w.Response.Fail(c, err.Error(), nil)
		return
	}
	w.Response.Success(c, webhook, "")
}

func (w *Webhook) Import(c *gin.Context) {
	file, _ := c.FormFile("file")

	var webhook model.Webhook
	f, _ := file.Open()
	content, _ := ioutil.ReadAll(f)
	err := json.Unmarshal(content, &webhook)
	if err != nil {
		w.Response.Fail(c, err.Error(), "")
		return
	}

	hook := model.Webhook{
		Name:                  webhook.Name,
		Description:           webhook.Description,
		Triggers:              webhook.Triggers,
		Actions:               webhook.Actions,
		PassArgumentsToAction: webhook.PassArgumentsToAction,
		RunCount:              0,
		CallCount:             0,
		AuthToken:             webhook.AuthToken,
		Debug:                 webhook.Debug,
		SaveRequest:           webhook.SaveRequest,
	}
	if err := w.Model.AddWebhook(&hook); err != nil {
		w.Response.Fail(c, err.Error(), nil)
		return
	}
	w.Response.Success(c, webhook, "")
}
