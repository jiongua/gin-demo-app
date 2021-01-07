package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuestionByID(t *testing.T) {
	type want struct {
		title string
		topicName string
		Attentions int64
	}
	var tests = []struct {
		questionID int
		offset int
		order string
		want
	}{{
		questionID: 1200,
		offset: 10,
		order: "created",
		want: want{
			title: "如何评价《中国2098》科幻系列作品水平?",
			topicName: "科技",
			Attentions: 96,
			},
		},
	{
		questionID: 1200,
		offset: 10,
		order: "vote_count",
		want: want{
			title: "如何评价《中国2098》科幻系列作品水平?",
			topicName: "科技",
			Attentions: 96,
		},
	},
	}
	log.Info("start test questionByID")
	for _, test := range tests {
		if got, err := QuestionByID(test.questionID, test.offset, 10, test.order); err == nil {
			assert.Equal(t, test.want.topicName, got.TopicName)
			assert.Equal(t, test.want.title, got.Title)
			assert.Equal(t, test.want.Attentions, got.AttentionCnt)
			t.Log(got.AnswersResult)
		} else {
			t.Logf("QuestionByID return error: %v", err)
			assert.EqualError(t, err, "record not found")
		}
	}
}