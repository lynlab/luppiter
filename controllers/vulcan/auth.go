package vulcan

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hellodhlyn/luppiter/controllers"
	"github.com/hellodhlyn/luppiter/services"
	"github.com/julienschmidt/httprouter"
)

type AuthController interface {
	AuthByGoogle(http.ResponseWriter, *http.Request, httprouter.Params)
	ActivateAccessToken(http.ResponseWriter, *http.Request, httprouter.Params)
	GetMe(http.ResponseWriter, *http.Request, httprouter.Params)
}

type AuthControllerImpl struct {
	accountSvc services.UserAccountService
	appSvc     services.ApplicationService
	tokenSvc   services.AccessTokenService
	authSvc    services.AuthenticationService
}

func NewAuthController(
	accountSvc services.UserAccountService,
	appSvc services.ApplicationService,
	tokenSvc services.AccessTokenService,
	authSvc services.AuthenticationService,
) (AuthController, error) {
	return &AuthControllerImpl{accountSvc, appSvc, tokenSvc, authSvc}, nil
}

type MeResBody struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type SignInReqBody struct {
	IDToken string `json:"idToken"`
	AppID   string `json:"appId"`
}

type SignInResBody struct {
	ActivationKey string `json:"activationKey"`
}

type ActivateReqBody struct {
	ActivationToken string `json:"activationToken"`
}

type ActivateResBody struct {
	AccessKey string     `json:"accessKey"`
	SecretKey string     `json:"secretKey"`
	ExpireAt  *time.Time `json:"expireAt"`
}

// GET /vulcan/auth/me
func (ctrl *AuthControllerImpl) GetMe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := ctrl.authSvc.Authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	controllers.JsonResponse(w, &MeResBody{UUID: user.UUID, Email: user.Email, Username: user.Username})
}

// POST /vulcan/auth/signin/google
func (ctrl *AuthControllerImpl) AuthByGoogle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var reqBody SignInReqBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account, err := ctrl.accountSvc.FindOrCreateByGoogleAccount(reqBody.IDToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app := ctrl.appSvc.FindByUUID(reqBody.AppID)
	if app == nil {
		http.Error(w, "invalid appId", http.StatusBadRequest)
		return
	}

	token, _ := ctrl.tokenSvc.CreateAccessToken(&account.Identity, app)
	controllers.JsonResponse(w, &SignInResBody{ActivationKey: token.ActivationKey})
}

// POST /vulcan/auth/activate
func (ctrl *AuthControllerImpl) ActivateAccessToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var reqBody ActivateReqBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := ctrl.tokenSvc.ActivateAccessToken(reqBody.ActivationToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	controllers.JsonResponse(w, &ActivateResBody{AccessKey: token.AccessKey, SecretKey: token.SecretKey, ExpireAt: token.ExpireAt})
}
