package odm

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

type GODM struct {
	Collection     *mongo.Collection
	Filter         bson.D
	OrFilter       []bson.M
	Model          interface{}
	LimitCount     int64
	SortFields     bson.D
	SkipCount      int64
	Projection     bson.M
	Ctx            context.Context
	CollectionName string
	DbName         string

	Observers []ModelObserver // 支援多個 observer
}
