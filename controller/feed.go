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

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	VideoRepo = db.NewMySQLVideoRepository()
	videos, _ := VideoRepo.GetVideosByTimestamp(time.Now())
	fmt.Println(len(videos))
	c.JSON(http.StatusOK, model.FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
}
