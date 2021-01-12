package entity

import (
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Answer struct {
	AnswerID      int `gorm:"primaryKey" json:"id"`
	OldAnswerID   int	`json:"-"`
	QuestionRefer int      `gorm:"index" json:"-"`
	Question      Question `json:"-" gorm:"foreignKey:QuestionRefer;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Excerpt       string	`json:"excerpt"`
	Content       string	`json:"-"`
	AuthorID      uuid.UUID `gorm:"index;type:uuid;" json:"author_id"`
	CommentCount  int		`json:"comment_count"`
	VoteCount     int		`json:"vote_count"`
	Created       int `gorm:"autoCreateTime" json:"created_time"`
	Updated       int `gorm:"autoUpdateTime" json:"update_time"`
	Deleted       gorm.DeletedAt 	`json:"-"`
}

type Answers []Answer

//创建新回答
func (a *Answer) Create() error {
	return Db().Create(a).Error
}