package entity

import (
	glog "gin_demo/interal/log"
	"gorm.io/gorm"
	"time"
)

var log = glog.Log
type Types map[string]interface{}

var Entities = Types{
	"topic": &Topic{},
	"question": &Question{},
	"answer": &Answer{},
	"comment": &Comment{},
	"user": &User{},
	"vote_answer": &VoteAnswer{},
	"vote_comment": &VoteComment{},
	"follow": &Follow{},
	"attention_question": &Attention{},
	"answer_notify": &AnswerNotify{},
	"vote_answer_notify": &VoteAnswerNotify{},
	"vote_comment_notify": &VoteCommentNotify{},
}

func (t Types) Migrate(db *gorm.DB)  {
	for _, entity := range t {
		if err := db.AutoMigrate(entity); err != nil {
			log.Debugf("entity: migrate %s (waiting 1s)", err.Error())
			time.Sleep(time.Second)
			if err := db.AutoMigrate(entity); err != nil {
				panic(err)
			}
		}
	}
}

//创建默认的数据
func CreateDefaultFixtures() {
	CreateDefaultUsers()
	//CreateDefaultMessageTypes()
}

func MigrateDB(db *gorm.DB)  {
	Entities.Migrate(db)
	CreateDefaultFixtures()
}
