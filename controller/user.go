package controller

import (
	"github.com/AdaQiao/simpleDouyin/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/rpc"
)

var UsersLoginInfo = map[string]model.User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

// UsersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 连接到远程RPC服务器
	client, err := rpc.Dial("tcp", "127.0.0.1:9091")
	if err != nil {
		log.Fatal("RPC连接失败：", err)
	}

	// 调用远程注册方法
	var reply model.UserLoginResponse
	err = client.Call("UserServiceImpl.Register", model.UserPassword{Username: username, Password: password}, &reply)
	if err != nil {
		log.Fatal("调用远程注册方法失败：", err)
	}
	c.JSON(http.StatusOK, reply)

}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	// 连接到远程RPC服务器
	client, err := rpc.Dial("tcp", "127.0.0.1:9091")
	if err != nil {
		log.Fatal("RPC连接失败：", err)
	}
	// 调用远程登录方法
	var reply model.UserLoginResponse
	err = client.Call("UserServiceImpl.Login", model.UserPassword{Username: username, Password: password}, &reply)
	if err != nil {
		log.Fatal("调用远程注册方法失败：", err)
	}
	c.JSON(http.StatusOK, reply)
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	//
	client, err := rpc.Dial("tcp", "127.0.0.1:9091")
	if err != nil {
		log.Fatal("RPC连接失败：", err)
	}
	//
	var reply model.UserResponse
	err = client.Call("UserServiceImpl.UserInfo", token, &reply)
	if err != nil {
		log.Fatal("调用远程注册方法失败：", err)
	}
	c.JSON(http.StatusOK, reply)
}
