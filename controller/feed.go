package controller

import (
	"fmt"
	"github.com/AdaQiao/simpleDouyin/db"
	"github.com/AdaQiao/simpleDouyin/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var VideoRepo *db.MySQLVideoRepository

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

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	lastestTime := c.Query("latest_time")
	curTime, err := strconv.ParseInt(lastestTime, 10, 64)
	if err != nil || curTime == int64(0) {
		curTime = time.Now().Unix()
	}
	fmt.Println("curTime:", curTime)
	videos, nextTime, _ := GetVideoList(curTime)
	c.JSON(http.StatusOK, model.FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  nextTime,
	})
}
