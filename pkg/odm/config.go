package odm

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// Use 設置模型和集合。如果通過 SetCollectionName 提供了自定義集合名稱，則使用該名稱；否則，默認使用模型名的小寫並附加 "s"。
// Use sets the model and collection. If a custom collection name is provided via SetCollectionName, it will use that; otherwise, it defaults to the lowercase model name with an appended "s".
func (o *GODM) Use(model interface{}) *GODM {
	// Determine the database name
	if DBName == "" && o.DBName == "" {
		panic("DBName is not set")
	} else if o.DBName == "" {
		o.DBName = DBName
	}
	modelType := reflect.TypeOf(model).Elem()
	defaultCollectionName := strings.ToLower(modelType.Name()) + "s"
	if o.CollectionName != "" {
		defaultCollectionName = o.CollectionName
	}
	o.Model = model
	o.Collection = MongoClient.Database(o.DBName).Collection(defaultCollectionName)
	o.Filter = bson.D{}
	o.OrFilter = []bson.M{}
	return o
}

func (o *GODM) SetDBName(name string) *GODM {
	o.DBName = name
	return o
}

// SetCollectionName 允許覆蓋默認的集合名稱。
// SetCollectionName allows overriding the default collection name.
func (o *GODM) SetCollectionName(name string) *GODM {
	o.CollectionName = name
	if o.Collection != nil {
		o.Collection = MongoClient.Database(o.DBName).Collection(name)
	}
	return o
}
