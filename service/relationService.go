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
	FollowList(req model.RelationListMessage, reply *model.UserListResponse) error
	FollowerList(req model.RelationListMessage, reply *model.UserListResponse) error
	FriendList(req model.RelationListMessage, reply *model.UserListResponse) error
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
		err = s.RelationRepo.AddFollower(toUserId, userId)
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
		err = s.RelationRepo.RemoveFollower(toUserId, userId)
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

func (s *RelationServiceImpl) FollowList(req model.RelationListMessage, reply *model.UserListResponse) error {
	_, err := s.UserRepo.GetUserId(req.Token)
	if err != nil {
		*reply = model.UserListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "user didn't exist",
			},
			UserList: nil,
		}
		return fmt.Errorf("user didn't uploaded")
	}
	followIds, err := s.RelationRepo.GetFollowById(req.UserId)
	if err != nil {
		*reply = model.UserListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "get followList failed",
			},
			UserList: nil,
		}
		return nil
	}

	var userList = make([]model.User, len(followIds))
	for i := 0; i < len(followIds); i++ {
		tempUser, err := s.UserRepo.GetUserByUserId(followIds[i])
		if err != nil {
			*reply = model.UserListResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  "获取关注用户失败",
				},
				UserList: nil,
			}
			return nil
		}
		tempUser.IsFollow = true
		userList[i] = *tempUser
	}

	*reply = model.UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: userList,
	}
	return nil
}
func (s *RelationServiceImpl) FollowerList(req model.RelationListMessage, reply *model.UserListResponse) error {
	_, err := s.UserRepo.GetUserId(req.Token)
	if err != nil {
		*reply = model.UserListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "user didn't exist",
			},
			UserList: nil,
		}
		return fmt.Errorf("user didn't uploaded")
	}
	followerIds, err := s.RelationRepo.GetFollowerById(req.UserId)
	if err != nil {
		*reply = model.UserListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "get followList failed",
			},
			UserList: nil,
		}
		return nil
	}

	var userList = make([]model.User, len(followerIds))
	for i := 0; i < len(followerIds); i++ {
		tempUser, err := s.UserRepo.GetUserByUserId(followerIds[i])
		if err != nil {
			*reply = model.UserListResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  "获取粉丝用户失败",
				},
				UserList: nil,
			}
			return nil
		}
		tempUser.IsFollow, err = s.RelationRepo.CheckFollow(req.UserId, followerIds[i])
		if err != nil {
			*reply = model.UserListResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  "获取粉丝用户失败",
				},
				UserList: nil,
			}
			return nil
		}
		userList[i] = *tempUser
	}

	*reply = model.UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: userList,
	}
	return nil
}
func (s *RelationServiceImpl) FriendList(req model.RelationListMessage, reply *model.UserListResponse) error {
	_, err := s.UserRepo.GetUserId(req.Token)
	if err != nil {
		*reply = model.UserListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "user didn't exist",
			},
			UserList: nil,
		}
		return fmt.Errorf("user didn't uploaded")
	}
	followerIds, err := s.RelationRepo.GetFollowerById(req.UserId)
	if err != nil {
		*reply = model.UserListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "get followList failed",
			},
			UserList: nil,
		}
		return nil
	}

	var userList []model.User
	for i := 0; i < len(followerIds); i++ {
		tempUser, err := s.UserRepo.GetUserByUserId(followerIds[i])
		if err != nil {
			*reply = model.UserListResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  "获取关注用户失败",
				},
				UserList: nil,
			}
			return nil
		}
		isFollow, err := s.RelationRepo.CheckFollow(req.UserId, followerIds[i])
		if err != nil {
			*reply = model.UserListResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  "检查是否关注失败",
				},
				UserList: nil,
			}
			return nil
		}
		if isFollow {
			tempUser.IsFollow = true
			userList = append(userList, *tempUser)
		}

	}

	*reply = model.UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: userList,
	}
	return nil
}

type RelationServiceImpl struct {
	UserRepo     *db.MySQLUserRepository
	RelationRepo *db.MySQLRelationRepository
}

func RunRelationServer() {
	// 创建服务实例
	relationService := &RelationServiceImpl{
		UserRepo:     db.NewMySQLUserRepository(),
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
