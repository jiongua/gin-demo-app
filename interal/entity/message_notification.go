package entity

import (
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

const (
	QuestionResource = uint8(0)
	AnswerResource = uint8(1)
	UserResource = uint8(2)
)
//消息类型表
type MessageTypes struct {
	ID uint8  `gorm:"primaryKey;autoIncrement" json:"id"`
	Describe string `gorm:"varchar(50)" json:"describe"`
}

// 消息共用结构
type NotifyMeta struct {
	ID             uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID       uuid.UUID `gorm:"type:uuid" json:"sender_id"`
	SenderName     string
	SenderURL      string `gorm:"varchar(100)" json:"sender_url"`
	ReceiverID     uuid.UUID `gorm:"type:uuid;index" json:"-"`
	CreatedAt      time.Time
	IsRead        bool 	`gorm:"default:false;index"`
	Deleted 	gorm.DeletedAt `json:"-"`
}

//消息通知表
type MessageNotification struct {
	ID             uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID       uuid.UUID `gorm:"type:uuid" json:"sender_id"`
	SenderName     uuid.UUID
	SenderURL      string `gorm:"varchar(100)" json:"sender_url"`
	Action         uint8	`json:"-"`
	ActionDescribe string `gorm:"varchar(50)" json:"action_describe"`
	ReceiverID     uuid.UUID `gorm:"type:uuid" json:"-"`
	ResourceType     uint8
	ResourceURL 	string	`gorm:"varchar(100)"`
	EventURL    string `gorm:"varchar(50)"`
	CreatedAt      time.Time
	IsRead        bool 	`gorm:"default:false"`
	Deleted 	gorm.DeletedAt `json:"-"`
}



type EventEntry struct {
	ID int
	Kind uint8		//消息类型: 消息通知类型、时间线
	Name string		//消息名称: 通知(回答、点赞、收藏、点赞、评论等)、时间线(提问、回答了、评论了、点赞了..)
	ResourceID int	//目的资源对象ID
}

var ActionAnswer = MessageTypes{
	ID:       0,
	Describe: "回答",
}

var ActionCommentAnswer = MessageTypes{
	ID:       1,
	Describe: "评论了回答",
}

var ActionVoteAnswer = MessageTypes{
	ID:       2,
	Describe: "赞同了回答",
}

var ActionFollow = MessageTypes{
	ID:       3,
	Describe: "关注了",
}

type MappingMessageType map[uint8]string
var MappingMessage = MappingMessageType{}

func (m MappingMessageType) GetKey(key uint8) string {
	return m[key]
}

func (m MappingMessageType) SetKey(key uint8, val string) {
	m[key] = val
}

func CreateDefaultMessageTypes()  {
	if result := FirstOrCreateMessageTypes(&ActionAnswer); result != nil {
		ActionAnswer = *result
		MappingMessage.SetKey(result.ID, result.Describe)
	}
	if result := FirstOrCreateMessageTypes(&ActionCommentAnswer); result != nil {
		ActionCommentAnswer = *result
		MappingMessage.SetKey(result.ID, result.Describe)
	}
	if result := FirstOrCreateMessageTypes(&ActionVoteAnswer); result != nil {
		ActionVoteAnswer = *result
		MappingMessage.SetKey(result.ID, result.Describe)
	}
	if result := FirstOrCreateMessageTypes(&ActionFollow); result != nil {
		ActionFollow = *result
		MappingMessage.SetKey(result.ID, result.Describe)
	}
}

func FirstOrCreateMessageTypes(m *MessageTypes) *MessageTypes {
	result := MessageTypes{}
	if err := Db().Where("id = ?", m.ID).First(&result).Error; err==nil {
		//找到id为m.ID的数据
		return &result
	} else if err := m.CreateMessageType(); err !=nil {
		//插入失败
		return nil
	}
	//插入m成功
	return m
}

//Create 创建一个消息通知类型
func (t *MessageTypes) CreateMessageType() error{
	return Db().Create(t).Error
}

func (m *MessageNotification) Create() error{
	return Db().Create(m).Error
}

func (m *MessageNotification) ReadOne() error{
	//UPDATE MessageNotification SET is_read = true WHERE id =
	return Db().Where("id = ", m.ID).UpdateColumn("IsRead", true).Error
}

func (m *MessageNotification) ReadAll() error{
	//UPDATE MessageNotification SET is_read = true WHERE is_read = false
	return Db().Not("is_read", true).Select("is_read").UpdateColumns(&MessageNotification{
		IsRead: true,
	}).Error
}


