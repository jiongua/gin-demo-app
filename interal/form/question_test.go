package form
//
//import (
//	"testing"
//)
//
//func TestQuestionForm_InvalidOrder(t *testing.T) {
//	driver := []struct {
//		input QuestionForm
//		want string
//	}{
//		{
//			input: QuestionForm{
//				AnswerOrder:  "xxxx",
//			},
//			want: "created",
//		},
//		{
//			input: QuestionForm{
//				AnswerOrder:  "byHot",
//			},
//			want: "vote_count",
//		},
//		{
//			input: QuestionForm{
//				AnswerOrder:  "byTime",
//			},
//			want: "created",
//		},
//
//	}
//	for _, test := range driver {
//		if got := test.input.OrderFormat(); got != test.want {
//			t.Errorf("got<%s>, but want<%s>", got, test.want)
//		}
//	}
//}
