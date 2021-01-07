package entity

import (
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type (
	Topic struct {
		TopicID int `gorm:"primaryKey"`
		Name    string
	}

	Question struct {
		QuestionID    int `gorm:"primaryKey;autoIncrement"`
		OldQuestionID int	`json:"-"`
		TopicID       int
		Title         string
		AuthorID      uuid.UUID `gorm:"index;type:uuid"`
		Created       int       `gorm:"autoCreateTime"`
		Updated       int       `gorm:"autoUpdateTime"`
		Deleted       gorm.DeletedAt `json:"-"`
		IsAnonymity   bool	`gorm:"default:false"`
		RealAuthorID  uuid.UUID `gorm:"-" json:"-"`
		AnswerCount int `gorm:"default:100"`
		AttentionCount int `gorm:"default:100"`
	}
)

//NewQuestion return a new question
func NewQuestion(topicID int, title string, authorID uuid.UUID, isAnonymity bool) *Question {
	if isAnonymity {
		authorID = uuid.Nil
	}
	result := &Question{
		TopicID:       topicID,
		Title:         title,
		AuthorID:      authorID,
		IsAnonymity: isAnonymity,
	}
	return result
}

func (q *Question) CreateQuestion() error {
	return Db().Create(q).Error
}
func (t *Topic) CreateTopic() error {
	return Db().Create(t).Error
}

func FirstOrCreateTopic(topic *Topic) *Topic {
	var result = Topic{}
	if err := Db().Where("name = ?", topic.Name).First(&result).Error; err == nil {
		return &result
	}else {
		log.Debug("topic not found, create one")
		if err := topic.CreateTopic(); err != nil {
			return nil
		}
		return topic
	}
}


