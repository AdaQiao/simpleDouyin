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
}

func (s *RelationServiceImpl) FollowAction(req model.FollowActionMessage, reply *model.Response) error {
	userId, err := s.UserRepo.GetUserId(req.Token)
	if err != nil {
		*reply = model.Response{
			StatusCode: 1,
			StatusMsg:  "user didn't exist",
		}
		return fmt.Errorf("user didn't uploaded")
	}
	toUserId := req.ToUserId
	if req.ActionType == 1 {
		err = s.RelationRepo.AddFollow(userId, toUserId)
		if err != nil {
			*reply = model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			}
			return nil
		}
		err = s.RelationRepo.AddFan(toUserId, userId)
		if err != nil {
			*reply = model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			}
			return nil
		}
		err = s.UserRepo.UpdateFollowCount(userId, 1)
		if err != nil {
			*reply = model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			}
			return nil
		}
		err = s.UserRepo.UpdateFollowerCount(toUserId, 1)
		if err != nil {
			*reply = model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			}
			return nil
		}
	} else if req.ActionType == 2 {
		err = s.RelationRepo.RemoveFollow(userId, toUserId)
		if err != nil {
			*reply = model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			}
			return nil
		}
		err = s.RelationRepo.RemoveFan(toUserId, userId)
		if err != nil {
			*reply = model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			}
			return nil
		}
		err = s.UserRepo.UpdateFollowCount(userId, 2)
		if err != nil {
			*reply = model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			}
			return nil
		}
		err = s.UserRepo.UpdateFollowerCount(toUserId, 2)
		if err != nil {
			*reply = model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			}
			return nil
		}
	}
	*reply = model.Response{
		StatusCode: 0,
		StatusMsg:  "Relation Action success",
	}
	return nil
}

type RelationServiceImpl struct {
	UserRepo     *db.MySQLUserRepository
	VideoRepo    *db.MySQLVideoRepository
	FavoriteRepo *db.MySQLFavoriteRepository
	RelationRepo *db.MySQLRelationRepository
}

func RunRelationServer() {
	// 创建服务实例
	relationService := &RelationServiceImpl{
		UserRepo:     db.NewMySQLUserRepository(),
		VideoRepo:    db.NewMySQLVideoRepository(),
		FavoriteRepo: db.NewMySQLFavoriteRepository(),
		RelationRepo: db.NewMySQLRelationRepository(),
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
