package main

import (
	"gin_demo/interal/server"
	"gin_demo/interal/task"
	"github.com/gin-gonic/gin"
)


func main() {
	router := gin.New()
	router.Use(server.Logger())
	server.RegisterRoutes(router)
	task.StartAnswerNotifyWorker()
	task.StartVoteAnswerNotifyWorker()
	router.Run()
}
