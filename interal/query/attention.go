package query

import (
	"gin_demo/interal/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func GetAllAttentionUsersOrNil(questionID int) []uuid.UUID {
	var users []uuid.UUID
	err := entity.Db().Select("user_id").Where("question_refer = ?", questionID).Find(&users).Error
	if err != nil {
		log.WithFields(logrus.Fields{
			"question_id": questionID,
		}).Error("查询问题关注者错误")
		return  nil
	}
	return users
}
