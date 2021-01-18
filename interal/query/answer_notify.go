package query

import (
	"gin_demo/interal/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func GetAnswerNotifyByReceiverID(receiverID uuid.UUID, offset, limit int)  []entity.AnswerNotify {
	var result []entity.AnswerNotify
	err := entity.Db().Model(&entity.AnswerNotify{}).Offset(offset).Limit(limit).Where("receiver_id=?", receiverID).Find(&result).Error
	if err != nil {
		log.WithFields(logrus.Fields{
			"receiverID": receiverID,
			"error": err.Error(),
		}).Error("GetAnswerNotifyByReceiverID error!")
	}
	return result
}