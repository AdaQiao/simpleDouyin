package db

import "github.com/cpl/simple-demo/controller"

type UserRepository interface {
	CreateUser(user controller.UserPassword) error
	GetUser(token string) (*controller.User, error)
}
