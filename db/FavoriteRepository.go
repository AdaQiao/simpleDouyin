package db

import (
	"database/sql"
	"errors"
	"log"
	"sync"
)

type FavoriteRepository interface {
	AddFavorite(userID, videoID int64) error
	RemoveFavorite(userID, videoID int64) error
	CheckFavorite(userID, videoID int64) (bool, error)
	GetFavoriteVideoIdByUserId(userId int64) ([]int64, error)
}

type MySQLFavoriteRepository struct {
	mutex sync.Mutex
}

func NewMySQLFavoriteRepository() *MySQLFavoriteRepository {
	return &MySQLFavoriteRepository{}
}

func (repo *MySQLFavoriteRepository) AddFavorite(userID, videoID int64) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 检查是否已存在相同的记录
	query := "SELECT is_favorite FROM favorite WHERE user_id = ? AND video_id = ?"
	row := dB.QueryRow(query, userID, videoID)
	var isFavorite int
	err := row.Scan(&isFavorite)
	if err == sql.ErrNoRows {
		// 不存在记录，插入新记录
		query = "INSERT INTO favorite (user_id, video_id, is_favorite) VALUES (?, ?, 1)"
		_, err = dB.Exec(query, userID, videoID)
		if err != nil {
			log.Println("插入喜欢记录失败:", err)
			return err
		}
	} else if err != nil {
		log.Println("查询喜欢记录失败:", err)
		return err
	} else {
		if isFavorite == 1 {
			return errors.New("已点赞")
		}
		// 已取消点赞，更新 is_favorite 为 1
		query = "UPDATE favorite SET is_favorite = 1 WHERE user_id = ? AND video_id = ?"
		_, err = dB.Exec(query, userID, videoID)
		if err != nil {
			log.Println("更新喜欢记录失败:", err)
			return err
		}
	}

	return nil
}

func (repo *MySQLFavoriteRepository) RemoveFavorite(userID, videoID int64) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// 更新is_favorite为0
	query := "UPDATE favorite SET is_favorite = 0 WHERE user_id = ? AND video_id = ?"
	_, err := dB.Exec(query, userID, videoID)
	if err != nil {
		log.Println("更新喜欢记录失败:", err)
		return err
	}

	return nil
}

func (repo *MySQLFavoriteRepository) CheckFavorite(userID, videoID int64) (bool, error) {
	query := "SELECT is_favorite FROM favorite WHERE user_id = ? AND video_id = ?"
	row := dB.QueryRow(query, userID, videoID)
	var isFavorite int
	err := row.Scan(&isFavorite)
	if err == sql.ErrNoRows {
		// 不存在记录，没有点过赞
		return false, nil
	} else if err != nil {
		log.Println("查询喜欢记录失败:", err)
		return false, err
	} else {
		if isFavorite == 1 {
			return true, nil
		}
		// 已取消点赞
		return false, nil
	}
}

func (repo *MySQLFavoriteRepository) GetFavoriteVideoIdByUserId(userId int64) ([]int64, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// 执行查询视频数据的SQL语句
	query := `
		SELECT video_id FROM favorite WHERE user_id = ? AND is_favorite = 1
	`
	rows, err := dB.Query(query, userId)
	if err != nil {
		log.Println("查询喜欢列表失败:", err)
		return nil, err
	}
	var videoIds []int64
	for rows.Next() {
		var videoId int64
		err := rows.Scan(&videoId)
		if err != nil {
			log.Println("扫描视频数据失败:", err)
			return nil, err
		}
		videoIds = append(videoIds, videoId)
	}

	if err := rows.Err(); err != nil {
		log.Println("遍历视频结果失败:", err)
		return nil, err
	}
	return videoIds, nil
}
