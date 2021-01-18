package query

import (
	"gin_demo/interal/entity"
	uuid "github.com/satori/go.uuid"
	"time"
)

//QuestionByID 获取指定问题信息, 包含
// - 问题所属的话题、问题标题、问题关注人数、问题浏览人数
//		- 问题答案列表
//			- 回答人信息：用户名、headline、是否关注、回答数、粉丝数
//			- 回答信息：回答内容、回答点赞数、评论数
/*
响应json:
	{
		"question": {
			"topic_name":
			"title":
			"attention_count":
			"view_count":
		},
		"answers": [
			{
				"id"
				"author": {},
				"data": {
					"content":
				}
			},
			{}
		]
	}
*/
type Answer struct {
	AnswerEntity entity.Answer	`json:"answer"`
	Author Author		`json:"author"`
	HumanTime time.Time `json:"created_time2"`
}
// Result is url /question/:qid response
type QuestionAllResult struct {
	TopicName string	`json:"topic_name"`
	Title string		`json:"title"`
	AttentionCnt int64	`json:"attention_count"`
	ViewCnt int64		`json:"view_count"`
	AnswersResult []Answer 	`json:"data"`
}

//QuestionByID 根据问题ID获取该问题的所有信息
func QuestionByID(qid int, offset int, _ int, order string) (QuestionAllResult, error){
	var finalResult QuestionAllResult
	var questionHeader struct{
		Name string
		Title string
	}
	if err := entity.Db().Model(&entity.Question{}).Select("name, title").Where("question_id = ?", qid).
		Joins("join topics on topics.topic_id = questions.topic_id").Find(&questionHeader).Error; err != nil {
			return finalResult, err
	}
	finalResult.TopicName = questionHeader.Name
	finalResult.Title = questionHeader.Title
	finalResult.AttentionCnt = GetAttentionCount(qid)
	finalResult.ViewCnt = GetViewCount(qid)
	//answers和相关的作者信息
	if answers, err := AnswersByQuestionID(qid, offset, order); err != nil {
		return finalResult, err
	} else {
		//组合回答和回答者信息
		for _, item := range answers {
			author, err := GetAuthorWithFollow(item.AuthorID, entity.Admin.ID)
			if err != nil {
				return finalResult, err
			}
			finalResult.AnswersResult = append(finalResult.AnswersResult, Answer{
				AnswerEntity: item,
				Author: author,
				HumanTime: time.Unix(int64(item.Created), 0),
			})
		}
		return finalResult, nil
	}
}

// GetAttentionCount 根据问题ID获取问题的关注人数
func GetAttentionCount(qid int) int64 {
	var count int64
	entity.Db().Table("attentions").Where("question_refer = ?", qid).Count(&count)
	return count
}

func GetViewCount(qid int) int64 {
	return 1024
}

func GetQuestionsByAuthorID(authorID uuid.UUID, offset int, isMine bool) ([]entity.Question, error) {
	result := make([]entity.Question, 0, 10)
	q := entity.Db().Limit(10).Order("created desc").Offset(offset).
		Where("author_id=?", authorID)
	if !isMine {
		//不是查看自己的主页，过滤掉匿名的提问
		q = q.Where("is_anonymity=?", false)
	}
	if err := q.Find(&result).Error; err!=nil {
		return nil, err
	}
	return result, nil
}




