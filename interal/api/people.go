package api

import (
	"gin_demo/interal/entity"
	"gin_demo/interal/query"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func getUidFromUri(c *gin.Context) string {
	return c.Param("uid")
}

//测试用
func getUserFromToken() uuid.UUID {
	return entity.Super1.ID
}

//HomePage 用户的个人主页
//GET /people/:uid -> /people/:uid/
func HomePage(router *gin.RouterGroup) {
	router.GET("/", func(c *gin.Context) {
		//get path parameter
		uid := getUidFromUri(c)

		c.JSON(200, gin.H{"path": c.FullPath(), "uid": uid})
	})
}

//Questions 用户的提问
//GET /people/:uid/questions/?pages=
func Questions(router *gin.RouterGroup)  {
	router.GET("/questions", func(c *gin.Context) {
		//get path parameter
		uid := uuid.FromStringOrNil(getUidFromUri(c))
		pages, _ := strconv.Atoi(c.DefaultQuery("pages", "0"))
		loginID := getUserFromToken()
		isMine := false
		if uuid.Equal(uid, loginID) {
			//查看自己的主页, 可以看到所有匿名的提问
			isMine = true
		}
		log.WithFields(logrus.Fields{
			"uid": uid,
			"pages": pages,
			"loginID": loginID,
			"isMine": isMine,
		}).Debug("query people pages")
		//从question表过滤uid的提问
		result, err := query.GetQuestionsByAuthorID(uid, pages*10, isMine)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": result})
	})
}

//Answers 用户的回答
//GET /people/:uid/answers
func Answers(router *gin.RouterGroup) {

}

//Votes 用户的点赞
//GET /people/:uid/votes
func Votes(router *gin.RouterGroup)  {
	
}
//Following 用户的关注列表
//GET /people/:uid/following
func Following(router *gin.RouterGroup)  {
	
}
//FollowingQuestions 用户关注的问题列表
//GET /people/:uid/following/questions
func FollowingQuestions(router *gin.RouterGroup) {

}
//Followers 用户的粉丝列表
//GET /people/:uid/followers
func Followers(router *gin.RouterGroup)  {

}


