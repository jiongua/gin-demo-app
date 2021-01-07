package form

import (
	"github.com/sirupsen/logrus"
)

//for query binding
type QuestionQuery struct {
	//QuestionID int 	`uri:"qid" binding:"required"`
	AnswerOffset *int	`form:"offset" binding:"required"`
	AnswerLimit int 	`form:"limit" binding:"-"`
	AnswerOrder string 	`form:"order" binding:"-"`
}

//for create binding
type QuestionCreate struct {
	TopicName string `json:"topic_name"`
	Title string `json:"title"`
	IsAnonymity bool `json:"is_anonymity"`
}

//OrderFormat 将查询参数order转换为通过数据库查询的字段
func (f *QuestionQuery) OrderFormat() {
	switch f.AnswerOrder {
	case "byTime", "created":
		f.AnswerOrder = "created"
	case "byHot", "vote_count":
		f.AnswerOrder = "vote_count"
	default:
		log.WithFields(logrus.Fields{
			"order": f.AnswerOrder,
			"default": "created",
		}).Warn("invalid order, set default")
		f.AnswerOrder = "created"
	}
}




