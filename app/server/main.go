package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type errorResponse struct {
	ErrorCode string `json:"errorCode"`
}

func ping(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte("pong"))
}

func respondError(w http.ResponseWriter, err error) {
	body, _ := json.Marshal(errorResponse{fmt.Sprintf("%v", err)})

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(body)
}

func main() {
	router := httprouter.New()
	router.GET("/ping", ping)
	router.GET("/storage/:namespace/:key", getStorageItem)
	router.POST("/storage/:namespace/:key", postStorageItem)

	log.Fatal(http.ListenAndServe(":8080", router))
}
