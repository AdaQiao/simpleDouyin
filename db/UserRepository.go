package db

import (
	"github.com/cpl/simple-demo/controller"
)

type UserRepository interface {
	CreateUser(user controller.UserPassword) error
	GetUser(token string) (*controller.User, error)
}

type MySQLUserRepository struct {
}

func NewMySQLUserRepository() *MySQLUserRepository {
	return &MySQLUserRepository{}
}

func (repo *MySQLUserRepository) CreateUser(user controller.UserPassword) error {
	// 执行插入用户数据的SQL语句
	// ...
	return nil
}

func (repo *MySQLUserRepository) GetUser(token string) (*controller.User, error) {
	// 执行查询用户数据的SQL语句
	// ...
	return nil, nil
}
