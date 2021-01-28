package value

type ActionID int

const (
	CREATE_QUESTION = iota
	VIEW_QUESTION
	ANSWER
	COMMENT
	VOTE_ANSWER
	VOTE_COMMENT
	)
func (a ActionID) GetByName()  {

}
