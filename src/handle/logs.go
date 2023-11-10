package handle

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
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
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 1000
	}
	logs, err := l.Model.QueryLogs(id, int64(limit))
	if err != nil {
		l.Response.Fail(c, err.Error(), nil)
		return
	}
	l.Response.Success(c, logs, "OK")
}

func (l *Log) Clear(c *gin.Context) {
	count, err := l.Model.ClearLogs(30)

	if err != nil {
		l.Response.Fail(c, err.Error(), nil)
		return
	}
	l.Response.Success(c, count, "")
}
