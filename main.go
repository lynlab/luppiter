package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
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

func main() {
	e := echo.New()

	/// GET, POST /apis/graphql
	/// Public GraphQL APIs endpoint.
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				// Auth queries.
				"apiKeyList": services.APIKeysQuery,
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootMutation",
			Fields: graphql.Fields{
				// Auth mutations.
				"createAPIKey":             services.CreateAPIKeyMutation,
				"addPermissionToAPIKey":    services.AddPermissionToAPIKeyMutation,
				"removePermissionToAPIKey": services.RemovePermissionToAPIKeyMutation,
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

		// Run query and return the result.
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: req,
			Context:       ctx,
		})

		return c.JSON(http.StatusOK, result)
	})

	/// HTTP web pages.
	/// GET /web/**
	e.Renderer = &Template{templates: template.Must(template.ParseGlob("public/*.html"))}
	e.GET("/web", func(c echo.Context) error { return c.Render(http.StatusOK, "index", nil) })

	/// Static files serving.
	e.Static("/statics", "public/statics")

	/// Set middlewares
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://127.0.0.1", "http://luppiter.lynlab.co.kr", "https://luppiter.lynlab.co.kr"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	/// ... and server starts!
	e.Logger.Fatal(e.Start(":1323"))
}
