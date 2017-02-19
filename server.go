package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/achiku/testsvr"
)

func helloHandler(logger testsvr.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := HelloResponse{
			StatusCode: SuccessStatusCode,
			Message:    "hello!!",
		}
		payload, err := json.Marshal(res)
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(payload))
		return
	}
}

// DefaultHandlerMap default url and handler map
var DefaultHandlerMap = map[string]testsvr.CreateHandler{
	"/v1/api/hello": helloHandler,
}
