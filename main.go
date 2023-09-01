package main

import (
	"github.com/AdaQiao/simpleDouyin/db"
	"github.com/AdaQiao/simpleDouyin/service"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatal("初始化数据库连接失败:", err)
	}
	defer db.CloseDB()
	go service.RunFeedServer()
	go service.RunPublishServer()
	go service.RunMessageServer()
	go service.RunUserServer()
	go service.RunFavoriteServer()
	go service.RunCommentServer()
	go service.RunRelationServer()
	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
