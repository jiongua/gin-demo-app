package entity

//回答被点赞
type VoteAnswerNotify struct {
	NotifyMeta
	AnswerID int
	AnswerURL string	`gorm:"varchar(200)"`
}

//评论被点赞
type VoteCommentNotify struct {
	NotifyMeta
	CommentID     int
	CommentURL    string	`gorm:"varchar(200)"`
}
