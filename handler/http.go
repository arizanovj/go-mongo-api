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
func (response *Response) CORS(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers,Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	if r.Method == "OPTIONS" {
		return
	}
	next(w, r)

}
