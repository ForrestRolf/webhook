package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"webhook/src"
	"webhook/src/model"
)

type Template struct {
	MongoClient *mongo.Client
	Model       *model.TemplateClient
	Response    *src.Response
	Logger      *logrus.Logger
}

func (t *Template) Query(c *gin.Context) {
	templates, err := t.Model.ListTemplates()
	if err != nil {
		t.Response.Fail(c, err.Error(), nil)
		return
	}
	t.Response.Success(c, templates, "")
}
func (t *Template) Store(c *gin.Context) {
	var template *model.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		t.Response.BadRequest(c, err.Error(), nil)
		return
	}
	if err := t.Model.AddTemplate(template); err != nil {
		t.Response.Fail(c, err.Error(), nil)
		return
	}
	t.Response.Success(c, template, "")
}
func (t *Template) Detail(c *gin.Context) {
	id := c.Param("id")
	template, err := t.Model.GetTemplate(id)
	if err != nil {
		t.Response.Fail(c, err.Error(), nil)
		return
	}
	t.Response.Success(c, template, "")
}
func (t *Template) Update(c *gin.Context) {
	id := c.Param("id")
	var template model.Template

	if err := c.ShouldBindJSON(&template); err != nil {
		t.Response.BadRequest(c, err.Error(), nil)
		return
	}
	_, err := t.Model.UpdateTemplate(id, template)
	if err != nil {
		t.Response.Fail(c, err.Error(), nil)
		return
	}
	t.Response.Success(c, template, "")
}
func (t *Template) Delete(c *gin.Context) {
	id := c.Param("id")
	deleteCount, err := t.Model.DeleteTemplate(id)
	if err != nil {
		t.Response.Fail(c, err.Error(), nil)
		return
	}
	t.Response.Success(c, deleteCount, "")
}
