package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type Response struct {
	Data    interface{}         `json:"data"`
	Message string              `json:"message"`
	Err     interface{}         `json:"error"`
	Code    int                 `json:"code"`
	W       http.ResponseWriter `json:"-"`
}

func (response *Response) Json() {
	response.W.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Printf("%+v\n", reflect.TypeOf(response.Err))
	response.W.WriteHeader(response.Code)
	json.NewEncoder(response.W).Encode(response)
}
