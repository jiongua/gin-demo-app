//question包获取指定话题下的精华问题集
package crawl

import (
	"fmt"
	"gin_demo/interal/entity"
	"io"
	"log"
	"net/http"
	"time"
)

//QuestionFetch 实现了Fetcher接口
type QuestionFetch struct {
	TopicID   int
	TopicName string
	Questions []entity.Question
	NextURL   string
	AnswerCh  chan []AnswerFetch
	CommentCh chan []CommentFetch
}

//Fetch 根据主题ID获取1000个问题列表
func (qf *QuestionFetch) Fetch() {
	log.Printf("0000===> Start Fetch questions for topic %s\n", qf.TopicName)
	db := entity.Db()
	topicItem := entity.Topic{
		TopicID: qf.TopicID,
		Name:    qf.TopicName,
	}
	db.Create(&topicItem)
	for qf.NextURL != "" {
		//请求topic_feed URL, 获取问题列表
		resp, err := http.Get(qf.NextURL)
		if err != nil {
			log.Printf("get Topic URL<%s> error: %v", qf.NextURL, err)
			break
		}
		if qf.parse(resp.Body) != nil {
			log.Printf("parse topic<%s> error: %v\ncontinue...", qf.NextURL, err)
			continue
			//panic("crawl refused")
		}
		//- 新问题写库
		//- 发布人关注问题
		//- 将问题写入chan
		db := entity.Db()
		var answerFetches []AnswerFetch
		for _, ques := range qf.Questions {
			db.Create(&ques)
			AttentionQuestionOrSkip(ques.AuthorID, ques.QuestionID)
			answerFetcher := AnswerFetch{
				QuestionID:       ques.QuestionID,
				QuestionAuthorID: ques.AuthorID,
				NextURL:          fmt.Sprintf(answerURLHead, ques.OldQuestionID) + answerURLEmbed + answerURLStartPage,
				Ch:               qf.CommentCh,
			}
			answerFetches = append(answerFetches, answerFetcher)
		}
		qf.AnswerCh <- answerFetches
		qf.Questions = nil
		//request next 10 questions...
		log.Printf("request next 10 questions: <%s>\n", qf.NextURL)
		time.Sleep(time.Second * 3)
	}
	log.Printf("0000===> END Fetch questions for topic %s, lastURL = %s\n", qf.TopicName, qf.NextURL)
}

//Parse 解析TopicURL响应, 每次得到10个待处理问题
func (qf *QuestionFetch) parse(respBody io.ReadCloser) error {
	js, err := streamToJson(respBody)
	if err != nil {
		return err
	}
	data := js.Get("data").MustArray()
	//捕获所有问题
	for i := range data {
		questionBody := js.Get("data").GetIndex(i).Get("target").Get("question")
		/*
			"question":{
					"title":"如何评价电影版《滚蛋吧！肿瘤君》？",
					"url":"http://www.zhihu.com/api/v4/questions/32373308",
					"created":1437290985,
					"is_following":false,
					"question_type":"normal",
					"type":"question",
					"id":32373308
				},
		*/
		//skip non question
		OldQuestionID := questionBody.Get("id").MustInt()
		if OldQuestionID == 0 {
			continue
		}
		questionItem := entity.Question{
			OldQuestionID: questionBody.Get("id").MustInt(),
			TopicID:       qf.TopicID,
			Title:         questionBody.Get("title").MustString(),
			Created:       questionBody.Get("created").MustInt(),
			Updated:       questionBody.Get("created").MustInt(), //update in answerAPI
			AuthorID:      selectRandomUserID(),
		}
		qf.Questions = append(qf.Questions, questionItem)
	}

	isEnd := js.Get("paging").Get("is_end").MustBool()
	if isEnd == false {
		qf.NextURL = js.Get("paging").Get("next").MustString()
	} else {
		log.Printf("<=== END crawl topic[%d], cause isEnd(%t)\n", qf.TopicID, isEnd)
		qf.NextURL = ""
	}
	return nil
}
