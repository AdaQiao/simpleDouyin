package db

import "github.com/RaymondCode/simple-demo/model"

type VideoRepository interface {
	CreateVideo(user model.Video) error
	GetVideoByToken(token string) (*model.Video, error)
}
