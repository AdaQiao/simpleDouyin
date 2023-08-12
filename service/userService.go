package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"strings"

	"log"
	"net"

	"net/rpc"

	"sync/atomic"
)

// 用户服务接口
type UserService interface {
	Register(user model.UserPassword, reply *model.UserLoginResponse) error
	Login(user model.UserPassword, reply *model.UserLoginResponse) error
	UserInfo(token string, reply *model.UserResponse)
}

var repo *db.MySQLUserRepository = db.NewMySQLUserRepository()

// 用户服务实现
type UserServiceImpl struct {
}

// 用户注册
func (s *UserServiceImpl) Register(user model.UserPassword, reply *model.UserLoginResponse) error {
	//检查用户名是否已存在
	token := user.Username + user.Password
	// 调用存储库的 CreateUser 函数执行插入操作
	if err := repo.CreateUser(user); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			// 处理唯一约束错误
			errorMsg := fmt.Sprintf("用户名 %s 已存在\n", user.Username)
			*reply = model.UserLoginResponse{
				Response: model.Response{StatusCode: 1, StatusMsg: errorMsg},
			}
			return nil
		}
		log.Println("插入用户数据失败:", err)
		return err
	}

	atomic.AddInt64(&controller.UserIdSequence, 1)
	newUser := model.User{
		Id:   controller.UserIdSequence,
		Name: user.Username,
	}
	controller.UsersLoginInfo[token] = newUser
	*reply = model.UserLoginResponse{
		Response: model.Response{StatusCode: 0},
		UserId:   controller.UserIdSequence,
		Token:    token,
	}
	return nil
}

// 用户登录
func (s *UserServiceImpl) Login(user model.UserPassword, reply *model.UserLoginResponse) error {
	token := user.Username + user.Password
	userInfo, err := repo.GetUser(token)
	if err == nil {
		*reply = model.UserLoginResponse{
			Response: model.Response{StatusCode: 0},
			UserId:   userInfo.Id,
			Token:    token,
		}
	} else {
		*reply = model.UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		}
	}
	return nil
}

// 用户信息
func (s *UserServiceImpl) UserInfo(token string, reply *model.UserResponse) error {
	userInfo, err := repo.GetUser(token)
	if err == nil {
		*reply = model.UserResponse{
			Response: model.Response{StatusCode: 0},
			User:     *userInfo,
		}
	} else {
		*reply = model.UserResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
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