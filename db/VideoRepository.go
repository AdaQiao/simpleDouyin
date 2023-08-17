package db

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/AdaQiao/simpleDouyin/model"
)

type VideoRepository interface {
	CreateVideo(video model.Video, token string) error
	GetVideoById(userId int64) ([]model.Video, error)
	GetVideoByVideoId(videoId int64) (model.Video, error)
	GetVideosByTimestamp(timestamp int64) ([]model.Video, int64, []string, error)
	UpdateFavoriteCount(videoId int64, mode int32) error
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

func (repo *MySQLVideoRepository) GetVideoByVideoId(videoId int64) (model.Video, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 执行查询视频数据的SQL语句
	query := `
		SELECT id, author_id, play_url, cover_url, favorite_count, comment_count, is_favorite, title FROM videos WHERE id = ?
	`
	var video model.Video
	rows, err := dB.Query(query, videoId)
	if err != nil {
		log.Println("查询视频失败:", err)
		return video, err
	}

	err = rows.Scan(
		&video.Id,
		&video.Author.Id,
		&video.PlayUrl,
		&video.CoverUrl,
		&video.FavoriteCount,
		&video.CommentCount,
		&video.IsFavorite,
		&video.Title,
	)
	return video, nil
}
func (repo *MySQLVideoRepository) GetVideosByTimestamp(timestamp int64) ([]model.Video, int64, []string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 执行查询视频数据的SQL语句
	query := `
		SELECT id, token, play_url, cover_url, favorite_count, comment_count, is_favorite, title, created_time
		FROM videos
		WHERE created_time < ? 
		ORDER BY created_time DESC
		LIMIT 1
	`
	rows, err := dB.Query(query, timestamp)
	if err != nil {
		log.Println("查询视频失败:", err)
		return nil, 0, nil, err
	}
	defer rows.Close()

	var videos []model.Video
	var tokens []string
	var tempTime int64
	var firstTime int64
	var token string
	for rows.Next() {
		var video model.Video
		err := rows.Scan(
			&video.Id,
			&token,
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
			return nil, 0, nil, err
		}
		videos = append(videos, video)
		tokens = append(tokens, token)
		if len(videos) == 1 {
			firstTime = tempTime
		}
		// 保存第一个视频的created_time
	}

	if err := rows.Err(); err != nil {
		log.Println("遍历视频结果失败:", err)
		return nil, 0, nil, err
	}
	return videos, firstTime, tokens, nil
}

func (repo *MySQLVideoRepository) UpdateFavoriteCount(videoId int64, mode int32) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 根据mode选择执行的操作
	var query string
	if mode == 2 {
		query = `
			UPDATE videos SET favorite_count = favorite_count - 1 WHERE id = ?
		`
	} else if mode == 1 {
		query = `
			UPDATE videos SET favorite_count = favorite_count + 1 WHERE id = ?
		`
	} else {
		return errors.New("无效的mode值")
	}

	// 执行更新操作
	_, err := dB.Exec(query, videoId)
	if err != nil {
		log.Println("更新收藏数失败:", err)
		return err
	}

	return nil
}
