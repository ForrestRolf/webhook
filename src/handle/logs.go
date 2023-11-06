package handle

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"webhook/src"
	"webhook/src/model"
)

type Log struct {
	MongoClient *mongo.Client
	Model       *model.LogClient
	Response    *src.Response
}

func (l *Log) Query(c *gin.Context) {
	id := c.Query("id")
	logs, err := l.Model.QueryLogs(id)
	if err != nil {
		l.Response.Fail(c, err.Error(), nil)
		return
	}
	l.Response.Success(c, logs, "OK")
}
