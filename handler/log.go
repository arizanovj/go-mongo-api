package handler

import (
	"encoding/json"
	"net/http"

	"github.com/arizanovj/go-mongo-api/env"
	"github.com/arizanovj/go-mongo-api/model"
)

type Log struct {
	Env *env.Env
}

func (a *Log) Create(w http.ResponseWriter, r *http.Request) {

	response := &Response{W: w}
	log := &model.Log{Env: a.Env}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&log)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	data, err := log.Create()

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	response.Data = data
	response.Code = 200
	response.Json()

}
