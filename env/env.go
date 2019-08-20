package env

import "go.mongodb.org/mongo-driver/mongo"

type Env struct {
	MDB    *mongo.Client
	DBName string
}
