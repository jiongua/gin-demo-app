package query

import (
	"fmt"
	"gin_demo/interal/entity"
	"github.com/sirupsen/logrus"
)

//AnswersByQuestionID 根据问题ID查询回答
func AnswersByQuestionID(qid int, offset int, order string) (answers entity.Answers, err error) {
	//select * from answers where question_refer = qid  ORDER BY created OFFSET 0 LIMIT 10
	//获取10个回答
	order = fmt.Sprintf("%s DESC", order)
	if err := entity.Db().Limit(10).Offset(offset).Order(order).Where("question_refer = ?", qid).Find(&answers).Error; err != nil {
		log.WithFields(logrus.Fields{
			"questionID": qid,
			"order": order,
		}).Error("get answers from database error")
	}
	return
}
