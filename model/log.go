package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/arizanovj/go-mongo-api/env"
)

//Log represents the log collection in db
type Log struct {
	ID          int64    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Env         *env.Env `json:"-"`
}

//Create a log
func (log *Log) Create() (*mongo.InsertOneResult, error) {
	collection := log.Env.MDB.Database("logs").Collection("log")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, log)
	return result, err
}
