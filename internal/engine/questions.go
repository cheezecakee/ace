package engine

type Question interface {
	Type() QuestionType
	GetPrompt() string
	GetAnswer() Answer
}

type ChoiceQuestion struct {
	Prompt  string
	Options []string
	Correct int
}

func (q ChoiceQuestion) Type() QuestionType { return Choice }
func (q ChoiceQuestion) GetPrompt() string  { return q.Prompt }
func (q ChoiceQuestion) GetAnswer() Answer {
	return ChoiceAnswer{Selected: q.Correct}
}

type MultipleChoiceQuestion struct {
	Prompt  string
	Options []string
	Correct []int
}

func (q MultipleChoiceQuestion) Type() QuestionType { return MultipleChoice }
func (q MultipleChoiceQuestion) GetPrompt() string  { return q.Prompt }
func (q MultipleChoiceQuestion) GetAnswer() Answer {
	return MultipleChoiceAnswer{Selected: q.Correct}
}

type TrueFalseQuestion struct {
	Prompt  string
	Correct bool
}

func (q TrueFalseQuestion) Type() QuestionType { return TrueFalse }
func (q TrueFalseQuestion) GetPrompt() string  { return q.Prompt }
func (q TrueFalseQuestion) GetAnswer() Answer {
	return TrueFalseAnswer{Answer: q.Correct}
}

type TextEntryQuestion struct {
	Prompt         string
	ExpectedAnswer string   // Idead answer for AI comparison
	Keywords       []string // Key concepts that should be present
}

func (q TextEntryQuestion) Type() QuestionType { return TextEntry }
func (q TextEntryQuestion) GetPrompt() string  { return q.Prompt }
func (q TextEntryQuestion) GetAnswer() Answer {
	return TextEntryAnswer{Text: q.ExpectedAnswer}
}
