package odm

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// ODM 提供了一個簡單的物件文件映射工具，支援自定義上下文、自定義集合名稱、資料庫名稱、AND/OR 條件組合、批量操作以及事務控制。
// ODM provides a simple Object-Document Mapping with support for custom context, custom collection names, database name, combined AND/OR filtering, bulk operations, and transactions.

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

// Use 設置模型和集合。如果通過 SetCollectionName 提供了自定義集合名稱，則使用該名稱；否則，默認使用模型名的小寫並附加 "s"。
// Use sets the model and collection. If a custom collection name is provided via SetCollectionName, it will use that; otherwise, it defaults to the lowercase model name with an appended "s".
func (o *GODM) Use(model interface{}) *GODM {
	modelType := reflect.TypeOf(model).Elem()
	defaultCollectionName := strings.ToLower(modelType.Name()) + "s"
	if o.CollectionName != "" {
		defaultCollectionName = o.CollectionName
	}
	// Determine the database name
	dbName := "your_db"
	if o.DbName != "" {
		dbName = o.DbName
	}
	o.Model = model
	o.Collection = MongoClient.Database(dbName).Collection(defaultCollectionName)
	o.Filter = bson.D{}
	o.OrFilter = []bson.M{}
	return o
}

// WithContext 設置 ODM 操作的自定義上下文。
// WithContext sets a custom context for the ODM operations.
func (o *GODM) WithContext(ctx context.Context) *GODM {
	o.Ctx = ctx
	return o
}

// getContext 返回設置的上下文，如果未設置則默認返回 context.TODO()。
// getContext returns the set context or defaults to context.TODO().
func (o *GODM) getContext() context.Context {
	if o.Ctx != nil {
		return o.Ctx
	}
	return context.TODO()
}

// SetCollectionName 允許覆蓋默認的集合名稱。
// SetCollectionName allows overriding the default collection name.
func (o *GODM) SetCollectionName(name string) *GODM {
	o.CollectionName = name
	if o.Collection != nil {
		o.Collection = MongoClient.Database("your_db").Collection(name)
	}
	return o
}

// Where 為過濾條件添加一個 AND 條件。
// Where adds an AND condition to the filter.
func (o *GODM) Where(field, op string, value interface{}) *GODM {
	var cond bson.E
	switch op {
	case "=":
		cond = bson.E{Key: field, Value: value}
	case ">":
		cond = bson.E{Key: field, Value: bson.M{"$gt": value}}
	case "<":
		cond = bson.E{Key: field, Value: bson.M{"$lt": value}}
	case "!=":
		cond = bson.E{Key: field, Value: bson.M{"$ne": value}}
	default:
		cond = bson.E{Key: field, Value: value}
	}
	o.Filter = append(o.Filter, cond)
	return o
}

// WhereIn 為過濾條件添加一個包含條件 (AND)。
// WhereIn adds an AND condition for inclusion.
func (o *GODM) WhereIn(field string, values []interface{}) *GODM {
	cond := bson.E{Key: field, Value: bson.M{"$in": values}}
	o.Filter = append(o.Filter, cond)
	return o
}

// WhereNotIn 為過濾條件添加一個排除條件 (AND)。
// WhereNotIn adds an AND condition for exclusion.
func (o *GODM) WhereNotIn(field string, values []interface{}) *GODM {
	cond := bson.E{Key: field, Value: bson.M{"$nin": values}}
	o.Filter = append(o.Filter, cond)
	return o
}

// OrWhere 添加一個 OR 條件。多次調用將累積到 OR 條件列表中。
// OrWhere appends an OR condition. Multiple calls will accumulate conditions in the OR list.
func (o *GODM) OrWhere(field, op string, value interface{}) *GODM {
	var cond bson.M
	switch op {
	case "=":
		cond = bson.M{field: value}
	case ">":
		cond = bson.M{field: bson.M{"$gt": value}}
	case "<":
		cond = bson.M{field: bson.M{"$lt": value}}
	case "!=":
		cond = bson.M{field: bson.M{"$ne": value}}
	default:
		cond = bson.M{field: value}
	}
	o.OrFilter = append(o.OrFilter, cond)
	return o
}

// OrWhereIn 添加一個包含條件的 OR 條件。
// OrWhereIn appends an OR condition for inclusion.
func (o *GODM) OrWhereIn(field string, values []interface{}) *GODM {
	cond := bson.M{field: bson.M{"$in": values}}
	o.OrFilter = append(o.OrFilter, cond)
	return o
}

// OrWhereNotIn 添加一個排除條件的 OR 條件。
// OrWhereNotIn appends an OR condition for exclusion.
func (o *GODM) OrWhereNotIn(field string, values []interface{}) *GODM {
	cond := bson.M{field: bson.M{"$nin": values}}
	o.OrFilter = append(o.OrFilter, cond)
	return o
}

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

// FilterToMap 將 AND 條件 (bson.D) 轉換為 bson.M 映射。
// FilterToMap converts the AND filter (bson.D) to a bson.M map.
func (o *GODM) FilterToMap() map[string]interface{} {
	m := make(map[string]interface{})
	for _, e := range o.Filter {
		m[e.Key] = e.Value
	}
	return m
}

// buildFinalFilter 將 AND 和 OR 條件組合為一個最終的過濾器。
// buildFinalFilter combines the AND and OR conditions into a single filter.
func (o *GODM) buildFinalFilter() bson.D {
	if len(o.Filter) > 0 && len(o.OrFilter) > 0 {
		return bson.D{
			{
				Key: "$and",
				Value: []bson.M{
					bson.M(o.FilterToMap()),
					{
						"$or": o.OrFilter,
					},
				},
			},
		}
	} else if len(o.OrFilter) > 0 {
		return bson.D{
			{
				Key:   "$or",
				Value: o.OrFilter,
			},
		}
	}
	return o.Filter
}

// First 獲取符合條件的第一個文檔。
// First retrieves the first document matching the filter.
func (o *GODM) First() error {
	findOptions := options.FindOne()
	if o.Projection != nil {
		findOptions.SetProjection(o.Projection)
	}
	return o.Collection.FindOne(o.getContext(), o.buildFinalFilter(), findOptions).Decode(o.Model)
}

// Create 將當前模型作為文檔插入到集合中。
// Create inserts the current model as a document into the collection.
func (o *GODM) Create() error {
	_, err := o.Collection.InsertOne(o.getContext(), o.Model)
	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}
	return nil
}

// BulkCreate 將多個文檔插入到集合中。
// BulkCreate inserts multiple documents into the collection.
func (o *GODM) BulkCreate(models []interface{}) error {
	if len(models) == 0 {
		return nil
	}
	_, err := o.Collection.InsertMany(o.getContext(), models)
	if err != nil {
		return fmt.Errorf("bulk create error: %w", err)
	}
	return nil
}

// Update 對符合條件的第一個文檔應用更新操作。
// Update applies the updates to the first document matching the filter.
func (o *GODM) Update(updates bson.M) error {
	_, err := o.Collection.UpdateOne(o.getContext(), o.buildFinalFilter(), bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}
	return nil
}

// Delete 刪除符合條件的第一個文檔。
// Delete removes the first document matching the filter.
func (o *GODM) Delete() error {
	_, err := o.Collection.DeleteOne(o.getContext(), o.buildFinalFilter())
	if err != nil {
		return fmt.Errorf("delete error: %w", err)
	}
	return nil
}

// Count 返回符合條件的文檔數量。
// Count returns the number of documents matching the filter.
func (o *GODM) Count() (int64, error) {
	count, err := o.Collection.CountDocuments(o.getContext(), o.buildFinalFilter())
	if err != nil {
		return 0, fmt.Errorf("count error: %w", err)
	}
	return count, nil
}

// All 獲取所有符合條件的文檔。
// All retrieves all documents matching the filter.
func (o *GODM) All(results interface{}) error {
	findOptions := options.Find()
	if o.LimitCount > 0 {
		findOptions.SetLimit(o.LimitCount)
	}
	if len(o.SortFields) > 0 {
		findOptions.SetSort(o.SortFields)
	}
	if o.SkipCount > 0 {
		findOptions.SetSkip(o.SkipCount)
	}
	if o.Projection != nil {
		findOptions.SetProjection(o.Projection)
	}
	cursor, err := o.Collection.Find(o.getContext(), o.buildFinalFilter(), findOptions)
	if err != nil {
		return fmt.Errorf("find error: %w", err)
	}
	defer cursor.Close(o.getContext())

	return cursor.All(o.getContext(), results)
}

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

// WithTransaction 為需要原子性操作的業務提供事務支持。
// WithTransaction provides transaction support for operations that require atomicity.
func (o *GODM) WithTransaction(callback func(sessCtx mongo.SessionContext) error) error {
	session, err := MongoClient.StartSession()
	if err != nil {
		return fmt.Errorf("start session error: %w", err)
	}
	defer session.EndSession(o.getContext())

	err = mongo.WithSession(o.getContext(), session, func(sessCtx mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return fmt.Errorf("start transaction error: %w", err)
		}
		if err := callback(sessCtx); err != nil {
			if abortErr := session.AbortTransaction(sessCtx); abortErr != nil {
				return fmt.Errorf("abort transaction error: %w", abortErr)
			}
			return fmt.Errorf("transaction error: %w", err)
		}
		if err := session.CommitTransaction(sessCtx); err != nil {
			return fmt.Errorf("commit transaction error: %w", err)
		}
		return nil
	})
	return err
}
