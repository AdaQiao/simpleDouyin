package db

import (
	"database/sql"
	"fmt"
	"github.com/AdaQiao/simpleDouyin/model"
	"log"
)

type UserRepository interface {
	CreateUser(user model.UserPassword) error
	GetUserByToken(token string) (*model.User, error)
	GetUserId(token string) (int64, error)
	GetUserByUserId(userId int64) (*model.User, error)
	UpdateWorkCount(token string) error
	UpdateFavoriteCount(token string, mode int32) error
	UpdateTotalFavorited(userId int64, mode int32) error
	UpdateFollowCount(userId int64, mode int32) error
	UpdateFollowerCount(userId int64, mode int32) error
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

func (repo *MySQLUserRepository) GetUserByToken(token string) (*model.User, error) {
	// 执行查询用户数据的SQL语句
	query := "SELECT id, name, follow_count, follower_count, is_follow ,total_favorited, work_count, favorite_count, avatar FROM users WHERE token = ?"
	row := dB.QueryRow(query, token)

	user := &model.User{}
	err := row.Scan(&user.Id, &user.Name, &user.FollowCount, &user.FollowerCount, &user.IsFollow, &user.TotalFavorited, &user.WorkCount, &user.FavoriteCount, &user.Avatar)
	if err != nil {
		log.Println("repository:", err)
		if err == sql.ErrNoRows {
			// 用户不存在
			fmt.Println("用户不存在:", token)
			return nil, fmt.Errorf("用户不存在")
		}
		log.Println("查询用户失败:", err)
		return nil, err
	}
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

func (repo *MySQLUserRepository) GetUserByUserId(userId int64) (*model.User, error) {
	// 执行查询用户数据的SQL语句
	query := "SELECT id, name, follow_count, follower_count, is_follow ,total_favorited, work_count, favorite_count, avatar FROM users WHERE id = ?"
	row := dB.QueryRow(query, userId)

	user := &model.User{}
	err := row.Scan(&user.Id, &user.Name, &user.FollowCount, &user.FollowerCount, &user.IsFollow, &user.TotalFavorited, &user.WorkCount, &user.FavoriteCount, &user.Avatar)
	if err != nil {
		if err == sql.ErrNoRows {
			// 用户不存在
			fmt.Println("用户不存在:", userId)
			return nil, fmt.Errorf("用户不存在")
		}
		log.Println("查询用户失败:", err)
		return nil, err
	}
	return user, nil
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

func (repo *MySQLUserRepository) UpdateFavoriteCount(token string, mode int32) error {
	// 获取用户信息
	query := "SELECT favorite_count FROM users WHERE token = ?"
	row := dB.QueryRow(query, token)
	var FavoriteCount int64
	err := row.Scan(&FavoriteCount)
	// 更新 favorite_count
	if mode == 1 {
		FavoriteCount++
	} else {
		FavoriteCount--
	}
	query = "UPDATE users SET favorite_count = ? WHERE token = ?"
	_, err = dB.Exec(query, FavoriteCount, token)
	if err != nil {
		log.Println("更新用户 favorite_count  失败:", err)
		return err
	}

	return nil
}
func (repo *MySQLUserRepository) UpdateTotalFavorited(userId int64, mode int32) error {
	query := "SELECT total_favorited FROM users WHERE id = ?"
	row := dB.QueryRow(query, userId)
	var favorited int64
	err := row.Scan(&favorited)
	// 更新 favorite_count
	if mode == 1 {
		favorited++
	} else {
		favorited--
	}
	query = "UPDATE users SET total_favorited = ? WHERE id = ?"
	_, err = dB.Exec(query, favorited, userId)
	if err != nil {
		log.Println("更新用户 total_favorited 失败:", err)
		return err
	}

	return nil
}

func (repo *MySQLUserRepository) UpdateFollowCount(userId int64, mode int32) error {
	query := "SELECT follow_count FROM users WHERE id = ?"
	row := dB.QueryRow(query, userId)
	var followCount int64
	err := row.Scan(&followCount)
	// 更新 follow_count
	if mode == 1 {
		followCount++
	} else if mode == 2 {
		followCount--
	}
	query = "UPDATE users SET follow_count = ? WHERE id = ?"
	_, err = dB.Exec(query, followCount, userId)
	if err != nil {
		log.Println("更新用户 follow_count 失败:", err)
		return err
	}
	return nil
}
func (repo *MySQLUserRepository) UpdateFollowerCount(userId int64, mode int32) error {
	query := "SELECT follower_count FROM users WHERE id = ?"
	row := dB.QueryRow(query, userId)
	var followerCount int64
	err := row.Scan(&followerCount)
	// 更新 follower_count
	if mode == 1 {
		followerCount++
	} else if mode == 2 {
		followerCount--
	}
	query = "UPDATE users SET follower_count = ? WHERE id = ?"
	_, err = dB.Exec(query, followerCount, userId)
	if err != nil {
		log.Println("更新用户 follower_count 失败:", err)
		return err
	}
	return nil
}
