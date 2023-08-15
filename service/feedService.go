package service

import (
	"github.com/AdaQiao/simpleDouyin/db"
	"github.com/AdaQiao/simpleDouyin/model"
)

func GetVideoList(curTime int64) ([]model.Video, int64, error) {
	VideoRepo := db.NewMySQLVideoRepository()
	UserRepo := db.NewMySQLUserRepository()
	videos, nextTime, tokens, err := VideoRepo.GetVideosByTimestamp(curTime)
	if err != nil {
		return nil, 0, err
	}
	for i := 0; i < len(videos); i++ {
		user, _ := UserRepo.GetUser(tokens[i])
		videos[i].Author = *user
	}
	return videos, nextTime, nil
}
