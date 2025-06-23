package odm

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client
var DBName string

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
	DBName         string

	Observers []ModelObserver // 支援多個 observer

	// 預先載入關聯的欄位名稱（例如 "posts", "comments"）
	WithRelations []string

	// 關聯欄位對應的設定，例如 localField, foreignField 等（未來可用來自定義 $lookup 行為）
	RelationConfigs map[string]RelationConfig
}

// RelationConfig 用來定義一個 $lookup 的設定
type RelationConfig struct {
	From         string // 關聯的 collection 名稱
	LocalField   string // 本地欄位
	ForeignField string // 關聯 collection 的欄位
	As           string // 最終回傳的欄位名稱
	IsArray      bool   // 是否為一對多（true）或一對一（false）
}
