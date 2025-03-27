package odm

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// Aggregate 執行聚合管道並解碼結果。
// Aggregate runs an aggregation pipeline and decodes the results.
func (o *GODM) Aggregate(pipeline mongo.Pipeline, results interface{}) error {
	cursor, err := o.Collection.Aggregate(o.getContext(), pipeline)
	if err != nil {
		return fmt.Errorf("aggregate error: %w", err)
	}
	defer cursor.Close(o.getContext())

	return cursor.All(o.getContext(), results)
}
