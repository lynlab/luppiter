package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	keyvalue "github.com/lynlab/luppiter/services/key_value"
)

type keyValueResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func getKeyValueItem(c echo.Context) error {
	key := c.Param("key")
	value, err := keyvalue.GetKeyValueItem("public", key)
	if err != nil {
		return respondError(c, err)
	}
	return c.JSON(http.StatusOK, keyValueResponse{Key: key, Value: value})
}

func postKeyValueItem(c echo.Context) error {
	key := c.Param("key")
	value := c.QueryParam("value")
	keyvalue.SetKeyValueItem("public", key, value)
	return c.JSON(http.StatusOK, keyValueResponse{Key: key, Value: value})
}
