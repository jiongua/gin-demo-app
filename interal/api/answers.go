package api

import (
	"gin_demo/interal/entity"
	"gin_demo/interal/form"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//CreateAnswer 新建回答
//POST /question/:qid/answer
func CreateAnswer(router *gin.RouterGroup)  {
	router.POST("/:qid/answer", func(c *gin.Context) {
		qid, err := strconv.Atoi(c.Param("qid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var f form.Answer
		if err:=c.ShouldBind(&f); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		content := f.Content
		excerpt := func() string {
			s := []rune(content)
			tail := len(s)
			if tail >= 50 {
				tail = 50
			}
			return string(s[:tail])
			//var result []rune
			//runeCount := 0
			//for _, r := range content {
			//	result = append(result, r)
			//	runeCount++
			//	if runeCount >= 50 {
			//		break
			//	}
			//}
			//return string(result)
		}()
		answer := entity.Answer{
			QuestionRefer: qid,
			Excerpt:       excerpt,
			Content:       content,
			AuthorID:      getUserFromToken(),
		}
		if err := answer.Create(); err !=nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		//发布消息
		// <- 写入消息队列，消费者将该消息写入回答者的粉丝时间线缓存
		// <- 写入消息队列，消费者将该回答设置为所属问题订阅者的通知消息
		//message := &MessageSendMeta{
		//	SenderID: answer.AuthorID,
		//	ActionID: entity.ActionAnswer.ID,
		//	ResourceID: answer.QuestionRefer,
		//}
		//mq.Publish(message)
		c.JSON(http.StatusOK, gin.H{
			"message": "成功提交回答",
		})
	})
}

//GetAnswer 获取一条回答
//GET /question/:qid/answer
func GetAnswer(router *gin.RouterGroup)  {
	router.GET("/:qid/answer/:answer_id", func(c *gin.Context) {
		
	})
}
