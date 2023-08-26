package db

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type MessageRepository interface {
	AddMessage(userID, videoID int64) error
	GetMessageByFromId(userId int64) ([]int64, error)
	GetMessageByToId(userId int64) ([]int64, error)
}

type MySQLMessageRepository struct {
	mutex sync.Mutex
}

func NewMySQLMessageRepository() *MySQLMessageRepository {
	return &MySQLMessageRepository{}
}

func (repo *MySQLMessageRepository) AddMessage(fromID, toID int64, content string) (int64, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	//添加新消息
	var err error = nil
	sqlDateFormat := "2006-01-02"
	current_time := time.Now().Format(sqlDateFormat)

	// 生成一个随机的 64 位整数作为id
	rand.Seed(time.Now().UnixNano())
	var id = rand.Int63()
	//插入数据库
	query := "INSERT INTO message (id, message, create_date,fromID,toID) VALUES (?, ?, ?, ?, ?)"
	result, err := dB.Exec(query, id, content, current_time, fromID, toID)
	if err != nil {
		log.Println("添加消息失败:", err)
		return 0, err
	}
	// 获取插入数据的自增ID
	commentID, err := result.LastInsertId()
	if err != nil {
		log.Println("获取插入数据的ID失败:", err)
		return 0, err
	}
	return commentID, nil
}

func (repo *MySQLMessageRepository) GetMessageByFromId(userId int64) ([]int64, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 执行查询消息的SQL语句
	query := `
		SELECT id FROM message WHERE from_user_id = ? OR to_user_id = ?
		ORDER BY create_date DESC
	`
	rows, err := dB.Query(query, userId)
	if err != nil {
		log.Println("获取消息列表失败:", err)
		return nil, err
	}
	var messageIds []int64
	for rows.Next() {
		var messageId int64
		err := rows.Scan(&messageId)
		if err != nil {
			log.Println("扫描消息数据失败:", err)
			return nil, err
		}
		messageIds = append(messageIds, messageId)
	}

	if err := rows.Err(); err != nil {
		log.Println("遍历消息列表失败:", err)
		return nil, err
	}
	return messageIds, nil
}
