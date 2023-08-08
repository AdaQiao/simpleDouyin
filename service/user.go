package service

import (
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/controller"
	"log"
	"net"
	"net/rpc"
)

// 用户结构体

// 用户服务接口
type UserService interface {
	Register(user controller.UserPassword, reply *string) error
	Login(user controller.UserPassword, reply *string) error
}

// 用户服务实现
type UserServiceImpl struct {
	users []controller.UserPassword
}

// 用户注册
func (s *UserServiceImpl) Register(user controller.UserPassword, reply *string) error {
	// 检查用户名是否已存在
	for _, u := range s.users {
		if u.Username == user.Username {
			return errors.New("用户名已存在")
		}
	}

	// 注册用户
	s.users = append(s.users, user)
	*reply = "注册成功"
	return nil
}

// 用户登录
func (s *UserServiceImpl) Login(user controller.UserPassword, reply *string) error {
	// 查找用户
	for _, u := range s.users {
		if u.Username == user.Username && u.Password == user.Password {
			*reply = "登录成功"
			return nil
		}
	}

	return errors.New("用户名或密码错误")
}

func main() {
	// 创建用户服务实例
	userService := &UserServiceImpl{}

	// 注册RPC服务
	rpc.Register(userService)

	// 启动RPC服务器
	listener, err := net.Listen("tcp", "127.0.0.1:9090")
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
