package db

import (
	"fmt"
	_ "fmt"
	"github.com/RaymondCode/simple-demo/model"
	"log"
)

type VideoRepository interface {
	CreateVideo(video model.Video, token string) error
	GetVideoByToken(token string) ([]model.Video, error)
}

type MySQLVideoRepository struct {
}

func NewMySQLVideoRepository() *MySQLVideoRepository {
	return &MySQLVideoRepository{}
}

func (repo *MySQLVideoRepository) CreateVideo(video model.Video, token string) error {
	// 执行插入视频数据的SQL语句
	query := `
		INSERT INTO videos (token, author_id, play_url, cover_url, favorite_count, comment_count, is_favorite, title)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	// 执行插入操作
	_, err := dB.Exec(query, token, video.Author.Id, video.PlayUrl, video.CoverUrl, 0, 0, 0, video.Title)
	if err != nil {
		log.Println("插入视频数据失败:", err)
		return err
	}

	return nil
}

func (repo *MySQLVideoRepository) GetVideoByToken(token string, userId string) ([]model.Video, error) {
	// 执行查询视频数据的SQL语句
	query := `
		SELECT author_id, play_url, cover_url, favorite_count, comment_count, is_favorite, title
		FROM videos
		WHERE  author_id = ?
	`
	rows, err := dB.Query(query, userId)
	if err != nil {
		log.Println("查询视频失败:", err)
		return nil, err
	}
	defer rows.Close()
	fmt.Printf("%v\n", rows)
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
	fmt.Println(len(videos))
	return videos, nil
}