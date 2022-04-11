package controller

import (
	"ast_set/app/service"
	"math/rand"
	"net/http"
	"strconv"
)

//MainController @Controller
type MainController struct {
	service *service.MainService
}

func (c *MainController) AfterCreate(service *service.MainService) {
	c.service = service
}

//Find @GET(:id)
func (c *MainController) Find(r *http.Request) interface{} {
	if id, err := strconv.Atoi(r.URL.Query().Get("id")); err == nil {
		return c.service.FindOne(id)
	}
	return nil
}

//FindMore @GET
func (c *MainController) FindMore() interface{} {
	return c.service.FindMany()
}

//Create @POST
func (c *MainController) Create(r *http.Request) interface{} {
	return c.service.Create(&service.User{
		Name:   r.URL.Query().Get("name"),
		Age:    rand.Int(),
		Salary: float64(rand.Int()),
	})
}

//Update @PUT
func (c *MainController) Update(r *http.Request) interface{} {
	return c.service.Create(&service.User{
		Name:   r.URL.Query().Get("name"),
		Age:    rand.Int(),
		Salary: float64(rand.Int()),
	})
}

//Delete @DELETE
func (c *MainController) Delete(r *http.Request) interface{} {
	if id, err := strconv.Atoi(r.URL.Query().Get("id")); err == nil {
		return c.service.Delete(id)
	}
	return nil
}
