package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"log"
	"net"
	"net/rpc"
)

type PublishService interface {
	Publish(req model.UploadViewReq, reply *model.Response) error
	PublishList(reply *model.VideoListResponse)
}

type PublishServiceImpl struct {
	repo *db.MySQLUserRepository
}

func (s *PublishServiceImpl) Publish(req model.UploadViewReq, reply *model.Response) error {
	user := s.repo.GetUser(req.Token)

	return nil
}
func (s *PublishServiceImpl) PublishList(reply *model.VideoListResponse) {

}

func RunPublishServer() {
	// 创建服务实例

	publishService := &PublishServiceImpl{
		repo: db.NewMySQLUserRepository(),
	}

	// 注册RPC服务
	rpc.Register(publishService)

	// 启动RPC服务器
	listener, err := net.Listen("tcp", "127.0.0.1:9092")
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
