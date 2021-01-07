package entity

import "github.com/satori/go.uuid"

type Follow struct {
	FollowerID uuid.UUID `gorm:"primaryKey;type:uuid;"`
	FolloweeID uuid.UUID `gorm:"primaryKey;type:uuid;"`
}

