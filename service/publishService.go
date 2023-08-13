package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"net"
	"net/rpc"
	"os"
)

type PublishService interface {
	Publish(req model.UploadViewReq, reply *model.Response) error
	PublishList(userIDToken model.UserIdToken, reply *model.VideoListResponse) error
	UploadVideoToOSS(file model.FilenameAndFilepath, reply *string) error
}

type PublishServiceImpl struct {
	UserRepo  *db.MySQLUserRepository
	VideoRepo *db.MySQLVideoRepository
}

func (s *PublishServiceImpl) Publish(req model.UploadViewReq, reply *model.Response) error {
	user, err := s.UserRepo.GetUser(req.Token)
	if err != nil {
		return fmt.Errorf("user doesn't exist")
	}
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
		return fmt.Errorf("user doesn't exist")
	}
	Videos, err := s.VideoRepo.GetVideoById(userIDToken.UserId)
	if err != nil {
		*reply = model.VideoListResponse{
			Response: model.Response{
				StatusCode: 0,
			},
			VideoList: nil,
		}
		return nil
	}
	*reply = model.VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: Videos,
	}
	return nil
}

func (s *PublishServiceImpl) UploadVideoToOSS(file model.FilenameAndFilepath, reply *string) error {
	accessKeyID := "LTAI5t7jPFXhiXgckbXHeWeR"
	accessKeySecret := "imAsfE1B4MF7VZTcgH6puYngVm0IwN"
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
	//defer file.Close()
	defer func() {
		fileToUpload.Close()
		// 删除本地保存的视频文件
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("Error deleting local file:", err)
		}
	}()

	// 设置上传到 OSS 的文件名
	objectKey := file.FileName
	fmt.Println("Final name: ", file.FileName)

	// 开始上传文件
	err = bucket.PutObject(objectKey, fileToUpload)
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return err
	}

	// 获取存储的网址
	objectURL, err := bucket.SignURL(objectKey, oss.HTTPGet, 3600)
	*reply = objectURL
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
