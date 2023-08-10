package main

import (
	"github.com/cpl/simple-demo/db"
	"github.com/cpl/simple-demo/service"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	go service.RunMessageServer()
	go service.RunUserServer()
	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	defer db.DB.Close()
}
