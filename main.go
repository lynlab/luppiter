package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/graphql-go/graphql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"luppiter/services"
)

type Template struct{ templates *template.Template }

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func generateContext(c echo.Context) context.Context {
	ctx := context.Background()
	headers := c.Request().Header["Authorization"]
	if len(headers) == 1 {
		splits := strings.Split(headers[0], " ")
		ctx = context.WithValue(ctx, "Authorization", splits[len(splits)-1])
	} else {
		ctx = context.WithValue(ctx, "Authorization", "")
	}

	key := c.QueryParam("key")
	if key != "" {
		ctx = context.WithValue(ctx, "APIKey", key)
	} else {
		ctx = context.WithValue(ctx, "APIKey", "")
	}

	return ctx
}

func main() {
	e := echo.New()

	/// GET, POST /apis/graphql
	/// Public GraphQL APIs endpoint.
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				// Auth query.
				"apiKeyList": services.APIKeysQuery,

				// KeyValue query.
				"keyValueItem": services.KeyValueItemQuery,

				// Storage query.
				"storageBucketList": services.StorageBucketListQuery,
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootMutation",
			Fields: graphql.Fields{
				// Auth mutations.
				"createAPIKey":             services.CreateAPIKeyMutation,
				"addPermissionToAPIKey":    services.AddPermissionToAPIKeyMutation,
				"removePermissionToAPIKey": services.RemovePermissionToAPIKeyMutation,

				// KeyValue mutation.
				"setKeyValueItem": services.SetKeyValueItemMutation,

				// Storage mutation.
				"createStorageBucket": services.CreateStorageBucketMutation,
			},
		}),
	})

	e.Any("/apis/graphql", func(c echo.Context) error {
		var req string
		if c.Request().Method == "GET" {
			req = c.QueryParam("query")
		} else if c.Request().Method == "POST" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(c.Request().Body)
			req = buf.String()
		}

		// Add authorization data.

		// Run query and return the result.
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: req,
			Context:       generateContext(c),
		})

		return c.JSON(http.StatusOK, result)
	})

	/// File apis.
	e.GET("/files/:bucketName/:itemName", func(c echo.Context) error {
		// Cache control
		c.Response().Header().Set("Cache-Control", "max-age=86400, public, immutable")
		if headers := c.Request().Header["Cache-Control"]; len(headers) == 1 {
			cacheControl := strings.Split(headers[0], "=")
			if len(cacheControl) == 2 && cacheControl[0] == "max-age" {
				age, _ := strconv.Atoi(cacheControl[1])
				if age < 86400 {
					return c.String(http.StatusNotModified, "")
				}
			}
		}

		// Download file from cloud storage and return it.
		file, cType, err := services.DownloadStorageItem(generateContext(c), c.Param("bucketName"), c.Param("itemName"))
		if err != nil {
			return c.String(http.StatusNotFound, "")
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			return c.String(http.StatusInternalServerError, "")
		}

		return c.Blob(http.StatusOK, cType, data)
	})

	e.POST("/files/:bucketName/:itemName", func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.String(http.StatusBadRequest, "")
		}
		src, err := file.Open()
		if err != nil {
			return c.String(http.StatusBadRequest, "")
		}

		err = services.UploadStorageItem(generateContext(c), c.Param("bucketName"), c.Param("itemName"), src)
		if err != nil {
			return c.String(http.StatusBadRequest, "")
		}
		return c.String(http.StatusCreated, "")
	})

	/// HTTP web pages.
	/// GET /web/**
	e.Renderer = &Template{templates: template.Must(template.ParseGlob("public/*.html"))}
	e.GET("/web", func(c echo.Context) error { return c.Render(http.StatusOK, "index", nil) })

	/// Static files serving.
	e.Static("/statics", "public/statics")

	/// Set middlewares
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://127.0.0.1",
			"http://luppiter.lynlab.co.kr",
			"https://luppiter.lynlab.co.kr",
		},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Cache-Control"},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	e.Use(middleware.BodyLimit("10M"))

	/// ... and server starts!
	e.Logger.Fatal(e.Start(":1323"))
}
