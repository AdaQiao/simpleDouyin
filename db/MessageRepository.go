package db

import (
	"github.com/AdaQiao/simpleDouyin/model"
	"log"
	"math/rand"
	"sync"
	"time"
)

type MessageRepository interface {
	AddMessage(fromID, toID int64, content string) (int64, error)
	GetMessageByUserId(userId int64) ([]int64, error)
	GetMessageById(Id int64) (*model.Message, error)
	GetMessagesBetweenUsers(userId1, userId2 int64) ([]model.Message, error)
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
	currentTime := time.Now().Format(sqlDateFormat)
	// 应该用message里的id
	// 生成一个随机的 64 位整数作为id
	rand.Seed(time.Now().UnixNano())
	var id = rand.Int63()
	//插入数据库
	query := "INSERT INTO message (id, msg_content, create_date,from_user_id ,to_user_id) VALUES (?, ?, ?, ?, ?)"
	result, err := dB.Exec(query, id, content, currentTime, fromID, toID)
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

func (repo *MySQLMessageRepository) GetMessageByUserId(userId int64) ([]int64, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// 执行查询消息的SQL语句
	query := `
		SELECT id FROM message WHERE from_user_id = ? OR to_user_id = ?
		ORDER BY create_date DESC
	`
	rows, err := dB.Query(query, userId, userId)
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

	err = rows.Close()
	if err != nil {
		return nil, err
	}
	return messageIds, nil
}

func (repo *MySQLMessageRepository) GetMessageById(Id int64) (*model.Message, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// 执行查找评论数据的SQL语句
	query := "SELECT id, msg_content, create_time, user_id,to_user_id FROM Message WHERE id = ?"
	var message model.Message
	rows, err := dB.Query(query, Id)
	if err != nil {
		log.Println("查询评论失败：", err)
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(
			&message.Id,
			&message.Content,
			&message.CreateTime,
			&message.FromUserId,
			&message.ToUserId,
		)
		if err != nil {
			log.Println("扫描评论失败:", err)
			return nil, err
		}
	} else {
		log.Println("未找到匹配的评论")
		return nil, nil
	}
	return &message, err
}

func (repo *MySQLMessageRepository) GetMessagesBetweenUsers(userId1, userId2 int64) ([]model.Message, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// 执行查找评论数据的SQL语句
	query := "SELECT id, content, create_time, from_user_id, to_user_id FROM Message WHERE (from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)  ORDER BY create_time ASC"
	var messages []model.Message
	rows, err := dB.Query(query, userId1, userId2, userId2, userId1)
	if err != nil {
		log.Println("查询评论失败：", err)
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var message model.Message
		err = rows.Scan(
			&message.Id,
			&message.Content,
			&message.CreateTime,
			&message.FromUserId,
			&message.ToUserId,
		)
		if err != nil {
			log.Println("扫描评论失败:", err)
			return nil, err
		}
		messages = append(messages, message)
	} else {
		log.Println("未找到匹配的评论")
		return nil, nil
	}
	return messages, err
}
