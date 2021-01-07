package entity

import (
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Comment struct {
	CommentID    uint64 `gorm:"primaryKey" json:"id"`
	OldCommentID uint64	`json:"-"`
	AnswerRefer  int    `gorm:"index" json:"-"`
	Answer       Answer `json:"-" gorm:"foreignKey:AnswerRefer;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Content      string `json:"content"`
	AuthorID     uuid.UUID `gorm:"index;type:uuid;" json:"author_id"`
	VoteCount    int `json:"vote_count"`
	Created      int `gorm:"autoCreateTime" json:"created_time"`
	Deleted      gorm.DeletedAt `json:"-"`
}

//type ChildComment struct {
//	CommentID uint64 `gorm:"primaryKey"`
//	ParentID uint64
//	Content string
//	AuthorID uuid.UUID `gorm:"index"`
//	ParentAuthorID uuid.UUID
//	VoteCount int
//	Created int	`gorm:"autoCreateTime"`
//	Deleted gorm.DeletedAt
//}
