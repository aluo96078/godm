package examples

import (
	"context"
	"godm/pkg/odm"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	odm.GODM `bson:"-"`
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserID   primitive.ObjectID `bson:"user_id"`
	Title    string             `bson:"title"`
	Body     string             `bson:"body"`

	User *User `bson:"user,omitempty"` // 新增 User 欄位

}

// NewUser 建立一個新的 User 實例，並初始化 ODM。
// NewUser creates a new User instance and initializes ODM.
func NewPost() *Post {
	p := &Post{}
	p.Use(p)
	return p
}

// SetCollectionName 允許覆蓋默認的集合名稱。
// SetCollectionName allows overriding the default collection name.
func (o *Post) SetCollectionName(name string) *Post {
	o.CollectionName = name
	if o.Collection != nil {
		o.Collection = odm.MongoClient.Database(o.DbName).Collection(name)
	}
	return o
}

// WithContext 允許覆蓋默認的集合名稱。
// WithContext allows overriding the default collection name.
func (o *Post) WithContext(ctx context.Context) *Post {
	o.Ctx = ctx
	return o
}
