package query

import (
	"gin_demo/interal/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Author struct {
	UserEntity entity.User 	`json:"author"`
	IsFollow bool	`json:"is_follow"`
}

func GetUserByID(uid uuid.UUID) (entity.User, error) {
	if IsAnonymousUser(uid) {
		//return 匿名账户
		return entity.Anonymity, nil
	}
	var result entity.User
	if err := entity.Db().First(&result, uid).Error; err != nil {
		log.WithFields(logrus.Fields{
			"userID": uid,
		}).Error("query user error")
		return result, err
	}
	return result, nil
}

func GetAuthorWithFollow(authorID uuid.UUID, loginID uuid.UUID) (author Author, err error)  {
	user, err := GetUserByID(authorID)
	if err != nil {
		return Author{}, err
	}
	return Author{
		UserEntity: user,
		IsFollow:   IsFollow(loginID, authorID),
	}, nil
}

//IsAnonymousUser检查指定用户ID是否是匿名
func IsAnonymousUser(id uuid.UUID) bool {
	return id == uuid.Nil
}

