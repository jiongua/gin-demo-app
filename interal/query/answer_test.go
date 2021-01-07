package query

import (
	"gin_demo/interal/entity"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestAnswersByQuestionID(t *testing.T) {
	type want struct {

	}
	var tests = []struct {
		qid int
		offset int
		order string
	}{
		{
			qid: 1200,
			offset: 0,
			order: "created",
		},
		{
			qid: 1200,
			offset: 0,
			order: "vote_count",
		},
	}
	for _, test := range tests {
		answers, _ := AnswersByQuestionID(test.qid, test.offset, test.order)
		if test.order == "created" {
			sortedByCreated := createdList(answers)
			sort.Ints(sortedByCreated)
			assert.Equal(t, sortedByCreated, createdList(answers))
		} else if test.order == "vote_count" {
			sortedByVote := voteList(answers)
			sort.Ints(sortedByVote)
			assert.Equal(t, sortedByVote, voteList(answers))
		}
	}
}

func createdList(result entity.Answers) (out []int) {
	for _, e := range result {
		out = append(out, e.Created)
	}
	return
}

func voteList(result entity.Answers) (out []int) {
	for _, e := range result {
		out = append(out, e.VoteCount)
	}
	return
}