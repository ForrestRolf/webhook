package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	webhooks, err := w.Model.ListWebhooks()
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
	modifiedCount, err := w.Model.UpdateWebhook(id, webhook)
	if err != nil {
		w.Response.Fail(c, err.Error(), nil)
		return
	}
	if modifiedCount == 0 {
		w.Response.NotFound(c, "")
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
