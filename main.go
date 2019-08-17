package main

import (
	"context"
	"net/http"
	"time"

	"github.com/arizanovj/go-mongo-api/auth"
	"github.com/arizanovj/go-mongo-api/env"
	"github.com/arizanovj/go-mongo-api/handler"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(ctx, clientOptions)

	env := env.Env{
		MDB: client,
	}
	bearerAuth := &auth.Bearer{}

	logHandler := &handler.Log{Env: &env}
	resp := &handler.Response{}
	r := mux.NewRouter()
	r.Handle("/log/", negroni.New(
		negroni.HandlerFunc(resp.CORS),
		negroni.HandlerFunc(bearerAuth.Validate),
		negroni.Wrap(http.HandlerFunc(logHandler.Create)),
	))
	http.Handle("/", r)

	http.ListenAndServe(":9001", nil)
}
