package crawl

import (
	"gin_demo/interal/entity"
	"github.com/satori/go.uuid"
	"io"
	"log"
	"net/http"
	"time"
)

type CommentFetch struct {
	AnswerID         int
	AnswerAuthorID   uuid.UUID
	QuestionAuthorID uuid.UUID
	Comments         []entity.Comment
	NextURL          string
	Users            []entity.User
}

func (cf *CommentFetch) Fetch() {
	//log.Printf("2222 ===> Start Fetch Comments For AnswerID %d\n", cf.AnswerID)
	//只取20条评论
	if cf.NextURL != "" {
		resp, err := http.Get(cf.NextURL)
		if err != nil {
			log.Printf("fetch comment url<%s> error: %v\n", cf.NextURL, err)
			return
		}
		if cf.parse(resp.Body) != nil {
			log.Printf("parse answer<%s> error: %v\n", cf.NextURL, err)
			return
		}
		//- 新用户写库
		//- 新用户关注回答者和问题发布者
		//- 评论写库
		//- 对匿名账号不作处理
		for _, user := range cf.Users {
			CreateUserOrSkip(&user)
			FollowUser(user.ID, cf.AnswerAuthorID)
			FollowUser(user.ID, cf.QuestionAuthorID)
		}
		for _, comm := range cf.Comments {
			CreateCommentOrSkip(&comm)
		}
		cf.Users = nil
		cf.Comments = nil
		//request next 20 comments...
		time.Sleep(time.Second * 3)
	}
	//log.Printf("2222 ===> END Fetch Comments For AnswerID %d\n", cf.AnswerID)
}

func (cf *CommentFetch) parse(body io.ReadCloser) error {
	js, err := streamToJson(body)
	if err != nil {
		//log.Printf("stream to json err: %v\n", err)
		return err
	}
	paging := js.Get("paging")
	dataBody := js.Get("data").MustArray()
	for i := range dataBody {
		commentBody := js.Get("data").GetIndex(i)
		//create user
		memberBody := commentBody.Get("author").Get("member")
		//id为"0"是匿名账号
		userOrigId := memberBody.Get("id").MustString()
		var userItem = entity.User{
			ID:       uuid.FromStringOrNil(userOrigId),
			Name:     memberBody.Get("name").MustString(),
			HeadLine: memberBody.Get("headline").MustString(),
		}
		cf.Users = append(cf.Users, userItem)
		var commentItem = entity.Comment{
			OldCommentID: commentBody.Get("id").MustUint64(),
			AnswerRefer:  cf.AnswerID,
			Content:      commentBody.Get("content").MustString(),
			//authorID为NIL UUID代表匿名评论
			AuthorID:  userItem.ID,
			VoteCount: commentBody.Get("vote_count").MustInt(),
			Created:   commentBody.Get("created_time").MustInt(),
		}
		cf.Comments = append(cf.Comments, commentItem)
	}
	isEnd := paging.Get("is_end").MustBool()
	if isEnd == false {
		cf.NextURL = paging.Get("next").MustString()
	} else {
		cf.NextURL = ""
	}
	return nil
}

/* not used
func ChildComments(c *colly.Collector, ParentCommentID uint64, ParentAuthorID uuid.UUID) {
	c.OnResponse(func(r *colly.Response) {
		js, err := simplejson.NewJson(r.Body)
		if err != nil {
			panic("invalid json")
		}
		paging := js.Get("paging")
		dataBody, err := js.Get("data").Array()
		if err != nil {
			panic("data is not array")
		}
		db := entity.GetDB()
		for i := range dataBody {
			commentBody := js.Get("data").GetIndex(i)
			//create user
			memberBody := commentBody.Get("author").Get("member")
			userOrigId := memberBody.Get("id").MustString()
			var userItem = entity.User{
				ID:       uuid.FromStringOrNil(userOrigId),
				Name:     memberBody.Get("name").MustString(),
				HeadLine: memberBody.Get("headline").MustString(),
			}
			//db.Create(&userItem)
			CreateUserOrSkip(&userItem)
			var commentItem = entity.ChildComment{
				CommentID:      commentBody.Get("id").MustUint64(),
				ParentID:       ParentCommentID,
				Content:        commentBody.Get("content").MustString(),
				AuthorID:       userItem.ID,
				ParentAuthorID: ParentAuthorID,
				VoteCount:      commentBody.Get("vote_count").MustInt(),
				Created:        commentBody.Get("created_time").MustInt(),
			}
			db.Create(&commentItem)
		}
		isEnd, err := paging.Get("is_end").Bool()
		if err != nil {
			log.Printf("err in parse paging by simplejson: %s\n", err.Error())
			return
		} else if isEnd == false {
			nextURL, _ := paging.Get("next").String()
			c.Visit(nextURL)
		}

	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visting %d_Child_Comment_for_parent: %s\n",
			ParentCommentID,
			r.URL.String())
	})

	ChildCommentURL := fmt.Sprintf(`https://www.zhihu.com/api/v4/comments/%d/child_comments?limit=20&offset=0`, ParentCommentID)
	c.Visit(ChildCommentURL)
}
*/
