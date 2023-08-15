package controller

import (
	"fmt"
	"github.com/AdaQiao/simpleDouyin/db"
	"github.com/AdaQiao/simpleDouyin/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var VideoRepo *db.MySQLVideoRepository

type FeedResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	VideoRepo = db.NewMySQLVideoRepository()
	videos, _ := VideoRepo.GetVideoById(3)
	fmt.Println(len(videos))
	c.JSON(http.StatusOK, FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
}
