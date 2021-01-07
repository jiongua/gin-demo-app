package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDefaultUsers(t *testing.T) {
	assert.Equal(t, true, Admin.IsAdmin())
	FirstOrCreateUser(&Admin)
}
