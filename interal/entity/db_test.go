package entity

import (
	"testing"
)


func TestDb(t *testing.T) {
	if Db() == nil {
		t.Error("connect db error")
	}
}
