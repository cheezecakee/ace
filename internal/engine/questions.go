package engine

type Questions []Question

type Question interface {
	Type() QuestionType
	GetPrompt() string
	GetAnswer() Answer
}

type BaseQuestion struct {
	ID     string
	Prompt string
}

func (q BaseQuestion) GetPrompt() string { return q.Prompt }

func (q BaseQuestion) GetID() string { return q.ID }

type ChoiceQuestion struct {
	BaseQuestion
	Options []string
	Correct int
}

func (q ChoiceQuestion) Type() QuestionType { return Choice }

func (q ChoiceQuestion) GetAnswer() Answer {
	return ChoiceAnswer{Selected: q.Correct}
}

type MultipleChoiceQuestion struct {
	BaseQuestion
	Options []string
	Correct []int
}

func (q MultipleChoiceQuestion) Type() QuestionType { return MultipleChoice }

func (q MultipleChoiceQuestion) GetAnswer() Answer {
	return MultipleChoiceAnswer{Selected: q.Correct}
}

type BoolQuestion struct {
	BaseQuestion
	Correct bool
}

func (q BoolQuestion) Type() QuestionType { return Bool }

func (q BoolQuestion) GetAnswer() Answer {
	return BoolAnswer{Answer: q.Correct}
}

type TextEntryQuestion struct {
	BaseQuestion
	ExpectedAnswer string   // Idead answer for AI comparison
	Keywords       []string // Key concepts that should be present
}

func (q TextEntryQuestion) Type() QuestionType { return TextEntry }

func (q TextEntryQuestion) GetAnswer() Answer {
	return TextEntryAnswer{Text: q.ExpectedAnswer}
}
