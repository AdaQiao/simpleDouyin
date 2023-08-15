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

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	lastestTime := c.Query("latest_time")
	VideoRepo = db.NewMySQLVideoRepository()
	curTime, err := strconv.ParseInt(lastestTime, 10, 64)
	if err != nil || curTime == int64(0) {
		curTime = time.Now().Unix()
	}
	fmt.Println("curTime:", curTime)
	fmt.Println("now:", time.Now().Unix())
	videos, nextTime, _ := VideoRepo.GetVideosByTimestamp(curTime)
	c.JSON(http.StatusOK, model.FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  nextTime,
	})
}
