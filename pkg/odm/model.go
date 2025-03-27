package odm

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

// GODM：MongoDB for Go 的簡易查詢映射器

// 🧩 簡介
//
// GODM（Go Object-Document Mapper）是一個用於 MongoDB 的輕量級查詢封裝工具，使用 Go 語言實作。它提供類似 ORM 的開發體驗，並針對常見的查詢條件與鏈式操作做了簡化，幫助你快速構建資料模型並執行 CRUD、聚合、事務等操作。
//
// 核心實作位於 [`pkg/odm`](./pkg/odm)，使用範例可見於 [`examples/`](./examples).
//
// ---
//
// ✨ 功能特色
//
// - 🚀 鏈式查詢語法（Where, OrWhere, WhereIn 等）
// - 🔧 自動關聯資料模型與集合（支援自定義集合名與資料庫名）
// - 💾 支援 CRUD 與 BulkCreate
// - 🧠 支援複雜查詢條件組合（AND / OR）
// - 🔁 支援 MongoDB 聚合管道
// - 💼 內建事務封裝 `WithTransaction`
// - 🧪 簡潔易測試，模組化設計便於擴展
//
// ---
//
// 🛠 使用方式（以 User 模型為例）
//
// `examples/user.go` 定義了一個使用 GODM 的 User 模型，以下是簡要示例：
//
// ```go
// type User struct {
//     GODM   `bson:"-"`
//     ID     primitive.ObjectID `bson:"_id,omitempty"`
//     Name   string             `bson:"name"`
//     Email  string             `bson:"email"`
// }
//
// func NewUser() *User {
//     u := &User{}
//     u.Use(u)
//     return u
// }
// ```
//
// ### 建立與查詢
//
// ```go
// user := NewUser()
// user.Name = "Test"
// user.Email = "test@example.com"
// _ = user.Create()
//
// // 查詢第一筆資料
// err := user.Where("email", "=", "test@example.com").First()
// ```
//
// ### 聚合與事務操作
//
// ```go
// pipeline := mongo.Pipeline{
//     {{"$match", bson.M{"email": bson.M{"$ne": ""}}}},
//     {{"$group", bson.M{"_id": "$name", "count": bson.M{"$sum": 1}}}},
// }
// var result []bson.M
// _ = user.Aggregate(pipeline, &result)
//
// _ = user.WithTransaction(func(sess mongo.SessionContext) error {
//     return user.Update(bson.M{"name": "Updated"})
// })
// ```
//
// ---
//
// 📂 專案結構
//
// ```
// gotest/
// ├── examples/        # 使用範例：main.go, user.go
// ├── pkg/
// │   └── odm/         # GODM 核心實作（已模組化）
// │       ├── crud.go
// │       ├── query.go
// │       ├── config.go
// │       ├── context.go
// │       ├── transaction.go
// │       ├── aggregate.go
// │       └── util.go
// ├── go.mod
// └── README.md
// ```
//
// ---
//
// 📄 授權
//
// 本專案採用 [MIT License](./LICENSE) 授權。

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
