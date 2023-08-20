package service

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/AdaQiao/simpleDouyin/db"
	"github.com/AdaQiao/simpleDouyin/model"
)

type CommentService interface {
	Comment(req model.CommentActionRequest, reply *model.CommentActionResponse) error
	CommentList(videoId int64, reply *model.CommentListResponse) error
}
type CommentServiceImpl struct {
	UserRepo    *db.MySQLUserRepository
	VideoRepo   *db.MySQLVideoRepository
	CommentRepo *db.MySQLCommentRepository
}

func (s *CommentServiceImpl) Comment(req model.CommentActionRequest, reply *model.CommentActionResponse) error {
	userId, err := s.UserRepo.GetUserId(req.Token)
	if err != nil {
		*reply = model.CommentActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "user didn't exist",
			},
		}
		return fmt.Errorf("user didn't uploaded")
	}

	// 数据库存入评论数据
	if req.ActionType == 1 {
		err = s.CommentRepo.AddComment(userId, req.VideoId, req.CommentId, req.CommentText)
		if err != nil {
			*reply = model.CommentActionResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			}
			return nil
		}
	} else if req.ActionType == 2 {
		err = s.CommentRepo.RemoveComment(userId, req.VideoId, req.CommentId)
		if err != nil {
			*reply = model.CommentActionResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			}
			return nil
		}
	}

	//被评论视频评论数更新
	err = s.VideoRepo.UpdateCommentCount(req.VideoId, req.ActionType)
	if err != nil {
		*reply = model.CommentActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		}
		return nil
	}
	*reply = model.CommentActionResponse{
		Response: model.Response{
			StatusCode: 0,
		},
	}
	return nil
}

func (s *CommentServiceImpl) CommentList(videoId int64, reply *model.CommentListResponse) error {
	CommentIds, err := s.CommentRepo.GetCommentIdByVideoId(videoId)
	if err != nil {
		*reply = model.CommentListResponse{
			Response: model.Response{
				StatusCode: 0,
			},
			CommentList: nil,
		}
		return nil
	}
	Comments := make([]model.Comment, len(CommentIds))
	for i := 0; i < len(CommentIds); i++ {
		comment, err := s.CommentRepo.GetCommentByCommentId(CommentIds[i])
		Comments[i] = *comment
		if err != nil {
			*reply = model.CommentListResponse{
				Response: model.Response{
					StatusCode: 0,
				},
				CommentList: nil,
			}
			return nil
		}
	}

	*reply = model.CommentListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		CommentList: Comments,
	}
	return nil
}

func RunCommentServer() {
	// 创建服务实例
	CommentService := &CommentServiceImpl{
		UserRepo:    db.NewMySQLUserRepository(),
		VideoRepo:   db.NewMySQLVideoRepository(),
		CommentRepo: db.NewMySQLCommentRepository(),
	}

	// 注册RPC服务
	rpc.Register(CommentService)

	// 启动RPC服务器
	listener, err := net.Listen("tcp", "127.0.0.1:9095")
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
