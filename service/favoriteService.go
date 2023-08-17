package service

import (
	"fmt"
	"github.com/AdaQiao/simpleDouyin/db"
	"github.com/AdaQiao/simpleDouyin/model"
	"log"
	"net"
	"net/rpc"
)

type FavoriteService interface {
	FavoriteVideo(req model.FavoriteMessage, reply string) error
}
type FavoriteServiceImpl struct {
	UserRepo  *db.MySQLUserRepository
	VideoRepo *db.MySQLVideoRepository
}

func (s *FavoriteServiceImpl) FavoriteVideo(req model.FavoriteMessage, reply string) error {

}
func RunFaServer() {
	// 创建服务实例

	favoriteService := &FavoriteServiceImpl{
		UserRepo:  db.NewMySQLUserRepository(),
		VideoRepo: db.NewMySQLVideoRepository(),
	}

	// 注册RPC服务
	rpc.Register(favoriteService)

	// 启动RPC服务器
	listener, err := net.Listen("tcp", "127.0.0.1:9094")
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
