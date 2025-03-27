package odm

import (
	"go.mongodb.org/mongo-driver/bson"
)

// Limit 設置要檢索的最大文檔數量。
// Limit sets the maximum number of documents to retrieve.
func (o *GODM) Limit(n int64) *GODM {
	o.LimitCount = n
	return o
}

// Offset 設置要跳過的文檔數量。
// Offset sets the number of documents to skip.
func (o *GODM) Offset(n int64) *GODM {
	o.SkipCount = n
	return o
}

// OrderBy 根據指定的欄位進行升序或降序排序。
// OrderBy sorts the results by the specified field in ascending or descending order.
func (o *GODM) OrderBy(field string, ascending bool) *GODM {
	order := 1
	if !ascending {
		order = -1
	}
	o.SortFields = append(o.SortFields, bson.E{Key: field, Value: order})
	return o
}

// Select 指定結果中要包含的欄位。
// Select specifies the fields to include in the results.
func (o *GODM) Select(fields ...string) *GODM {
	o.Projection = bson.M{}
	for _, field := range fields {
		o.Projection[field] = 1
	}
	return o
}

// Exclude 指定結果中要排除的欄位。
// Exclude specifies the fields to exclude from the results.
func (o *GODM) Exclude(fields ...string) *GODM {
	if o.Projection == nil {
		o.Projection = bson.M{}
	}
	for _, field := range fields {
		o.Projection[field] = 0
	}
	return o
}
