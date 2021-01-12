package interal

import (
	"github.com/bmizerany/assert"
	"os"
	"testing"
)

func TestGetEnvWithDefault(t *testing.T) {
	os.Setenv("key1", "v1")
	assert.Equal(t, "v1", GetEnvWithDefault("key1", "v2"))
	os.Unsetenv("key1")
	assert.Equal(t, "v2", GetEnvWithDefault("key1", "v2"))
}
