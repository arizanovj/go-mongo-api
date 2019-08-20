package model

import (
	"context"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/arizanovj/go-mongo-api/env"
)

var collection = "log"

//Log represents the log collection in db
type Log struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Description  string             `json:"description,omitempty" bson:"description,omitempty"`
	Request      interface{}        `json:"request,omitempty" bson:"request,omitempty"`
	RequestTime  time.Time          `json:"request_time,omitempty" bson:"request_time,omitempty"`
	Response     interface{}        `json:"response,omitempty" bson:"response,omitempty"`
	ResponseTime time.Time          `json:"response_time,omitempty" bson:"response_time,omitempty"`
	Type         string             `json:"type,omitempty" bson:"type,omitempty"`
	Browser      string             `json:"browser,omitempty" bson:"browser,omitempty"`
	IP           net.IP             `json:"IP,omitempty" bson:"IP,omitempty"`
	UserID       uint32             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Env          *env.Env           `json:"-" bson:"-"`
}

//Create a log
func (log *Log) Create() (*mongo.InsertOneResult, error) {
	collection := log.Env.MDB.Database(log.Env.DBName).Collection(collection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, log)
	return result, err
}

func (log *Log) Validate() bool {
	return true
}
