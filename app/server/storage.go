package main

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/lynlab/luppiter/services/storage"
)

func getStorageItem(c echo.Context) error {
	reader, contentType, err := storage.ReadFile(c.Param("namespace"), c.Param("key"))
	if err != nil {
		return respondError(c, err)
	}

	return c.Stream(http.StatusOK, contentType, reader)
}

func postStorageItem(c echo.Context) error {
	upload, err := c.FormFile("file")
	if err != nil {
		return respondError(c, errors.New("bad request"))
	}
	src, err := upload.Open()
	if err != nil {
		return respondError(c, errors.New("bad request"))
	}

	err = storage.WriteFile(c.Param("namespace"), c.Param("key"), src)
	if err != nil {
		return respondError(c, err)
	}
	return c.String(http.StatusCreated, "")
}
