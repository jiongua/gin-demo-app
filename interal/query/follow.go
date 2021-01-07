package query

import (
	"gin_demo/interal/entity"
	uuid "github.com/satori/go.uuid"
)

func IsFollow(me uuid.UUID, target uuid.UUID) bool {
	if entity.Db().Model(&entity.Follow{}).Where("follower_id = ? and followee_id = ?", me, target).RowsAffected == 0 {
		return false
	}else {
		return true
	}
}
