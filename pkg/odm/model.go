package odm

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

type GODM struct {
	Collection     *mongo.Collection
	Filter         bson.D   // AND 條件 // AND conditions
	OrFilter       []bson.M // OR 條件 // OR conditions
	Model          interface{}
	LimitCount     int64
	SortFields     bson.D
	SkipCount      int64
	Projection     bson.M
	Ctx            context.Context
	CollectionName string
	DbName         string // 資料庫名稱 // Database name
}
