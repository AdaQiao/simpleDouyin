package service

import (
	"fmt"
	"github.com/AdaQiao/simpleDouyin/db"
	"github.com/AdaQiao/simpleDouyin/model"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"net"
	"net/rpc"
	"os"
)

type PublishService interface {
	Publish(req model.UploadViewReq, reply *model.Response) error
	PublishList(userIDToken model.UserIdToken, reply *model.VideoListResponse) error
	UploadVideoToOSS(file model.FilenameAndFilepath, reply *model.CoverAndVideoURL) error
}

type PublishServiceImpl struct {
	UserRepo     *db.MySQLUserRepository
	VideoRepo    *db.MySQLVideoRepository
	FavoriteRepo *db.MySQLFavoriteRepository
}

func (s *PublishServiceImpl) Publish(req model.UploadViewReq, reply *model.Response) error {
	user, err := s.UserRepo.GetUserByToken(req.Token)
	if err != nil {
		return fmt.Errorf("user doesn't exist")
	}
	s.UserRepo.UpdateWorkCount(req.Token)
	newVideo := model.Video{
		Author:        *user,
		PlayUrl:       req.ViewUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		CoverUrl:      req.CoverUrl,
		IsFavorite:    false,
		Title:         req.Title,
	}
	err = s.VideoRepo.CreateVideo(newVideo, req.Token)
	if err != nil {
		return err
	}
	*reply = model.Response{
		StatusCode: 0,
		StatusMsg:  req.Title + " uploaded successfully",
	}
	return nil
}
func (s *PublishServiceImpl) PublishList(userIDToken model.UserIdToken, reply *model.VideoListResponse) error {
	_, err := s.UserRepo.GetUserId(userIDToken.Token)
	if err != nil {
		reply = nil
		return fmt.Errorf("user doesn't exist")
	}
	Videos, err := s.VideoRepo.GetVideoById(userIDToken.UserId)
	for i := 0; i < len(Videos); i++ {
		user, _ := s.UserRepo.GetUserByToken(userIDToken.Token)
		Videos[i].Author = *user
	}
	if err != nil {
		*reply = model.VideoListResponse{
			Response: model.Response{
				StatusCode: 0,
			},
			VideoList: nil,
		}
		return nil
	}
	//检查是否点赞
	for i := 0; i < len(Videos); i++ {
		isLike, err := s.FavoriteRepo.CheckFavorite(userIDToken.UserId, Videos[i].Id)
		if err != nil {
			fmt.Println("查询是否点赞失败:", err)
		}
		Videos[i].IsFavorite = isLike
	}

	*reply = model.VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: Videos,
	}
	return nil
}

func (s *PublishServiceImpl) UploadVideoToOSS(file model.FilenameAndFilepath, reply *model.CoverAndVideoURL) error {
	accessKeyID := os.Getenv("OSS_ID")
	accessKeySecret := os.Getenv("OSS_PASSWORD")
	fmt.Println("id:", os.Getenv("OSS_ID"))
	fmt.Println("pw:", os.Getenv("OSS_PASSWORD"))
	endpoint := "oss-cn-beijing.aliyuncs.com"
	bucketName := "simple-douyin"
	// 创建 OSS 客户端实例
	client1, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		fmt.Println("Error creating OSS client:", err)
		return err
	}

	// 获取存储空间（Bucket）实例
	bucket, err := client1.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error getting OSS bucket:", err)
		return err
	}

	// 要上传的文件路径
	filePath := file.FilePath

	// 打开要上传的文件
	fileToUpload, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}

	// 设置上传到 OSS 的文件名
	objectKey := file.FileName

	// 开始上传文件
	err = bucket.PutObject(objectKey, fileToUpload)
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return err
	}

	// 获取存储的网址
	objectURL := fmt.Sprintf("https://%s.%s/%s", bucketName, endpoint, objectKey)

	*reply = model.CoverAndVideoURL{
		CoverURL: "https://simple-douyin.oss-cn-beijing.aliyuncs.com/douyin.png",
		VideoURL: objectURL,
	}
	return nil

}

func RunPublishServer() {
	// 创建服务实例

	publishService := &PublishServiceImpl{
		UserRepo:  db.NewMySQLUserRepository(),
		VideoRepo: db.NewMySQLVideoRepository(),
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
