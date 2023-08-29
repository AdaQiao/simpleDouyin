package db

import (
	"database/sql"
	"errors"
	"log"
	"sync"
)

type RelationRepository interface {
	AddFan(userId int64, fanId int64) error
	RemoveFan(userId int64, fanId int64) error
	AddFollow(userId int64, followId int64) error
	RemoveFollow(userId int64, followId int64) error
}

type MySQLRelationRepository struct {
	mutex sync.Mutex
}

func NewMySQLRelationRepository() *MySQLRelationRepository {
	return &MySQLRelationRepository{}
}

func (repo *MySQLRelationRepository) AddFan(userId int64, fanId int64) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 检查是否已存在相同的记录
	query := "SELECT is_fan FROM fan WHERE user_id = ?"
	row := dB.QueryRow(query, userId)
	var isFan bool
	err := row.Scan(&isFan)
	if err == sql.ErrNoRows {
		// 不存在记录，插入新记录
		query = "INSERT INTO fan (user_id, fan_id, is_fan) VALUES (?, ?, true)"
		_, err = dB.Exec(query, userId, fanId)
		if err != nil {
			log.Println("插入粉丝记录失败:", err)
			return err
		}
	} else if err != nil {
		log.Println("查询粉丝记录失败:", err)
		return err
	} else {
		if isFan == true {
			return errors.New("已有该粉丝")
		}
		// 已取关，重新关注
		query = "UPDATE fan SET is_fan = true WHERE user_id = ? AND fan_id = ?"
		_, err = dB.Exec(query, userId, fanId)
		if err != nil {
			log.Println("增加粉丝记录失败:", err)
			return err
		}
	}
	return nil
}

func (repo *MySQLRelationRepository) RemoveFan(userId int64, fanId int64) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// 更新is_fan为0
	query := "UPDATE fan SET is_fan = false WHERE user_id = ? AND fan_id = ?"
	_, err := dB.Exec(query, userId, fanId)
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
	query := "SELECT is_following FROM follow WHERE user_id = ?"
	row := dB.QueryRow(query, userId)
	var isFollow bool
	err := row.Scan(&isFollow)
	if err == sql.ErrNoRows {
		// 不存在记录，插入新记录
		query = "INSERT INTO fan (user_id, follow_id, is_following) VALUES (?, ?, true)"
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
