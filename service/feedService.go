package service

import (
	"fmt"
	"github.com/AdaQiao/simpleDouyin/db"
	"github.com/AdaQiao/simpleDouyin/model"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"time"
)

type FeedService interface {
	GetVideoList(feedReq model.FeedRequest, reply *model.FeedResponse) error
}

type FeedServiceImpl struct {
	UserRepo     *db.MySQLUserRepository
	VideoRepo    *db.MySQLVideoRepository
	FavoriteRepo *db.MySQLFavoriteRepository
}

func (s *FeedServiceImpl) GetVideoList(feedReq model.FeedRequest, reply *model.FeedResponse) error {
	curTime, err := strconv.ParseInt(feedReq.LatestTime, 10, 64)
	if err != nil || curTime == int64(0) {
		curTime = time.Now().Unix()
	}
	videos, nextTime, tokens, err := s.VideoRepo.GetVideosByTimestamp(curTime)
	if err != nil {
		log.Println("获取视频流失败:", err)
		return err
	}
	for i := 0; i < len(videos); i++ {
		user, _ := s.UserRepo.GetUserByToken(tokens[i])
		videos[i].Author = *user
	}
	//如果用户已登录，查询点赞状态
	if feedReq.Token != "" {
		userId, err := s.UserRepo.GetUserId(feedReq.Token)
		if err != nil {
			fmt.Println("查询用户id失败:", err)
		}
		for i := 0; i < len(videos); i++ {
			isLike, err := s.FavoriteRepo.CheckFavorite(userId, videos[i].Id)
			if err != nil {
				fmt.Println("查询是否点赞失败:", err)
			}
			videos[i].IsFavorite = isLike
		}
	}
	*reply = model.FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  nextTime,
	}
	return nil
}

func RunFeedServer() {
	// 创建服务实例

	feedService := &FeedServiceImpl{
		UserRepo:  db.NewMySQLUserRepository(),
		VideoRepo: db.NewMySQLVideoRepository(),
	}

	// 注册RPC服务
	rpc.Register(feedService)

	// 启动RPC服务器
	listener, err := net.Listen("tcp", "127.0.0.1:9093")
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
