package service

import (
	"fmt"
	"github.com/cpl/simple-demo/controller"
	"log"
	"net"

	"net/rpc"

	"github.com/cpl/simple-demo/db"
)

// 用户服务接口
type UserService interface {
	Register(user controller.UserPassword, reply *controller.UserLoginResponse) error
	Login(user controller.UserPassword, reply *controller.UserLoginResponse) error
	UserInfo(token string, reply *controller.UserResponse)
}

// 用户服务实现
type UserServiceImpl struct {
	repo db.UserRepository
}

// 用户注册
func (s *UserServiceImpl) Register(user controller.UserPassword, reply *controller.UserLoginResponse) error {
	//检查用户名是否已存在
	token := user.Username + user.Password
	/*if _, exist := controller.UsersLoginInfo[token]; exist {
	*reply = controller.UserLoginResponse{
		Response: controller.Response{StatusCode: 1, StatusMsg: "User already exist"},
	}*/
	_, err := s.repo.GetUser(token)
	if err == nil {
		*reply = controller.UserLoginResponse{
			Response: controller.Response{StatusCode: 1, StatusMsg: "User already exists"},
		}
	} else {

		// 调用存储库的 CreateUser 函数执行插入操作
		if err := s.repo.CreateUser(user); err != nil {
			log.Println("插入用户数据失败:", err)
			return err
		}

		/*		atomic.AddInt64(&controller.UserIdSequence, 1)
				newUser := controller.User{
					Id:   controller.UserIdSequence,
					Name: user.Username,
				}
				controller.UsersLoginInfo[token] = newUser
				*reply = controller.UserLoginResponse{
					Response: controller.Response{StatusCode: 0},
					UserId:   controller.UserIdSequence,
					Token:    token,
				}*/

	}
	return nil
}

// 用户登录
func (s *UserServiceImpl) Login(user controller.UserPassword, reply *controller.UserLoginResponse) error {
	/*token := user.Username + user.Password
	if userInfo, exist := controller.UsersLoginInfo[token]; exist {
		*reply = controller.UserLoginResponse{
			Response: controller.Response{StatusCode: 0},
			UserId:   userInfo.Id,
			Token:    token,
		}
	} else {
		*reply = controller.UserLoginResponse{
			Response: controller.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		}
	}
	return nil*/

	token := user.Username + user.Password
	userInfo, err := s.repo.GetUser(token)
	if err == nil {
		*reply = controller.UserLoginResponse{
			Response: controller.Response{StatusCode: 0},
			UserId:   userInfo.Id,
			Token:    token,
		}
	} else {
		*reply = controller.UserLoginResponse{
			Response: controller.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		}
	}
	return nil
}

// 用户信息
func (s *UserServiceImpl) UserInfo(token string, reply *controller.UserResponse) error {
	/*if userInfo, exist := controller.UsersLoginInfo[token]; exist {
		*reply = controller.UserResponse{
			Response: controller.Response{StatusCode: 0},
			User:     userInfo,
		}
	} else {
		*reply = controller.UserResponse{
			Response: controller.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		}
	}
	return nil*/

	userInfo, err := s.repo.GetUser(token)
	if err == nil {
		*reply = controller.UserResponse{
			Response: controller.Response{StatusCode: 0},
			User:     *userInfo,
		}
	} else {
		*reply = controller.UserResponse{
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
