package service

import (
	"fmt"

	"github.com/AdaQiao/simpleDouyin/db"
	"github.com/AdaQiao/simpleDouyin/model"
)

type CommentService interface {
	Comment(req model.CommentActionRequest, reply *model.CommentActionResponse) error
	FavoriteList(userIDToken model.UserIdToken, reply *model.CommentListResponse) error
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
			Comment: model.Comment{},
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
				Comment: model.Comment{},
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
				Comment: model.Comment{},
			}
			return nil
		}
	}

	authorId, err := s.VideoRepo.GetAuthorIdByVideoId(req.VideoId)
}
