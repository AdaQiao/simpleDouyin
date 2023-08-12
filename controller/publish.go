package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/rpc"
	"path/filepath"
)

var repo *db.MySQLUserRepository = db.NewMySQLUserRepository()

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("title")
	fmt.Println(title)
	user, err := repo.GetUser(token)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 连接到远程RPC服务器
	client, err := rpc.Dial("tcp", "127.0.0.1:9092")
	if err != nil {
		log.Fatal("RPC连接失败：", err)
	}

	// 调用远程注册方法
	var reply model.Response
	err = client.Call("PublishServiceImpl.Publish", model.UploadViewReq{Title: title, Token: token, ViewUrl: saveFile, CoverUrl: "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"}, &reply)
	if err != nil {
		log.Fatal("调用远程注册方法失败：", err)
	}
	c.JSON(http.StatusOK, reply)

	/*c.JSON(http.StatusOK, model.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})*/
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	token := c.PostForm("token")
	// 连接到远程RPC服务器
	client, err := rpc.Dial("tcp", "127.0.0.1:9092")
	if err != nil {
		log.Fatal("RPC连接失败：", err)
	}
	user_id := c.PostForm("user_id")
	// 调用远程注册方法
	var reply model.VideoListResponse
	err = client.Call("PublishServiceImpl.PublishList", model.UserIdToken{
		Token:  token,
		UserId: user_id,
	}, &reply)
	if err != nil {
		log.Fatal("调用远程注册方法失败：", err)
	}
	c.JSON(http.StatusOK, reply)
	/*c.JSON(http.StatusOK, model.VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: model.DemoVideos,
	})*/
}
