package db

import (
	"log"
	"sync"
	"time"

	"github.com/AdaQiao/simpleDouyin/model"
)

type CommentRepository interface {
	AddComment(userID, videoID int64, comment_text string) (int64, error)
	RemoveComment(userID, videoID, comment_id int64) error
	GetCommentIdByVideoId(videoId int64) ([]int64, error)
	GetCommentByCommentId(commentId int64) (*model.Comment, int64, error)
}

type MySQLCommentRepository struct {
	mutex sync.Mutex
}

func NewMySQLCommentRepository() *MySQLCommentRepository {
	return &MySQLCommentRepository{}
}

func (repo *MySQLCommentRepository) AddComment(userID, videoID int64, comment_text string) (int64, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	//添加新评论
	var err error = nil
	sqlDateFormat := "2006-01-02"
	current_time := time.Now().Format(sqlDateFormat)
	query := "INSERT INTO comment (user_id, video_id, create_date, comment_text) VALUES (?, ?, ?, ?)"
	result, err := dB.Exec(query, userID, videoID, current_time, comment_text)
	if err != nil {
		log.Println("添加评论失败:", err)
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

func (repo *MySQLCommentRepository) RemoveComment(userID, videoID, comment_id int64) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 删除指定Comment行
	query := "DELETE FROM comment WHERE user_id = ? AND video_id = ? AND comment_id = ?"
	_, err := dB.Exec(query, userID, videoID, comment_id)
	if err != nil {
		log.Println("删除评论失败:", err)
		return err
	}
	return nil
}

func (repo *MySQLCommentRepository) GetCommentIdByVideoId(videoId int64) ([]int64, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 执行查询视频数据的SQL语句
	query := "SELECT comment_id FROM comment WHERE video_id = ? ORDER BY create_date DESC"
	rows, err := dB.Query(query, videoId)
	if err != nil {
		log.Println("查询评论列表失败:", err)
		return nil, err
	}

	var commentIds []int64

	for rows.Next() {
		var commentId int64
		err := rows.Scan(&commentId)
		if err != nil {
			log.Println("扫描评论数据失败:", err)
			return nil, err
		}
		commentIds = append(commentIds, commentId)
	}

	if err := rows.Err(); err != nil {
		log.Println("遍历评论结果失败:", err)
		return nil, err
	}
	return commentIds, nil
}

func (repo *MySQLCommentRepository) GetCommentByCommentId(commentId int64) (*model.Comment, int64, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	var userId int64
	// 执行查找评论数据的SQL语句
	query := "SELECT user_id, comment_id, comment_text, create_date FROM comment WHERE comment_id = ?"
	var comment model.Comment
	rows, err := dB.Query(query, commentId)
	if err != nil {
		log.Println("查询评论失败：", err)
		return nil, 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&userId,
			&comment.Id,
			&comment.Content,
			&comment.CreateDate,
		)
		if err != nil {
			log.Println("扫描评论失败:", err)
			return nil, 0, err
		}
	} else {
		log.Println("未找到匹配的评论")
		return nil, 0, nil
	}
	return &comment, userId, nil
}
