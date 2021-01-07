package server

import (
	"gin_demo/interal/api"
	"github.com/gin-gonic/gin"
)

//registerRoutes 注册所有router handler
func registerRoutes(router *gin.Engine) {
	// JSON-REST API Version 1
	v1 := router.Group("/api/v1")
	{
		api.CreateQuestion(v1)
		api.GetQuestion(v1)
		//api.UpdateQuestion(v1)
		//api.DeleteQuestion(v1)
		//more...
		//api.CreateAnswer(v1)
	}
	people := v1.Group("/people/:uid")
	{
		api.HomePage(people)
		api.Questions(people)
		api.Answers(people)
		api.Following(people)
		api.Followers(people)
		api.FollowingQuestions(people)
		api.Votes(people)
	}

}


