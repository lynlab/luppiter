package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hellodhlyn/luppiter/connectors"
	"github.com/hellodhlyn/luppiter/controllers/vulcan"
	"github.com/hellodhlyn/luppiter/repositories"
	"github.com/hellodhlyn/luppiter/services"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	db, err := connectors.NewDatabaseConnection()
	if err != nil {
		panic(err)
	}

	accountRepo, err := repositories.NewUserAccountRepository(db)
	if err != nil {
		panic(err)
	}
	identityRepo, err := repositories.NewUserIdentityRepository(db)
	if err != nil {
		panic(err)
	}
	tokenRepo, err := repositories.NewAccessTokenRepository(db)
	if err != nil {
		panic(err)
	}
	appRepo, err := repositories.NewApplicationRepository(db)
	if err != nil {
		panic(err)
	}

	accountSvc, err := services.NewUserAccountService(accountRepo, identityRepo)
	if err != nil {
		panic(err)
	}
	tokenSvc, err := services.NewAccessTokenService(tokenRepo)
	if err != nil {
		panic(err)
	}
	appSvc, err := services.NewApplicationService(appRepo)
	if err != nil {
		panic(err)
	}
	authSvc, err := services.NewAuthenticationService(tokenRepo)
	if err != nil {
		panic(err)
	}

	appCtrl, _ := vulcan.NewApplicationsController(appSvc)
	authCtrl, _ := vulcan.NewAuthController(accountSvc, appSvc, tokenSvc, authSvc)

	router := httprouter.New()
	router.GET("/ping", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		_, _ = w.Write([]byte("pong"))
	})

	router.GET("/vulcan/applications/:uuid", appCtrl.Get)
	router.GET("/vulcan/auth/me", authCtrl.GetMe)
	router.POST("/vulcan/auth/signin/google", authCtrl.AuthByGoogle)
	router.POST("/vulcan/auth/activate", authCtrl.ActivateAccessToken)

	origins := strings.Split(os.Getenv("LUPPITER_ALLOWED_ORIGINS"), ",")
	handler := cors.New(cors.Options{
		AllowedOrigins: origins,
		AllowedHeaders: []string{"*"},
	}).Handler(router)

	fmt.Println("Start and listening 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
