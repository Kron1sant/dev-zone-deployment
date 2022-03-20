package mongodb

import (
	"devZoneDeployment/db"

	"go.mongodb.org/mongo-driver/bson"
)

type mongoFilter struct {
	Conditions []filterCondition
}

type filterCondition struct {
	Type  condType
	Key   string
	Value interface{}
}

type condType int

const (
	eqType condType = iota
	nEqType
)

func (f *mongoFilter) AddEq(field string, value interface{}) db.Filter {
	f.Conditions = append(f.Conditions, filterCondition{
		eqType,
		field,
		value,
	})

	return f
}

func (f *mongoFilter) AddNEq(field string, value interface{}) db.Filter {
	f.Conditions = append(f.Conditions, filterCondition{
		nEqType,
		field,
		value,
	})

	return f
}

func (f *mongoFilter) Compose() interface{} {
	d := bson.D{}
	for _, v := range f.Conditions {
		if v.Type == eqType {
			d = append(d, bson.E{Key: v.Key, Value: v.Value})
		} else if v.Type == nEqType {
			d = append(d, bson.E{Key: v.Key, Value: bson.M{"$ne": v.Value}})
		}
	}
	return d
}

func (f *mongoFilter) Empty() bool {
	return len(f.Conditions) == 0
}
