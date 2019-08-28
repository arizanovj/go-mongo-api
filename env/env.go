package env

import "go.mongodb.org/mongo-driver/mongo"

//Env is app env struct
type Env struct {
	MDB    *mongo.Client
	DBName string
}
