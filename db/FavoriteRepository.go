package db

import (
	"database/sql"
	"log"
	"sync"
)

type FavoriteRepository interface {
	AddFavorite(userID, videoID int64) error
	RemoveFavorite(userID, videoID int64) error
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
	query := "SELECT id FROM favorite WHERE user_id = ? AND video_id = ?"
	row := dB.QueryRow(query, userID, videoID)
	var id int
	err := row.Scan(&id)
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
		// 已存在记录，更新is_favorite为1
		query = "UPDATE user_likes SET is_favorite = 1 WHERE id = ?"
		_, err = dB.Exec(query, id)
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
