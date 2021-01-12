package api

import (
	"gin_demo/interal/entity"
	"gin_demo/interal/form"
	glog "gin_demo/interal/log"
	"gin_demo/interal/query"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var log = glog.Log

//GetQuestion 获取问题及答案
// GET /api/v1/question/:qid/answers
func GetQuestion(router *gin.RouterGroup) {
	router.GET("/:qid/answers", func(c *gin.Context) {
		qid, _ := strconv.Atoi(c.Param("qid"))
		//limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		//offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		//order := c.DefaultQuery("order", "created")

		var f form.QuestionQuery
		if err := c.ShouldBind(&f); err != nil {
			log.Errorf("question bind err: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		f.OrderFormat()
		log.WithFields(logrus.Fields{
			"qid": qid,
			"limit": f.AnswerLimit,
			"offset": *f.AnswerOffset,
			"order": f.AnswerOrder,
		}).Debug("get question request")

		result, err := query.QuestionByID(qid, *f.AnswerOffset, f.AnswerLimit, f.AnswerOrder)
		if err != nil {
			log.Debugf("QuestionByID error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": result})
	})
}

//CreateQuestion 提问
//POST /api/v1/question
func CreateQuestion(router *gin.RouterGroup)  {
	router.POST("/", func(c *gin.Context) {
		var f form.QuestionCreate
		if err := c.BindJSON(&f); err!=nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		topic := entity.FirstOrCreateTopic(&entity.Topic{
			Name:    f.TopicName,
		})
		if topic == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建topic失败"})
			return
		}
		authorID := entity.Super1.ID
		var q = entity.NewQuestion(topic.TopicID, f.Title, authorID, f.IsAnonymity)
		if err := q.CreateQuestion(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//创建者自动关注问题
		a := entity.Attention{
			UserID:        authorID,
			QuestionRefer: q.QuestionID,
		}
		if err := a.Create(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "问题创建成功"})
	})
}

//UpdateQuestion 修改问题
//PUT /api/v1/question/:qid
func UpdateQuestion(router *gin.RouterGroup)  {
	router.PUT("/:qid", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"url": c.Request.URL, "method": c.Request.Method})
	})
}

//DeleteAnswer 删除问题
//Delete /api/v1/question/:qid
func DeleteQuestion(router *gin.RouterGroup)  {
	router.DELETE("/:qid", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"url": c.Request.URL, "method": c.Request.Method})
	})
}