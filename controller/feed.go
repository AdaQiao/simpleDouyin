package controller

import (
	"github.com/AdaQiao/simpleDouyin/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/rpc"
)

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	lastestTime := c.Query("latest_time")
	// 连接到远程RPC服务器
	client, err := rpc.Dial("tcp", "127.0.0.1:9093")
	if err != nil {
		log.Fatal("RPC连接失败：", err)
	}
	// 调用远程登录方法
	var reply model.FeedResponse
	err = client.Call("FeedServiceImpl.GetVideoList", lastestTime, &reply)
	if err != nil {
		log.Fatal("调用远程注册方法失败：", err)
	}
	c.JSON(http.StatusOK, reply)
}
