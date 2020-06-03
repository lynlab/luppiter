package vulcan

import (
	"net/http"
	"time"

	"github.com/hellodhlyn/luppiter/controller"
	"github.com/hellodhlyn/luppiter/service"
	"github.com/julienschmidt/httprouter"
)

type ApplicationsController interface {
	Get(http.ResponseWriter, *http.Request, httprouter.Params)
}

type ApplicationsControllerImpl struct {
	svc service.ApplicationService
}

func NewApplicationsController(appSvc service.ApplicationService) (ApplicationsController, error) {
	return &ApplicationsControllerImpl{appSvc}, nil
}

type ApplicationBody struct {
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"createdAt"`
}

// GET /vulcan/applications/:uuid
func (ctrl *ApplicationsControllerImpl) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	app := ctrl.svc.FindByUUID(p.ByName("uuid"))

	var resBody *ApplicationBody
	if app != nil {
		resBody = &ApplicationBody{UUID: app.UUID, Name: app.Name, CreatedAt: app.CreatedAt}
	}
	controller.JsonResponse(w, resBody)
}
