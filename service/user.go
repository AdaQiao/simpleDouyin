package service

import (
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync/atomic"
)

// 用户服务接口
type UserService interface {
	Register(c *gin.Context, reply *string) error
	Login(c *gin.Context, reply *string) error
}

// 用户服务实现
type UserServiceImpl struct {
}

// 用户注册
func (s *UserServiceImpl) Register(c *gin.Context, reply *string) error {
	// 检查用户名是否已存在

	username := c.Query("username")
	password := c.Query("password")

	token := username + password
	if _, exist := controller.UsersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, controller.UserLoginResponse{
			Response: controller.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
    fmt.Println("afdasfasfadfa")
		atomic.AddInt64(&controller.UserIdSequence, 1)
		newUser := controller.User{
			Id:   controller.UserIdSequence,
			Name: username,
		}
		controller.UsersLoginInfo[token] = newUser
		c.JSON(http.StatusOK, controller.UserLoginResponse{
			Response: controller.Response{StatusCode: 0},
			UserId:   controller.UserIdSequence,
			Token:    username + password,
		})
	}

	return nil
}

// 用户登录
func (s *UserServiceImpl) Login(c *gin.Context, reply *string) error {
	// 查找用户
	//for _, u := range s.users {
	//	if u.Username == user.Username && u.Password == user.Password {
	//		*reply = "登录成功"
	//		return nil
	//	}
	//}
	//
	return errors.New("用户名或密码错误")
}

func RunUserServer() {
	// 创建用户服务实例
	userService := &UserServiceImpl{}

	// 注册RPC服务
	rpc.Register(userService)

	// 启动RPC服务器
	listener, err := net.Listen("tcp", "127.0.0.1:9091")
	if err != nil {
		log.Fatal("RPC服务器启动失败:", err)
	}

	fmt.Println("RPC服务器已启动，等待远程调用...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("接受连接失败:", err)
		}
		go rpc.ServeConn(conn)
	}
}
