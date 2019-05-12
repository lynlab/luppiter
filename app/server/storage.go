package main

import (
	"errors"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/lynlab/luppiter/services/storage"
)

func getStorageItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	reader, contentType, err := storage.ReadFile(p.ByName("namespace"), p.ByName("key"))
	if err != nil {
		respondError(w, err)
		return
	}

	w.Header().Set("Content-Type", contentType)
	io.Copy(w, reader)
}

func postStorageItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	upload, _, err := r.FormFile("file")
	if err != nil {
		respondError(w, errors.New("bad request"))
		return
	}

	err = storage.WriteFile(p.ByName("namespace"), p.ByName("key"), upload)
	if err != nil {
		respondError(w, err)
		return
	}
}
