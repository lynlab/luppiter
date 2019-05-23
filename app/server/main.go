package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func respondError(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, fmt.Sprintf("%v", err))
}

func main() {
	e := echo.New()

	e.GET("/ping", ping)
	e.GET("/storage/:namespace/:key", getStorageItem)
	e.POST("/storage/:namespace/:key", postStorageItem)
	e.GET("/vulcan/key_values/:key", getKeyValueItem)
	e.POST("/vulcan/key_values/:key", postKeyValueItem)

	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(os.Getenv("LUPPITER_ALLOWED_ORIGINS"), ","),
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Logger.Fatal(e.Start(":1323"))
}
