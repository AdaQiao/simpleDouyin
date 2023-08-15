package controller

import (
	"fmt"
	"github.com/AdaQiao/simpleDouyin/db"
	"github.com/AdaQiao/simpleDouyin/model"
	"github.com/AdaQiao/simpleDouyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var VideoRepo *db.MySQLVideoRepository

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	lastestTime := c.Query("latest_time")
	curTime, err := strconv.ParseInt(lastestTime, 10, 64)
	if err != nil || curTime == int64(0) {
		curTime = time.Now().Unix()
	}
	fmt.Println("curTime:", curTime)
	videos, nextTime, _ := service.GetVideoList(curTime)
	c.JSON(http.StatusOK, model.FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  nextTime,
	})
}
