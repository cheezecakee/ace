package engine

type Answer interface {
	Type() QuestionType
	Value() any // The actual answer data
}

type ChoiceAnswer struct {
	Selected int
}

func (a ChoiceAnswer) Type() QuestionType { return Choice }
func (a ChoiceAnswer) Value() any         { return a.Selected }

type MultipleChoiceAnswer struct {
	Selected []int
}

func (a MultipleChoiceAnswer) Type() QuestionType { return MultipleChoice }
func (a MultipleChoiceAnswer) Value() any         { return a.Selected }

type TrueFalseAnswer struct {
	Answer bool
}

func (a TrueFalseAnswer) Type() QuestionType { return TrueFalse }
func (a TrueFalseAnswer) Value() any         { return a.Answer }

type TextEntryAnswer struct {
	Text string
}

func (a TextEntryAnswer) Type() QuestionType { return TextEntry }
func (a TextEntryAnswer) Value() any         { return a.Text }
