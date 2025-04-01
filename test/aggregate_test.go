package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"godm/pkg/odm"
)

func TestGODM_Aggregate(t *testing.T) {
	g := &odm.GODM{}
	g.Where("type", "=", "admin")
	g.OrWhere("status", "=", "active")
	g.OrWhere("status", "=", "pending")

	actual := g.ToBson()
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

	assert.Equal(t, expected, actual)
}
