package controller

import "github.com/fobus1289/go-genadi/example/app/service"

//UserController @Controller(user)
type UserController struct {
}

//Get @GET
func (u *UserController) Get(mainService *service.PostService) string {
	return "hello i user"
}
