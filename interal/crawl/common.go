package crawl

import (
	"bytes"
	"gin_demo/interal/entity"
	"github.com/bitly/go-simplejson"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
)

type Fetcher interface {
	Fetch()
}

func streamToJson(stream io.ReadCloser) (*simplejson.Json, error) {
	defer stream.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	js, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		log.Printf("parse json err:%v\n", err)
		return nil, err
	}
	if js.Get("error").Get("need_login").MustBool() == true {
		//log.Println("")
		//return nil, errors.New("fetch refuse")
		panic("fetch refuse by zhihu...")
	}
	return js, nil
}

//IsAnonymousUser检查指定用户ID是否是匿名
func IsAnonymousUser(id uuid.UUID) bool {
	if id == uuid.Nil {
		//log.Printf("id <%v> is anonymous", id)
		return true
	}
	return false
}

//AttentionQuestionOrSkip 关注问题
func AttentionQuestionOrSkip(userID uuid.UUID, questionID int) {
	if IsAnonymousUser(userID) {
		return
	}
	db := entity.Db()
	//检查用户是否已关注该问题
	//select count(1) from attentions where user_id = userID and question_refer = questionID
	//var count int64
	//db.Model(&entity.Attention{}).Where("user_id = ? AND question_refer = ?", userID, questionID).Count(&count)
	//if count == 0 {
	//	db.Create(&entity.Attention{
	//		UserID:        userID,
	//		QuestionRefer: questionID,
	//	})
	//}
	//return
	newAttentionItem := entity.Attention{
		UserID:        userID,
		QuestionRefer: questionID,
	}
	db.Where(entity.Attention{
		UserID:        userID,
		QuestionRefer: questionID,
	}).FirstOrCreate(&newAttentionItem)
}

//FollowUser 关注用户
func FollowUser(follower, followee uuid.UUID) {
	//处理匿名账号
	if IsAnonymousUser(follower) || IsAnonymousUser(followee) {
		return
	}
	if uuid.Equal(follower, followee) {
		return
	}
	//先检查follower是否关注了followee
	//select count(1) from follows where follower_id = follower AND followee = followee_id
	//var count int64
	//entity.GetDB().Transaction(func(tx *gorm.DB) error {
	//	tx.Model(&entity.Follow{}).Where("follower_id = ? AND followee_id = ?", follower, followee).Count(&count)
	//	if count == 0 {
	//		if err := tx.Create(&entity.Follow{
	//			FollowerID:  	follower,
	//			FolloweeID: 	followee,
	//		}).Error; err != nil {
	//			return err
	//		}
	//	}
	//	return nil
	//})
	FollowItem := entity.Follow{
		FollowerID: follower,
		FolloweeID: followee,
	}
	entity.Db().Where(entity.Follow{
		FollowerID: follower,
		FolloweeID: followee,
	}).FirstOrCreate(&FollowItem)
	//db.Model(&entity.Follow{}).Where("follower_id = ? AND followee_id = ?", follower, followee).Count(&count)
	//if count == 0 {
	//	db.Create(&entity.Follow{
	//		FollowerID:  	follower,
	//		FolloweeID: 	followee,
	//	})
	//}
	//return
}

//CreateUserOrSkip 创建用户
func CreateUserOrSkip(user *entity.User) {
	//事务
	//var count int64
	//entity.GetDB().Transaction(func(tx *gorm.DB) error {
	//	tx.Model(&entity.User{}).Where("id = ?", user.ID).Count(&count)
	//	if count == 0 {
	//		if err := tx.Create(user).Error; err != nil {
	//			return err
	//		}
	//	}
	//	return nil
	//})
	if IsAnonymousUser(user.ID) {
		//匿名用户，无序创建新账号
		//log.Printf("don't create userID<%v-%s-%s>\n", user.ID, user.Name, user.HeadLine)
		return
	}
	entity.Db().Where(entity.User{
		ID: user.ID,
	}).FirstOrCreate(user)
	return
}

//创建评论
func CreateCommentOrSkip(comment *entity.Comment) {
	//var count int64
	//entity.GetDB().Transaction(func(tx *gorm.DB) error {
	//	tx.Model(&entity.Comment{}).Where("comment_id = ?", comment.CommentID).Count(&count)
	//	if count == 0 {
	//		if err := tx.Create(comment).Error; err != nil {
	//			return err
	//		}
	//	}
	//	return nil
	//})
	//if count > 0 {
	//	return false
	//}
	//return true
	entity.Db().Where(entity.Comment{
		OldCommentID: comment.OldCommentID,
	}).FirstOrCreate(comment)
}

//UpdateAnswerCountOfUser 更新用户评论数
func UpdateAnswerCountOfUser(userID uuid.UUID) {
	if IsAnonymousUser(userID) {
		return
	}
	//原子更新
	db := entity.Db()
	db.Exec("UPDATE users SET answer_count = answer_count + 1 where id = ?", userID)
}

//随机选择一个用户ID作为问题发布人
func selectRandomUserID() uuid.UUID {
	db := entity.Db()
	var userItem entity.User
	result := db.Raw("select * from users order by random() limit 1").Scan(&userItem)
	if result.RowsAffected == 0 {
		userItem = entity.User{
			ID:       uuid.NewV4(),
			HeadLine: "测试",
			Name:     "JIONGUA",
		}
		db.Create(&userItem)
	}
	return userItem.ID
}
