package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/rpc"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

var UserIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

type UserPassword struct {
	Username string
	Password string
}

func Register(c *gin.Context) {

	// 连接到远程RPC服务器
	client, err := rpc.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		log.Fatal("RPC连接失败：", err)
	}

	// 调用远程注册方法
	err = client.Call("UserServiceImpl.Register", UserPassword{Username: "john", Password: "password"}, c)
	if err != nil {
		log.Fatal("调用远程注册方法失败：", err)
	}

}

func Login(c *gin.Context) {
	//username := c.Query("username")
	//password := c.Query("password")
	//token := username + password
	//
	//if user, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   user.Id,
	//		Token:    token,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	//	})
	//}
}

func UserInfo(c *gin.Context) {
	//token := c.Query("token")
	//
	//if user, exist := UsersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: Response{StatusCode: 0},
	//		User:     user,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	//	})
	//}
}
