package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"godm/pkg/odm"
)

func TestGODM_ToBson_SimpleWhere(t *testing.T) {
	q := &odm.GODM{}
	q.Where("age", ">", 30)
	expected := bson.D{{Key: "age", Value: bson.M{"$gt": 30}}}
	assert.Equal(t, expected, q.ToBson())
}

func TestGODM_ToBson_WhereIn(t *testing.T) {
	q := &odm.GODM{}
	q.WhereIn("status", []interface{}{"active", "pending"})
	expected := bson.D{{Key: "status", Value: bson.M{"$in": []interface{}{"active", "pending"}}}}
	assert.Equal(t, expected, q.ToBson())
}

func TestGODM_ToBson_WhereAndOr(t *testing.T) {
	q := &odm.GODM{}
	q.Where("type", "=", "admin")
	q.OrWhere("status", "=", "active")
	q.OrWhere("status", "=", "pending")
	expected := bson.D{{
		Key: "$and",
		Value: []bson.M{
			{"type": "admin"},
			{"$or": []bson.M{
				{"status": "active"},
				{"status": "pending"},
			}},
		},
	}}
	assert.Equal(t, expected, q.ToBson())
}
