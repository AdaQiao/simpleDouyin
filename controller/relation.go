package controller

import (
	"github.com/AdaQiao/simpleDouyin/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/rpc"
	"strconv"
)

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	var action_type int32
	if c.Query("action_type") == "1" {
		action_type = 1
	} else {
		action_type = 2
	}
	mes := model.FollowActionMessage{
		ToUserId:   toUserId,
		Token:      token,
		ActionType: action_type,
	}
	client, err := rpc.Dial("tcp", "127.0.0.1:9096")
	if err != nil {
		log.Fatal("RPC连接失败：", err)
	}
	// 调用远程登录方法
	var reply model.Response
	err = client.Call("RelationServiceImpl.FollowAction", mes, &reply)
	if err != nil {
		log.Fatal("调用远程注册方法失败：", err)
	}
	c.JSON(http.StatusOK, reply)
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	c.JSON(http.StatusOK, model.UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: []model.User{model.DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, model.UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: []model.User{model.DemoUser},
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, model.UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: []model.User{model.DemoUser},
	})
}
