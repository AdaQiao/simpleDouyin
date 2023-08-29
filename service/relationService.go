package service

import (
	"fmt"
	"github.com/AdaQiao/simpleDouyin/db"
	"github.com/AdaQiao/simpleDouyin/model"
	"log"
	"net"
	"net/rpc"
)

type RelationService interface {
	FollowAction(req model.FollowActionMessage, reply *model.Response) error
	FollowList(userIDToken model.UserIdToken, reply *model.VideoListResponse) error
}

func (s *FavoriteServiceImpl) FollowAction(req model.FollowActionMessage, reply *model.Response) error {
	return nil
}

type RelationServiceImpl struct {
	UserRepo     *db.MySQLUserRepository
	VideoRepo    *db.MySQLVideoRepository
	FavoriteRepo *db.MySQLFavoriteRepository
}

func RunRelationServer() {
	// 创建服务实例
	relationService := &RelationServiceImpl{
		UserRepo:     db.NewMySQLUserRepository(),
		VideoRepo:    db.NewMySQLVideoRepository(),
		FavoriteRepo: db.NewMySQLFavoriteRepository(),
	}

	// 注册RPC服务
	rpc.Register(relationService)

	// 启动RPC服务器
	listener, err := net.Listen("tcp", "127.0.0.1:9096")
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
