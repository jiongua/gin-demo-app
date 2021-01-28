package api

import (
	"encoding/json"
	"fmt"
	"gin_demo/interal/client"
	"gin_demo/interal/entity"
	"gin_demo/interal/form"
	"gin_demo/interal/task"
	"gin_demo/interal/value"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"time"
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
		//<- 写入消息队列，消费者将该回答设置为所属问题订阅者的通知消息
		NotifyNewAnswer(answer, c.Copy())
		ReporterAnswer(answer, c.Copy())
		//eventURL := fmt.Sprintf("%s/%d", c.Request.URL.Path, answer.AnswerID)
		//log.Debugf("eventURL: %s\n", eventURL)
		//task.NotifyTask.AddTask(mq.NewAnswerNotify(
		//	answer.AuthorID,
		//	entity.ActionAnswer.ID,
		//	entity.QuestionResource,
		//	eventURL,
		//	answer.QuestionRefer))

		//task.AsyncPublishNotify(mq.NewAnswerNotify(
		//	answer.AuthorID,
		//	entity.ActionAnswer.ID,
		//	entity.QuestionResource,
		//	eventURL,
		//	answer.QuestionRefer))

		//if err != nil {
		//	log.WithFields(logrus.Fields{
		//		"author": answer.AuthorID,
		//		"actions": entity.MappingMessage[entity.ActionAnswer.ID],
		//		"resourceURL": entity.QuestionResource,
		//		"questionID": answer.QuestionRefer,
		//	}).Error("publish subscribe_notify error")
		//}
		c.JSON(http.StatusOK, gin.H{
			"message": "成功提交回答",
		})
	})
}
func NotifyNewAnswer(answer entity.Answer, c *gin.Context) {
	go func(answer entity.Answer) {
		message := entity.AnswerNotify{
			NotifyMeta: entity.NotifyMeta{
				SenderID:   answer.AuthorID,
				CreatedAt:  time.Now(),
			},
			QuestionID: answer.QuestionRefer,
			QuestionURL: fmt.Sprintf("%s/%d", c.Request.URL.Path, answer.AnswerID),
		}
		//data, _ := json.Marshal(message)
		//note:
		//publishing from multiple goroutines on the same channel is not safe.
		//Avoid publishing on the same channel from multiple threads/goroutines/processes
		log.Debugf("send answer notify to AnswerNotifyChan[%d]\n", cap(task.AnswerNotifyChan))
		task.AnswerNotifyChan <-message
	}(answer)
}

//ReporterAnswer 上报用户行为记录
func ReporterAnswer(answer entity.Answer, c *gin.Context) {
	go func(answer entity.Answer) {
		event := struct {
			UserID     uuid.UUID
			QuestionID int
			ActionID   int
			CreatedAt  int
		}{
			UserID: answer.AuthorID,
			QuestionID: answer.QuestionRefer,
			ActionID: value.ANSWER,
			CreatedAt: answer.Created,
		}
		data, err := json.Marshal(event)
		if err != nil {
			log.Error(err.Error())
		}else {
			log.Debugf("send user_action event to EventUserActionChain\n")
			client.ReportUserAction("user_action_collection", data)
		}
	}(answer)
}

//GetAnswer 获取一条回答
//GET /question/:qid/answer
func GetAnswer(router *gin.RouterGroup)  {
	router.GET("/:qid/answer/:answer_id", func(c *gin.Context) {
		
	})
}
