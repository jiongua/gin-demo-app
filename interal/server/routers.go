package server

import (
	"gin_demo/interal/api"
	"github.com/gin-gonic/gin"
)

//registerRoutes 注册所有router handler
func RegisterRoutes(router *gin.Engine) {
	// JSON-REST API Version 1
	v1 := router.Group("/api/v1")
	question := v1.Group("/question")
	{
		api.CreateQuestion(question)
		api.GetQuestion(question)
		//api.UpdateQuestion(v1)
		//api.DeleteQuestion(v1)
		//more...
		api.CreateAnswer(question)
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
	notification := v1.Group("/notification")
	{
		api.NotificationAttentions(notification)
		api.NotificationVotes(notification)
		//api.NotificationCommentAnswers(notification)
		//api.NotificationVoteComments(notification)
		api.NotificationFollowers(notification)
	}
}


