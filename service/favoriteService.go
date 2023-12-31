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
	FavoriteVideo(req model.FavoriteMessage, reply *model.Response) error
	FavoriteList(userIDToken model.UserIdToken, reply *model.VideoListResponse) error
}
type FavoriteServiceImpl struct {
	UserRepo     *db.MySQLUserRepository
	VideoRepo    *db.MySQLVideoRepository
	FavoriteRepo *db.MySQLFavoriteRepository
}

func (s *FavoriteServiceImpl) FavoriteVideo(req model.FavoriteMessage, reply *model.Response) error {
	userId, err := s.UserRepo.GetUserId(req.Token)
	if err != nil {
		*reply = model.Response{
			StatusCode: 1,
			StatusMsg:  "user didn't exist",
		}
		return fmt.Errorf("user didn't uploaded")
	}
	authorId, err := s.VideoRepo.GetAuthorIdByVideoId(req.VideoId)

	//数据库存入点赞数据
	if req.ActionType == 1 {
		err = s.FavoriteRepo.AddFavorite(userId, req.VideoId)
		if err != nil {
			*reply = model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			}
			return nil
		}
	} else if req.ActionType == 2 {
		err = s.FavoriteRepo.RemoveFavorite(userId, req.VideoId)
		if err != nil {
			*reply = model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			}
			return nil
		}
	}

	//点赞用户点赞数更新
	err = s.UserRepo.UpdateFavoriteCount(req.Token, req.ActionType)
	if err != nil {
		*reply = model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		}
		return nil
	}

	//被点赞用户被点赞数更新
	err = s.UserRepo.UpdateTotalFavorited(authorId, req.ActionType)
	if err != nil {
		*reply = model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		}
		return nil
	}

	//被点赞视频点赞数更新
	err = s.VideoRepo.UpdateFavoriteCount(req.VideoId, req.ActionType)
	if err != nil {
		*reply = model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		}
		return nil
	}
	*reply = model.Response{
		StatusCode: 0,
	}
	return nil
}

func (s *FavoriteServiceImpl) FavoriteList(userIDToken model.UserIdToken, reply *model.VideoListResponse) error {
	_, err := s.UserRepo.GetUserId(userIDToken.Token)
	if err != nil {
		reply = nil
		return fmt.Errorf("user doesn't exist")
	}
	VideoIds, err := s.FavoriteRepo.GetFavoriteVideoIdByUserId(userIDToken.UserId)
	if err != nil {
		*reply = model.VideoListResponse{
			Response: model.Response{
				StatusCode: 0,
			},
			VideoList: nil,
		}
		return nil
	}
	Videos := make([]model.Video, len(VideoIds))
	for i := 0; i < len(VideoIds); i++ {
		video, err := s.VideoRepo.GetVideoByVideoId(VideoIds[i])
		Videos[i] = *video
		if err != nil {
			*reply = model.VideoListResponse{
				Response: model.Response{
					StatusCode: 0,
				},
				VideoList: nil,
			}
			return nil
		}
		Videos[i].IsFavorite = true
	}

	*reply = model.VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: Videos,
	}
	return nil
}

func RunFavoriteServer() {
	// 创建服务实例
	favoriteService := &FavoriteServiceImpl{
		UserRepo:     db.NewMySQLUserRepository(),
		VideoRepo:    db.NewMySQLVideoRepository(),
		FavoriteRepo: db.NewMySQLFavoriteRepository(),
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
