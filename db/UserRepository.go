package db

import (
	"database/sql"
	"fmt"
	"github.com/RaymondCode/simple-demo/controller"
	"strings"

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
		INSERT INTO users (token, name, is_follow,follow_count, follower_count)
		VALUES (?, ?, ?, ?, ?)
	`
	// 执行插入操作
	token := user.Username + user.Password
	_, err := dB.Exec(query, token, user.Username, 0, 0, 0)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			// 处理唯一约束错误
			log.Printf("用户名 %s 已存在\n", user.Username)
			return fmt.Errorf("用户名 %s 已存在", user.Username)
		}
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
			fmt.Println("用户不存在:", token)
			return nil, fmt.Errorf("用户不存在")
		}
		log.Println("查询用户失败:", err)
		return nil, err
	}
	fmt.Printf("查询结果 - ID: %d, Name: %s, FollowCount: %d, FollowerCount: %d, IsFollow: %t\n", user.Id, user.Name, user.FollowCount, user.FollowerCount, user.IsFollow)
	return user, nil
}
