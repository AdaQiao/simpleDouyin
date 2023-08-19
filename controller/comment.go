package controller

import (
	"log"
	"net/http"
	"net/rpc"
	"strconv"

	"github.com/AdaQiao/simpleDouyin/model"
	"github.com/gin-gonic/gin"
)

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	var actionType int32
	if c.Query("action_type") == "1" {
		actionType = 1
	} else {
		actionType = 2
	}
	commentTxet := c.Query("comment_text")
	commentId, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)

	mes := model.CommentActionRequest{
		Token:       token,
		VideoId:     videoId,
		ActionType:  actionType,
		CommentText: commentTxet,
		CommentId:   commentId,
	}

	client, err := rpc.Dial("tcp", "127.0.0.1:9094")
	if err != nil {
		log.Fatal("RPC连接失败：", err)
	}
	// 调用远程登录方法
	var reply model.CommentActionResponse
	err = client.Call("CommentServiceImpl.Comment", mes, &reply)
	if err != nil {
		log.Fatal("调用远程注册方法失败：", err)
	}
	c.JSON(http.StatusOK, reply)
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	token := c.Query("token")
	userIDStr := c.Query("user_id")
	user_id, _ := strconv.ParseInt(userIDStr, 10, 64)

	// 连接到远程RPC服务器
	client, err := rpc.Dial("tcp", "127.0.0.1:9094")
	if err != nil {
		log.Println("RPC连接失败：", err)
	}

	// 调用远程注册方法
	var reply model.CommentListResponse
	err = client.Call("CommentServiceImpl.CommentList", model.UserIdToken{
		Token:  token,
		UserId: user_id,
	}, &reply)
	if err != nil {
		log.Println("调用远程注册方法失败：", err)
	}
	c.JSON(http.StatusOK, reply)
}
