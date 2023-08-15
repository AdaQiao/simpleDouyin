package db

import (
	"database/sql"
	"fmt"
	"github.com/AdaQiao/simpleDouyin/model"
	"log"
)

type UserRepository interface {
	CreateUser(user model.UserPassword) error
	GetUser(token string) (*model.User, error)
	GetUserId(token string) (int64, error)
	UpdateWorkCount(token string) error
}

type MySQLUserRepository struct {
}

func NewMySQLUserRepository() *MySQLUserRepository {
	return &MySQLUserRepository{}
}

func (repo *MySQLUserRepository) CreateUser(user model.UserPassword) error {
	// 执行插入用户数据的SQL语句
	query := `
		INSERT INTO users (token, name, is_follow, follow_count, follower_count, total_favorited, work_count, favorite_count, avatar)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	// 执行插入操作
	token := user.Username + user.Password
	_, err := dB.Exec(query, token, user.Username, 0, 0, 0, 0, 0, 0, "https://simple-douyin.oss-cn-beijing.aliyuncs.com/douyin.png")
	if err != nil {
		log.Println("插入用户数据失败:", err)
		return err
	}

	return nil
}

func (repo *MySQLUserRepository) GetUser(token string) (*model.User, error) {
	// 执行查询用户数据的SQL语句
	query := "SELECT id, name, follow_count, follower_count, is_follow ,total_favorited, work_count, favorite_count FROM users WHERE token = ?"
	row := dB.QueryRow(query, token)

	user := &model.User{}
	err := row.Scan(&user.Id, &user.Name, &user.FollowCount, &user.FollowerCount, &user.IsFollow, &user.TotalFavorited, &user.WorkCount, &user.FavoriteCount)
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
func (repo *MySQLUserRepository) GetUserId(token string) (int64, error) {
	// 执行查询用户数据的SQL语句
	query := "SELECT id FROM users WHERE token = ?"
	row := dB.QueryRow(query, token)

	var userId int64
	err := row.Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			// 用户不存在
			fmt.Println("用户不存在:", token)
			return userId, fmt.Errorf("用户不存在")
		}
		log.Println("查询用户失败:", err)
		return -1, err
	}
	return userId, nil
}

func (repo *MySQLUserRepository) UpdateWorkCount(token string) error {
	// 获取用户信息
	query := "SELECT work_count FROM users WHERE token = ?"
	row := dB.QueryRow(query, token)
	var workCount int64
	err := row.Scan(&workCount)
	// 更新 work_count
	workCount++
	query = "UPDATE users SET work_count = ? WHERE token = ?"
	_, err = dB.Exec(query, workCount, token)
	if err != nil {
		log.Println("更新用户 work_count 失败:", err)
		return err
	}

	return nil
}
