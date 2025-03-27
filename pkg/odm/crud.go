package odm

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create inserts the current model as a document into the collection.
func (o *GODM) Create() error {
	_, err := o.Collection.InsertOne(o.getContext(), o.Model)
	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}
	return nil
}

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

// First retrieves the first document matching the filter.
func (o *GODM) First() error {
	findOptions := options.FindOne()
	if o.Projection != nil {
		findOptions.SetProjection(o.Projection)
	}
	return o.Collection.FindOne(o.getContext(), o.buildFinalFilter(), findOptions).Decode(o.Model)
}

// Update applies the updates to the first document matching the filter.
func (o *GODM) Update(updates bson.M) error {
	_, err := o.Collection.UpdateOne(o.getContext(), o.buildFinalFilter(), bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}
	return nil
}

// Delete removes the first document matching the filter.
func (o *GODM) Delete() error {
	_, err := o.Collection.DeleteOne(o.getContext(), o.buildFinalFilter())
	if err != nil {
		return fmt.Errorf("delete error: %w", err)
	}
	return nil
}

// Count returns the number of documents matching the filter.
func (o *GODM) Count() (int64, error) {
	count, err := o.Collection.CountDocuments(o.getContext(), o.buildFinalFilter())
	if err != nil {
		return 0, fmt.Errorf("count error: %w", err)
	}
	return count, nil
}

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
