package router

import (
	"ast_set/app/controller"
	"ast_set/app/service"
	"ast_set/config"
	"encoding/json"
	"net/http"
	"regexp"
)

var (
	config_config_DatabaseOptions = config.NewDatabaseOptions()
)

var (
	app_service_service_MainService = &service.MainService{}
	app_service_service_PostService = &service.PostService{}
	app_service_service_UserService = &service.UserService{}
)

var (
	app_controller_controller_MainController = &controller.MainController{}
	app_controller_controller_UserController = &controller.UserController{}
)

var (
	pattern_method_app_controller_controller_MainController_Find_GET      = regexp.MustCompile(`^(/?([0-9a-zA-Z]+)/?)$`)
	pattern_method_app_controller_controller_MainController_FindMore_GET  = regexp.MustCompile(`^(/)$`)
	pattern_method_app_controller_controller_MainController_Create_POST   = regexp.MustCompile(`^(/)$`)
	pattern_method_app_controller_controller_MainController_Update_PUT    = regexp.MustCompile(`^(/)$`)
	pattern_method_app_controller_controller_MainController_Delete_DELETE = regexp.MustCompile(`^(/)$`)
	pattern_method_app_controller_controller_UserController_Get_GET       = regexp.MustCompile(`^(/?user/?)$`)
)

func NewServer() *Server {
	app_controller_controller_MainController.AfterCreate(app_service_service_MainService)
	app_service_service_MainService.AfterCreate(app_service_service_MainService, config_config_DatabaseOptions)

	return &Server{
		Routers: []*Router{},
	}
}

type Router struct {
	path   string
	Method string
	match  *regexp.Regexp
}

type Server struct {
	Routers []*Router
}

func (s *Server) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {

	switch request.Method {

	case http.MethodGet:

		if pattern_method_app_controller_controller_MainController_Find_GET.MatchString(request.URL.Path) {
			var result = app_controller_controller_MainController.Find(request)
			if out, err := json.Marshal(result); err == nil {
				_, _ = responseWriter.Write(out)
			}
		} else if pattern_method_app_controller_controller_MainController_FindMore_GET.MatchString(request.URL.Path) {
			var result = app_controller_controller_MainController.FindMore()
			if out, err := json.Marshal(result); err == nil {
				_, _ = responseWriter.Write(out)
			}
		} else if pattern_method_app_controller_controller_UserController_Get_GET.MatchString(request.URL.Path) {
			var result = app_controller_controller_UserController.Get()
			if out, err := json.Marshal(result); err == nil {
				_, _ = responseWriter.Write(out)
			}
		}

	case http.MethodPost:

		if pattern_method_app_controller_controller_MainController_Create_POST.MatchString(request.URL.Path) {
			var result = app_controller_controller_MainController.Create(request)
			if out, err := json.Marshal(result); err == nil {
				_, _ = responseWriter.Write(out)
			}
		}

	case http.MethodPut:

		if pattern_method_app_controller_controller_MainController_Update_PUT.MatchString(request.URL.Path) {
			var result = app_controller_controller_MainController.Update(request)
			if out, err := json.Marshal(result); err == nil {
				_, _ = responseWriter.Write(out)
			}
		}

	case http.MethodDelete:

		if pattern_method_app_controller_controller_MainController_Delete_DELETE.MatchString(request.URL.Path) {
			var result = app_controller_controller_MainController.Delete(request)
			if out, err := json.Marshal(result); err == nil {
				_, _ = responseWriter.Write(out)
			}
		}

	default:
		http.NotFound(responseWriter, request)
	}

}
