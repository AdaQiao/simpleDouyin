package db

import (
	"database/sql"
	"github.com/cpl/simple-demo/controller"
	"log"
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
	query := `
		INSERT INTO users (token, name, is_follow)
		VALUES (?, ?, ?)
	`
	// 执行插入操作
	token := user.Username + user.Password
	_, err := dB.Exec(query, token, user.Username, 0)
	if err != nil {
		log.Println("插入用户数据失败:", err)
		return err
	}

	return nil
}

func (repo *MySQLUserRepository) GetUser(token string) (*controller.User, error) {
	// 执行查询用户数据的SQL语句
	query := "SELECT id, name, follow_count, follower_count, is_follow FROM users WHERE token = ?"
	row := dB.QueryRow(query, token)

	user := &controller.User{}
	err := row.Scan(&user.Id, &user.Name, &user.FollowCount, &user.FollowerCount, &user.IsFollow)
	if err != nil {
		if err == sql.ErrNoRows {
			// 用户不存在
			return nil, nil
		}
		log.Println("查询用户失败:", err)
		return nil, err
	}

	return user, nil
}
