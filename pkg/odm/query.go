package odm

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

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

// WhereIn adds an AND condition for inclusion.
func (o *GODM) WhereIn(field string, values []interface{}) *GODM {
	cond := bson.E{Key: field, Value: bson.M{"$in": values}}
	o.Filter = append(o.Filter, cond)
	return o
}

// WhereNotIn adds an AND condition for exclusion.
func (o *GODM) WhereNotIn(field string, values []interface{}) *GODM {
	cond := bson.E{Key: field, Value: bson.M{"$nin": values}}
	o.Filter = append(o.Filter, cond)
	return o
}

// OrWhere appends an OR condition.
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

// OrWhereIn appends an OR condition for inclusion.
func (o *GODM) OrWhereIn(field string, values []interface{}) *GODM {
	cond := bson.M{field: bson.M{"$in": values}}
	o.OrFilter = append(o.OrFilter, cond)
	return o
}

// OrWhereNotIn appends an OR condition for exclusion.
func (o *GODM) OrWhereNotIn(field string, values []interface{}) *GODM {
	cond := bson.M{field: bson.M{"$nin": values}}
	o.OrFilter = append(o.OrFilter, cond)
	return o
}

// buildFinalFilter combines the AND and OR conditions into a single filter.
func (o *GODM) buildFinalFilter() bson.D {
	if len(o.Filter) > 0 && len(o.OrFilter) > 0 {
		return bson.D{{
			Key: "$and",
			Value: []bson.M{
				bson.M(o.FilterToMap()),
				{"$or": o.OrFilter},
			},
		}}
	} else if len(o.OrFilter) > 0 {
		return bson.D{{Key: "$or", Value: o.OrFilter}}
	}
	return o.Filter
}

// FilterToMap converts the AND filter (bson.D) to a bson.M map.
func (o *GODM) FilterToMap() map[string]interface{} {
	m := make(map[string]interface{})
	for _, e := range o.Filter {
		m[e.Key] = e.Value
	}
	return m
}

// WhereID filters by _id field. Returns error if model does not contain a primitive.ObjectID _id field.
func (o *GODM) WhereID(id interface{}) *GODM {
	if !o.hasObjectIDField() {
		log.Println("WhereID error: model does not contain a _id field of type primitive.ObjectID")
		return nil
	}
	objectID, err := parseObjectID(id)
	if err != nil {
		log.Printf("WhereID error: %v\n", err)
		return nil
	}
	return o.Where("_id", "=", objectID)
}
