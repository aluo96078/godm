package examples

import (
	"context"
	"godm/pkg/odm"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 定義使用者模型，內嵌 ODM，並包含使用者專屬欄位。
// User defines the user model, embedding ODM and including user-specific fields.
type User struct {
	odm.GODM `bson:"-"`
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`

	Posts []Post `bson:"posts,omitempty"`
}

// NewUser 建立一個新的 User 實例，並初始化 ODM。
// NewUser creates a new User instance and initializes ODM.
func NewUser() *User {
	u := &User{}
	u.Use(u)
	return u
}

// SetCollectionName 允許覆蓋默認的集合名稱。
// SetCollectionName allows overriding the default collection name.
func (o *User) SetCollectionName(name string) *User {
	o.CollectionName = name
	if o.Collection != nil {
		o.Collection = odm.MongoClient.Database(o.DbName).Collection(name)
	}
	return o
}

// WithContext 允許覆蓋默認的集合名稱。
// WithContext allows overriding the default collection name.
func (o *User) WithContext(ctx context.Context) *User {
	o.Ctx = ctx
	return o
}
