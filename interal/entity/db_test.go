package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigPath(t *testing.T)  {
	assert.Equal(t, "/Users/jiongua/zhihu/gin-demo-app/config/pgsql.json", ConfigPath())
}

func TestDb(t *testing.T) {
	if Db() == nil {
		t.Error("connect db error")
	}
}
