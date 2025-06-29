package odm

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// crud.go - 封裝對 MongoDB 的基本操作（Create、Read、Update、Delete）與 Observer 整合
// Encapsulates basic MongoDB operations (Create, Read, Update, Delete) with integrated observer support.

// Create inserts the current model as a document into the collection.
// 創建將當前模型作為文檔插入集合中。
func (o *GODM) Create() error {
	if m, ok := o.Model.(ObservedModel); ok {
		o.Observers = append(o.Observers, m.Observers()...)
	}
	if err := o.notifyCreating(); err != nil {
		return fmt.Errorf("observer creating error: %w", err)
	}

	_, err := o.Collection.InsertOne(o.getContext(), o.Model)
	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}

	if err := o.notifyCreated(); err != nil {
		return fmt.Errorf("observer created error: %w", err)
	}
	return nil
}

// BulkCreate inserts multiple documents into the collection.
// BulkCreate 將多個文檔插入集合中。
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

// First retrieves the first document matching the filter.
// First 根據過濾條件檢索第一個文檔。
func (o *GODM) First() error {
	if len(o.WithRelations) > 0 {
		pipeline := []bson.M{
			{"$match": o.buildFinalFilter()},
		}

		for _, rel := range o.WithRelations {
			conf, ok := o.RelationConfigs[rel]
			if !ok {
				continue
			}
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         conf.From,
					"localField":   conf.LocalField,
					"foreignField": conf.ForeignField,
					"as":           conf.As,
				},
			})
			if !conf.IsArray {
				pipeline = append(pipeline, bson.M{
					"$unwind": bson.M{
						"path":                       "$" + conf.As,
						"preserveNullAndEmptyArrays": true,
					},
				})
			}
		}

		pipeline = append(pipeline, bson.M{"$limit": 1})
		cursor, err := o.Collection.Aggregate(o.getContext(), pipeline)
		if err != nil {
			return fmt.Errorf("aggregate error: %w", err)
		}
		defer cursor.Close(o.getContext())

		if cursor.Next(o.getContext()) {
			if err := cursor.Decode(o.Model); err != nil {
				return fmt.Errorf("decode error: %w (type = %T)", err, o.Model)
			}
			return nil
		}
		return mongo.ErrNoDocuments
	}

	findOptions := options.FindOne()
	if o.Projection != nil {
		findOptions.SetProjection(o.Projection)
	}

	return o.Collection.FindOne(o.getContext(), o.buildFinalFilter()).Decode(o.Model)
}

// Update applies the updates to the first document matching the filter.
// Update 將更新應用於第一個符合過濾條件的文檔。
func (o *GODM) Update(updates bson.M) error {
	if m, ok := o.Model.(ObservedModel); ok {
		o.Observers = append(o.Observers, m.Observers()...)
	}
	if err := o.notifyUpdating(); err != nil {
		return fmt.Errorf("observer updating error: %w", err)
	}

	_, err := o.Collection.UpdateOne(o.getContext(), o.buildFinalFilter(), bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	if err := o.notifyUpdated(); err != nil {
		return fmt.Errorf("observer updated error: %w", err)
	}
	return nil
}

// Delete removes the first document matching the filter.
// Delete 刪除第一個符合過濾條件的文檔。
func (o *GODM) Delete() error {
	if m, ok := o.Model.(ObservedModel); ok {
		o.Observers = append(o.Observers, m.Observers()...)
	}
	if err := o.notifyDeleting(); err != nil {
		return fmt.Errorf("observer deleting error: %w", err)
	}

	_, err := o.Collection.DeleteOne(o.getContext(), o.buildFinalFilter())
	if err != nil {
		return fmt.Errorf("delete error: %w", err)
	}

	if err := o.notifyDeleted(); err != nil {
		return fmt.Errorf("observer deleted error: %w", err)
	}
	return nil
}

// Count returns the number of documents matching the filter.
// Count 返回符合過濾條件的文檔數量。
func (o *GODM) Count() (int64, error) {
	count, err := o.Collection.CountDocuments(o.getContext(), o.buildFinalFilter())
	if err != nil {
		return 0, fmt.Errorf("count error: %w", err)
	}
	return count, nil
}

// All retrieves all documents matching the filter.
// All 根據過濾條件檢索所有文檔。
func (o *GODM) All(results interface{}) error {
	if len(o.WithRelations) > 0 {
		pipeline := []bson.M{
			{"$match": o.buildFinalFilter()},
		}

		for _, rel := range o.WithRelations {
			conf, ok := o.RelationConfigs[rel]
			if !ok {
				continue
			}
			pipeline = append(pipeline, bson.M{
				"$lookup": bson.M{
					"from":         conf.From,
					"localField":   conf.LocalField,
					"foreignField": conf.ForeignField,
					"as":           conf.As,
				},
			})
			if !conf.IsArray {
				pipeline = append(pipeline, bson.M{
					"$unwind": bson.M{
						"path":                       "$" + conf.As,
						"preserveNullAndEmptyArrays": true,
					},
				})
			}
		}

		if len(o.SortFields) > 0 {
			pipeline = append(pipeline, bson.M{"$sort": o.SortFields})
		}
		if o.SkipCount > 0 {
			pipeline = append(pipeline, bson.M{"$skip": o.SkipCount})
		}
		if o.LimitCount > 0 {
			pipeline = append(pipeline, bson.M{"$limit": o.LimitCount})
		}

		cursor, err := o.Collection.Aggregate(o.getContext(), pipeline)
		if err != nil {
			return fmt.Errorf("aggregate error: %w", err)
		}
		defer cursor.Close(o.getContext())

		return cursor.All(o.getContext(), results)
	}

	// fallback to regular Find
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

// Exists 檢查是否存在符合過濾條件的文檔。
// Exists checks if a document exists that matches the filter.
func (o *GODM) Exists() (bool, error) {
	count, err := o.Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
