package main

import (
	"github.com/cpl/simple-demo/db"
	"github.com/cpl/simple-demo/service"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	db.InitDB()
	go service.RunMessageServer()
	go service.RunUserServer()
	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	err := db.GetDB().Close()
	if err != nil {
		log.Fatal("关闭数据库连接失败:", err)
	}
}
