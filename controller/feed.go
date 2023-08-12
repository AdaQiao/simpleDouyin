package controller

import (
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
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
	viseos, _ := VideoRepo.GetVideoByToken("qjwe123456", 3)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: viseos,
		NextTime:  time.Now().Unix(),
	})
}
