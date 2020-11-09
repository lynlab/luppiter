package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/hellodhlyn/luppiter/models"
)

func Authorized(w http.ResponseWriter, r *http.Request, fn func(*models.UserIdentity)) {

}

func JsonResponse(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json; encode=utf-8")
	_ = json.NewEncoder(w).Encode(res)
}
