package handle

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"webhook/src/model"
)

type Log struct {
	MongoClient *mongo.Client
	Model       *model.LogClient
}

func (l *Log) Query(c *gin.Context) {}
