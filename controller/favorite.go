package controller

import (
	"github.com/AdaQiao/simpleDouyin/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/rpc"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	var action_type int32
	if c.Query("action_type") == "1" {
		action_type = 1
	} else {
		action_type = 2
	}
	token := c.Query("token")
	mes := model.FavoriteMessage{
		VideoId:    videoId,
		Token:      token,
		ActionType: action_type,
	}
	client, err := rpc.Dial("tcp", "127.0.0.1:9094")
	if err != nil {
		log.Fatal("RPC连接失败：", err)
	}
	// 调用远程登录方法
	var reply model.Response
	err = client.Call("FavoriteServiceImpl.FavoriteVideo", mes, &reply)
	if err != nil {
		log.Fatal("调用远程注册方法失败：", err)
	}
	c.JSON(http.StatusOK, reply)
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	userIDStr := c.Query("user_id")
	// 连接到远程RPC服务器
	client, err := rpc.Dial("tcp", "127.0.0.1:9094")
	if err != nil {
		log.Println("RPC连接失败：", err)
	}

	user_id, err := strconv.ParseInt(userIDStr, 10, 64)
	// 调用远程注册方法
	var reply model.VideoListResponse
	err = client.Call("FavoriteServiceImpl.FavoriteList", model.UserIdToken{
		Token:  token,
		UserId: user_id,
	}, &reply)
	if err != nil {
		log.Println("调用远程注册方法失败：", err)
	}
	c.JSON(http.StatusOK, reply)

}
