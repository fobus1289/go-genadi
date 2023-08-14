package service

import "github.com/fobus1289/go-genadi/example/config"

//MainService @Service
type MainService struct {
}

func (s *MainService) AfterCreate(service *MainService, options *config.DatabaseOptions) {

}

func (m *MainService) MainService() {

}

type User struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Salary float64 `json:"salary"`
}

type Users []*User

var users = Users{
	{
		Id:     1,
		Name:   "user 1",
		Age:    18,
		Salary: 5448484848,
	},
	{
		Id:     2,
		Name:   "user 2",
		Age:    25,
		Salary: 41151515,
	},
	{
		Id:     3,
		Name:   "user 3",
		Age:    35,
		Salary: 5559595,
	},
	{
		Id:     4,
		Name:   "user 4",
		Age:    29,
		Salary: 42112456,
	},
}

func (s *MainService) FindOne(id int) *User {

	for _, user := range users {
		if user.Id == id {
			return user
		}
	}

	return nil
}

func (s *MainService) FindMany() Users {
	return users
}

func (s *MainService) Create(user *User) *User {

	var lastUser = users[len(users)-1]

	user.Id = lastUser.Id + 1

	users = append(users, user)

	return user
}

func (s *MainService) Update(user *User) *User {

	var findUser *User

	for _, usr := range users {
		if usr.Id == user.Id {
			findUser = user
			break
		}
	}

	return findUser
}

func (s *MainService) Delete(id int) bool {

	var index = -1

	for i, user := range users {
		if user.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		return false
	}

	users = append(users[:index], users[index+1:]...)

	return true
}
