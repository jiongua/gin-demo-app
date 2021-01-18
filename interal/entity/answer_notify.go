package entity

import (
	"gin_demo/interal/query"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

//新回答通知
type AnswerNotify struct {
	NotifyMeta
	QuestionID     int
	QuestionURL    string	`gorm:"varchar(100)"`
}

// CopyAndCreateAll 生成每个关注人的通知，除发布人和禁用通知的接收者
func (m *AnswerNotify) CopyAndCreateAll()  {
	sender, err := query.GetUserByID(m.SenderID)
	if err != nil {
		log.Errorf("get SenderName by SenderID[%s] error: %s\n", m.SenderID, err.Error())
		return
	}
	m.SenderName = sender.Name

	receiverIDs := query.GetAllAttentionUsersOrNil(m.QuestionID)
	batchSize := len(receiverIDs)
	notifies := make([]AnswerNotify, 0, batchSize)
	for i, id := range receiverIDs {
		//过滤接收者是SenderID
		if uuid.Equal(m.SenderID, id) {
			continue
		}
		//todo 过滤屏蔽通知的用户
		notifyItem := AnswerNotify{}
		err := copier.Copy(notifyItem, m)
		if err != nil {
			log.Errorf("copier error: %s\n", err.Error())
			continue
		}
		notifyItem.ReceiverID = id
		notifies[i] = notifyItem
	}
	if batchSize > 0 {
		Db().CreateInBatches(&notifies, batchSize)
	}
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





