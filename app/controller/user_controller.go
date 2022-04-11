package controller

//UserController @Controller(user)
type UserController struct {
}

//Get @GET
func (u *UserController) Get() string {
	return "hello i user"
}
