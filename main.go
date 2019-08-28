package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/user"
	"time"

	"github.com/arizanovj/go-mongo-api/auth"
	"github.com/arizanovj/go-mongo-api/env"
	"github.com/arizanovj/go-mongo-api/handler"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	configDir := usr.HomeDir + "/.config/api/"
	configName := "api"

	viper.SetConfigName(configName) // name of config file (without extension)
	viper.AddConfigPath(configDir)  // path to look for the config file in

	err = viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
	//check congfigs

	if !viper.IsSet("db.user") {
		log.Fatal("missing db user")
	}
	if !viper.IsSet("db.password") {
		log.Fatal("missing db password")
	}
	if !viper.IsSet("db.host") {
		log.Fatal("missing db host")
	}
	if !viper.IsSet("db.port") {
		log.Fatal("missing port number")
	}
	if !viper.IsSet("db.name") {
		log.Fatal("missing db name")
	}
	if !viper.IsSet("db.auth") {
		log.Fatal("missing db auth")
	}
	dbUser := viper.GetString("db.user")
	dbPassword := viper.GetString("db.password")
	dbHost := viper.GetString("db.host")
	dbPort := viper.GetString("db.port")
	dbName := viper.GetString("db.name")
	dbAuth := viper.GetString("db.auth")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	connectionURI := "mongodb://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?authSource=" + dbAuth
	clientOptions := options.Client().ApplyURI(connectionURI)
	client, _ := mongo.Connect(ctx, clientOptions)
	defer client.Disconnect(ctx)
	defer cancel()

	env := env.Env{
		MDB:    client,
		DBName: dbName,
	}
	bearerAuth := &auth.Bearer{Env: &env}
	logHandler := &handler.Log{Env: &env}
	r := mux.NewRouter()
	r.Handle("/log/", negroni.New(
		negroni.HandlerFunc(bearerAuth.Validate),
		negroni.Wrap(http.HandlerFunc(logHandler.Create)),
	))

	http.Handle("/", r)

	http.ListenAndServe(":9001", nil)
}
