package entity

import (
	"github.com/sirupsen/logrus"
)

//新回答通知
type AnswerNotify struct {
	NotifyMeta
	QuestionID     int
	QuestionURL    string	`gorm:"varchar(100)"`
}

func UpdateRead(id int)  {
	err := Db().Model(&AnswerNotify{}).Where("id=", id).UpdateColumn("is_read = ?", true).Error
	if err != nil {
		log.WithFields(logrus.Fields{
			"id": id,
			"error": err.Error(),
		}).Error("UpdateRead error!")
	}
}





