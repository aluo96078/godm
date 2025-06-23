package odm

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// hasObjectIDField 檢查模型是否包含 _id 欄位且型別為 primitive.ObjectID。
// hasObjectIDField checks if the model has a _id field of type primitive.ObjectID.
func (o *GODM) hasObjectIDField() bool {
	typ := reflect.TypeOf(o.Model).Elem()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if strings.HasPrefix(field.Tag.Get("bson"), "_id") && field.Type == reflect.TypeOf(primitive.ObjectID{}) {
			return true
		}
	}
	return false
}

// parseObjectID 嘗試將輸入轉換為 primitive.ObjectID。
// parseObjectID attempts to convert input into a primitive.ObjectID.
func parseObjectID(id interface{}) (primitive.ObjectID, error) {
	switch v := id.(type) {
	case string:
		return primitive.ObjectIDFromHex(v)
	case primitive.ObjectID:
		return v, nil
	default:
		return primitive.NilObjectID, fmt.Errorf("unsupported id type: %T", id)
	}
}
