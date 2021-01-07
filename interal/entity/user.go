package entity

import (
	"github.com/satori/go.uuid"
)

type User struct {
	ID            uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	HeadLine      string `json:"head_line"`
	Name          string `json:"name"`
	AnswerCount   int `json:"answer_count"`
	FolloweeCount int `json:"followee_count"`
	FollowerCount int	`json:"follower_count"`
	VoteCount     int `json:"vote_count"`
}

var Admin = User{
	ID:            uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"),
	HeadLine:      "管理员",
	Name:          "Admin",
}


var Super1 = User{
	ID:            uuid.FromStringOrNil("00000000-0000-0000-0000-000000000002"),
	HeadLine:      "大V1",
	Name:          "Super1",
}

var Super2 = User{
	ID:            uuid.FromStringOrNil("00000000-0000-0000-0000-000000000003"),
	HeadLine:      "大V2",
	Name:          "Super2",
}

var Guest = User{
	ID:            uuid.FromStringOrNil("00000000-0000-0000-0000-000000000004"),
	HeadLine:      "zhihu-web-go测试账号",
	Name:          "zhihu-web-go测试账号",
}

var Anonymity = User{
	ID:             uuid.FromStringOrNil("00000000-0000-0000-0000-000000000005"),
	HeadLine:      "zhihu-web-go匿名账号",
	Name:          "zhihu-web-go匿名账号",
}

// CreateDefaultUsers initializes the database with default user accounts.
func CreateDefaultUsers() {
	if user := FirstOrCreateUser(&Admin); user != nil {
		Admin = *user
	}
	if user := FirstOrCreateUser(&Super1); user != nil {
		Super1 = *user
	}

	if user := FirstOrCreateUser(&Super2); user != nil {
		Super2 = *user
	}
	if user := FirstOrCreateUser(&Anonymity); user != nil {
		Anonymity = *user
	}
	if user := FirstOrCreateUser(&Guest); user != nil {
		Guest = *user
	}
}

//插入数据库
func (m *User) Create() error {
	//如果插入成功，则更新了m
	return Db().Create(m).Error
}

//创建用户，如果存在直接返回
func FirstOrCreateUser(m *User) *User {
	result := User{}

	if err := Db().Where("id = ?", m.ID).First(&result).Error; err == nil {
		return &result
	} else if err := m.Create(); err != nil {
		log.Debugf("create user err: %s", err)
		return nil
	}
	//success，此时m是插入数据库后gorm赋值的
	return m
}

func (m *User) IsAdmin() bool{
	return m.Name == "Admin" || m.ID == uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
}

func (m *User) IsAnonymity() bool {
	return m.Name == "zhihu-web-go匿名账号" || m.ID == uuid.Nil
}

