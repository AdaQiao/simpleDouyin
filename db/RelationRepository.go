package db

import (
	"database/sql"
	"errors"
	"log"
	"sync"
)

type RelationRepository interface {
	AddFollower(userId int64, followerId int64) error
	RemoveFollower(userId int64, followerId int64) error
	AddFollow(userId int64, followId int64) error
	RemoveFollow(userId int64, followId int64) error
	CheckFollow(userId, followId int64) (bool, error)
	GetFollowById(userId int64) ([]int64, error)
	GetFollowerById(userId int64) ([]int64, error)
}

type MySQLRelationRepository struct {
	mutex sync.Mutex
}

func NewMySQLRelationRepository() *MySQLRelationRepository {
	return &MySQLRelationRepository{}
}

func (repo *MySQLRelationRepository) AddFollower(userId int64, followerId int64) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 检查是否已存在相同的记录
	query := "SELECT is_following FROM follower WHERE user_id = ? and follower_id = ?"
	row := dB.QueryRow(query, userId, followerId)
	var isFollowing bool
	err := row.Scan(&isFollowing)
	if err == sql.ErrNoRows {
		// 不存在记录，插入新记录
		query = "INSERT INTO follower (user_id, follower_id, is_following) VALUES (?, ?, true)"
		_, err = dB.Exec(query, userId, followerId)
		if err != nil {
			log.Println("插入粉丝记录失败:", err)
			return err
		}
	} else if err != nil {
		log.Println("查询粉丝记录失败:", err)
		return err
	} else {
		if isFollowing == true {
			return errors.New("已有该粉丝")
		}
		// 已取关，重新关注
		query = "UPDATE follower SET is_following = true WHERE user_id = ? AND follower_id = ?"
		_, err = dB.Exec(query, userId, followerId)
		if err != nil {
			log.Println("增加粉丝记录失败:", err)
			return err
		}
	}
	return nil
}

func (repo *MySQLRelationRepository) RemoveFollower(userId int64, followerId int64) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// 更新is_following为0
	query := "UPDATE follower SET is_following = false WHERE user_id = ? AND follower_id = ?"
	_, err := dB.Exec(query, userId, followerId)
	if err != nil {
		log.Println("移除粉丝记录失败:", err)
		return err
	}
	return nil
}

func (repo *MySQLRelationRepository) AddFollow(userId int64, followId int64) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 检查是否已存在相同的记录
	query := "SELECT is_following FROM follow WHERE user_id = ? and follow_id = ?"
	row := dB.QueryRow(query, userId, followId)
	var isFollow bool
	err := row.Scan(&isFollow)
	if err == sql.ErrNoRows {
		// 不存在记录，插入新记录
		query = "INSERT INTO follow (user_id, follow_id, is_following) VALUES (?, ?, true)"
		_, err = dB.Exec(query, userId, followId)
		if err != nil {
			log.Println("插入关注记录失败:", err)
			return err
		}
	} else if err != nil {
		log.Println("查询关注记录失败:", err)
		return err
	} else {
		if isFollow == true {
			return errors.New("已关注该人")
		}
		// 已取关，重新关注
		query = "UPDATE follow SET is_following = true WHERE user_id = ? AND follow_id = ?"
		_, err = dB.Exec(query, userId, followId)
		if err != nil {
			log.Println("更新关注记录失败:", err)
			return err
		}
	}
	return nil
}

func (repo *MySQLRelationRepository) RemoveFollow(userId int64, followId int64) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// 更新is_fan为0
	query := "UPDATE follow SET is_following = false WHERE user_id = ? AND follow_id = ?"
	_, err := dB.Exec(query, userId, followId)
	if err != nil {
		log.Println("移除关注记录失败:", err)
		return err
	}
	return nil
}

func (repo *MySQLRelationRepository) CheckFollow(userId, followId int64) (bool, error) {
	query := "SELECT is_following FROM follow WHERE user_id = ? AND follow_id = ?"
	row := dB.QueryRow(query, userId, followId)
	var isFollow bool
	err := row.Scan(&isFollow)
	if err == sql.ErrNoRows {
		// 不存在记录，没有关注
		return false, nil
	} else if err != nil {
		log.Println("查询关注记录失败:", err)
		return false, err
	} else {
		if isFollow == true {
			return true, nil
		}
		// 已取消关注
		return false, nil
	}
}

func (repo *MySQLRelationRepository) GetFollowById(userId int64) ([]int64, error) {
	query := "SELECT follow_id FROM follow WHERE user_id = ? AND is_following = true"
	rows, err := dB.Query(query, userId)
	if err != nil {
		log.Println("查询关注列表失败:", err)
		return nil, err
	}
	var followIds []int64
	for rows.Next() {
		var followId int64
		err := rows.Scan(&followId)
		if err != nil {
			log.Println("扫描关注数据失败:", err)
			return nil, err
		}
		followIds = append(followIds, followId)
	}
	if err = rows.Err(); err != nil {
		log.Println("遍历关注结果失败:", err)
		return nil, err
	}
	return followIds, nil
}

func (repo *MySQLRelationRepository) GetFollowerById(userId int64) ([]int64, error) {
	query := "SELECT follower_id FROM follower WHERE user_id = ? AND is_following = true"
	rows, err := dB.Query(query, userId)
	if err != nil {
		log.Println("查询粉丝列表失败:", err)
		return nil, err
	}
	var followerIds []int64
	for rows.Next() {
		var followerId int64
		err := rows.Scan(&followerId)
		if err != nil {
			log.Println("扫描粉丝数据失败:", err)
			return nil, err
		}
		followerIds = append(followerIds, followerId)
	}
	if err = rows.Err(); err != nil {
		log.Println("遍历粉丝结果失败:", err)
		return nil, err
	}
	return followerIds, nil
}
