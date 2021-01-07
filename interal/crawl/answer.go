package crawl

import (
	"fmt"
	"gin_demo/interal/entity"
	"github.com/satori/go.uuid"
	"io"
	"log"
	"net/http"
	"time"
)

type AnswerFetch struct {
	QuestionID       int
	QuestionAuthorID uuid.UUID
	Answers          []entity.Answer
	NextURL          string
	Users            []entity.User
	Ch               chan []CommentFetch
}

//Start 提取一个问题下的所有回答
func (af *AnswerFetch) Fetch() {
	//log.Printf("1111===> Start Fetch Answer For QuestionID %d\n", af.QuestionID)
	totalAnswers := 0
	for af.NextURL != "" && totalAnswers <= 80 {
		resp, err := http.Get(af.NextURL)
		if err != nil {
			log.Printf("requst %s error: %v", af.NextURL, err)
			break
		}
		if af.parse(resp.Body) != nil {
			//time.Sleep(time.Minute)
			//continue
			//panic("crawl refused")
			log.Printf("parse questions<%s> error: %v\ncontinue...", af.NextURL, err)
			continue
		}

		//save to db
		db := entity.Db()
		//create new users
		for _, user := range af.Users {
			CreateUserOrSkip(&user)
			UpdateAnswerCountOfUser(user.ID)
			AttentionQuestionOrSkip(user.ID, af.QuestionID)
			FollowUser(user.ID, af.QuestionAuthorID)
		}
		af.Users = nil
		//create new answers
		var commentFetches []CommentFetch
		for _, ans := range af.Answers {
			db.Create(&ans)
			commentFetch := CommentFetch{
				AnswerID:         ans.AnswerID,
				AnswerAuthorID:   ans.AuthorID,
				NextURL:          fmt.Sprintf(commentURL, ans.OldAnswerID),
				QuestionAuthorID: af.QuestionAuthorID,
			}
			commentFetches = append(commentFetches, commentFetch)
			//af.Ch<-commentFetch
		}
		totalAnswers += len(af.Answers)
		af.Ch <- commentFetches
		af.Answers = nil
		//request next 20 answers...
		time.Sleep(time.Second * 3)
	}
	//log.Printf("1111===> END Fetch Answer For QuestionID %d\n", af.QuestionID)
}

//Parse 解析answerURL响应
func (af *AnswerFetch) parse(body io.ReadCloser) error {
	js, err := streamToJson(body)
	if err != nil {
		return err
	}
	jsAnswerSlice := js.Get("data").MustArray()
	jsData := js.Get("data")
	for i := range jsAnswerSlice {
		//each answer
		jsIndex := jsData.GetIndex(i)
		author := jsIndex.Get("author")
		var userItem = entity.User{
			ID:       uuid.FromStringOrNil(author.Get("id").MustString()),
			Name:     author.Get("name").MustString(),
			HeadLine: author.Get("headline").MustString(),
		}
		af.Users = append(af.Users, userItem)
		var answerItem = entity.Answer{
			OldAnswerID:   jsIndex.Get("id").MustInt(),
			QuestionRefer: af.QuestionID,
			Excerpt:       jsIndex.Get("excerpt").MustString(),
			Content:       jsIndex.Get("content").MustString(),
			AuthorID:      userItem.ID,
			CommentCount:  jsIndex.Get("comment_count").MustInt(),
			VoteCount:     jsIndex.Get("voteup_count").MustInt(),
			Created:       jsIndex.Get("created_time").MustInt(),
			Updated:       jsIndex.Get("updated_time").MustInt(),
		}
		af.Answers = append(af.Answers, answerItem)
	}

	isEnd := js.Get("paging").Get("is_end").MustBool()
	if isEnd == false {
		af.NextURL = js.Get("paging").Get("next").MustString()
	} else {
		af.NextURL = ""
	}
	return nil
}
