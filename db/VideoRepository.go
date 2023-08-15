package db

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/AdaQiao/simpleDouyin/model"
)

type VideoRepository interface {
	CreateVideo(video model.Video, token string) error
	GetVideoById(userId int64) ([]model.Video, error)
}

type MySQLVideoRepository struct {
	mutex sync.Mutex
}

func NewMySQLVideoRepository() *MySQLVideoRepository {
	return &MySQLVideoRepository{}
}

func (repo *MySQLVideoRepository) CreateVideo(video model.Video, token string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 执行插入视频数据的SQL语句
	query := `
		INSERT INTO videos (token, author_id, play_url, cover_url, favorite_count, comment_count, is_favorite, title, created_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	// 执行插入操作
	_, err := dB.Exec(query, token, video.Author.Id, video.PlayUrl, video.CoverUrl, 0, 0, 0, video.Title, time.Now().Unix())
	if err != nil {
		log.Println("插入视频数据失败:", err)
		return err
	}

	return nil
}

func (repo *MySQLVideoRepository) GetVideoById(userId int64) ([]model.Video, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 执行查询视频数据的SQL语句
	query := `
		SELECT author_id, play_url, cover_url, favorite_count, comment_count, is_favorite, title FROM videos WHERE  author_id = ?
	`
	rows, err := dB.Query(query, userId)
	if err != nil {
		log.Println("查询视频失败:", err)
		return nil, err
	}
	var videos []model.Video
	for rows.Next() {
		var video model.Video
		err := rows.Scan(
			&video.Author.Id,
			&video.PlayUrl,
			&video.CoverUrl,
			&video.FavoriteCount,
			&video.CommentCount,
			&video.IsFavorite,
			&video.Title,
		)
		if err != nil {
			log.Println("扫描视频数据失败:", err)
			return nil, err
		}
		videos = append(videos, video)
	}

	if err := rows.Err(); err != nil {
		log.Println("遍历视频结果失败:", err)
		return nil, err
	}
	return videos, nil
}

func (repo *MySQLVideoRepository) GetVideosByTimestamp(timestamp time.Time) ([]model.Video, int64, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 执行查询视频数据的SQL语句
	query := `
		SELECT author_id, play_url, cover_url, favorite_count, comment_count, is_favorite, title, created_time
		FROM videos
		WHERE created_time <= ? 
		ORDER BY created_time DESC
		LIMIT 30
	`
	rows, err := dB.Query(query, timestamp)
	if err != nil {
		log.Println("查询视频失败:", err)
		return nil, 0, err
	}
	defer rows.Close()

	var videos []model.Video
	var tempTime int64
	var firstTime int64
	for rows.Next() {
		var video model.Video
		err := rows.Scan(
			&video.Author.Id,
			&video.PlayUrl,
			&video.CoverUrl,
			&video.FavoriteCount,
			&video.CommentCount,
			&video.IsFavorite,
			&video.Title,
			&tempTime,
		)
		if err != nil {
			log.Println("扫描视频数据失败:", err)
			return nil, 0, err
		}
		videos = append(videos, video)
		fmt.Println(video.PlayUrl + "sfadafafafasdfasf")
		// 保存第一个视频的created_time
	}
	firstTime = tempTime
	if err := rows.Err(); err != nil {
		log.Println("遍历视频结果失败:", err)
		return nil, 0, err
	}
	return videos, firstTime, nil
}
