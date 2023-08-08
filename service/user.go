package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/controller"

	"github.com/gin-gonic/gin"
	"log"
	"net"

	"net/rpc"

	"sync/atomic"
)

// 用户服务接口
type UserService interface {
	Register(user controller.UserPassword, reply *controller.UserLoginResponse) error
	Login(c *gin.Context, reply *string) error
}

// 用户服务实现
type UserServiceImpl struct {
}

// 用户注册
func (s *UserServiceImpl) Register(user controller.UserPassword, reply *controller.UserLoginResponse) error {
	//检查用户名是否已存在
	token := user.Username + user.Password
	if _, exist := controller.UsersLoginInfo[token]; exist {
		reply = &controller.UserLoginResponse{
			Response: controller.Response{StatusCode: 1, StatusMsg: "User already exist"},
		}
	} else {
		atomic.AddInt64(&controller.UserIdSequence, 1)
		newUser := controller.User{
			Id:   controller.UserIdSequence,
			Name: user.Username,
		}
		controller.UsersLoginInfo[token] = newUser
		reply = &controller.UserLoginResponse{
			Response: controller.Response{StatusCode: 0},
			UserId:   controller.UserIdSequence,
			Token:    token,
		}
	}
	return nil
}

// 用户登录
func (s *UserServiceImpl) Login(user controller.UserPassword, reply *controller.UserLoginResponse) error {

	token := user.Username + user.Password
	//
	fmt.Println(token)
	if userInfo, exist := controller.UsersLoginInfo[token]; exist {
		reply = &controller.UserLoginResponse{
			Response: controller.Response{StatusCode: 0},
			UserId:   userInfo.Id,
			Token:    token,
		}
	} else {
		reply = &controller.UserLoginResponse{
			Response: controller.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		}
	}
	return nil
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
