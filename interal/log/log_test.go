package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileName(t *testing.T) {
	assert.Equal(t, "/Users/jiongua/zhihu/gin-demo-app/logs/api.log", FileName())
}
