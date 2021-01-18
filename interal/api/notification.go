package api

import (
	"gin_demo/interal/query"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//NotificationAttentions 关注的通知:新回答、回答的评论
//GET /notification/attentions?offset=0&limit=10
func NotificationAttentions(router *gin.RouterGroup)  {
	router.GET("/attentions", func(c *gin.Context) {
		offset, err := strconv.Atoi(c.Param("offset"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		limit, err := strconv.Atoi(c.Param("limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		userID := getUserFromToken()
		result := query.GetAnswerNotifyByReceiverID(userID, offset, limit)
		c.JSON(200, result)
	})
}

//NotificationFollowers 有新的粉丝通知
//GET /notification/followers?offset=0&limit=10
func NotificationFollowers(router *gin.RouterGroup)  {
	router.GET("/followers", func(c *gin.Context) {
		//offset, err := strconv.Atoi(c.Param("offset"))
		//if err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"error": err.Error(),
		//	})
		//	return
		//}
		//limit, err := strconv.Atoi(c.Param("limit"))
		//if err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"error": err.Error(),
		//	})
		//	return
		//}
		//userID := getUserFromToken()


	})
}

//NotificationVotes 点赞通知，回答和评论被点赞
//GET /notification/votes?offset=0&limit=10
func NotificationVotes(router *gin.RouterGroup) {
	router.GET("/votes", func(c *gin.Context) {
		//offset, err := strconv.Atoi(c.Param("offset"))
		//if err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"error": err.Error(),
		//	})
		//	return
		//}
		//limit, err := strconv.Atoi(c.Param("limit"))
		//if err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"error": err.Error(),
		//	})
		//	return
		//}
		//userID := getUserFromToken()



	})
}
