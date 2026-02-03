package pack

import (
	"fmt"

	"github.com/cheezecakee/ace/internal/engine"
)

type (
	Type  int
	Types []Type
)

const (
	TypeChoice Type = iota
	TypeMulti
	TypeBool
	TypeText
)

func (t Type) String() string {
	switch t {
	case TypeChoice:
		return "choice"
	case TypeMulti:
		return "multi"
	case TypeBool:
		return "bool"
	case TypeText:
		return "text"
	default:
		return ""
	}
}

// FromEngineType converts an engine.QuestionType to pack.Type
func FromEngineType(et engine.QuestionType) Type {
	switch et {
	case engine.Choice:
		return TypeChoice
	case engine.MultipleChoice:
		return TypeMulti
	case engine.Bool:
		return TypeBool
	case engine.TextEntry:
		return TypeText
	default:
		return TypeChoice // or panic/error
	}
}

// FromEngineTypes converts multiple engine types to pack types
func FromEngineTypes(engineTypes []engine.QuestionType) []Type {
	packTypes := make([]Type, 0, len(engineTypes))
	for _, et := range engineTypes {
		packTypes = append(packTypes, FromEngineType(et))
	}
	return packTypes
}

type Questions []Question

type Question struct {
	ID         string
	Difficulty engine.Difficulty
	Category   string
	Type       Type
	Prompt     string
	Answer     Answer
}

func (q Question) Validate() error {
	if q.ID == "" {
		return fmt.Errorf("question missing ID")
	}
	if q.Prompt == "" {
		return fmt.Errorf("question %s missing prompt", q.ID)
	}
	if !q.Difficulty.IsValid() {
		return fmt.Errorf("question %s has invalid difficulty", q.ID)
	}
	if q.Answer == nil {
		return fmt.Errorf("question %s missing answer", q.ID)
	}
	return nil
}

func (q Question) ToEngine() engine.Question {
	base := engine.BaseQuestion{
		ID:     q.ID,
		Prompt: q.Prompt,
	}

	switch q.Type {

	case TypeChoice:
		ans := q.Answer.(ChoiceAnswer)
		return engine.ChoiceQuestion{
			BaseQuestion: base,
			Options:      ans.Options,
			Correct:      ans.Correct,
		}

	case TypeMulti:
		ans := q.Answer.(MultiAnswer)
		return engine.MultipleChoiceQuestion{
			BaseQuestion: base,
			Options:      ans.Options,
			Correct:      ans.Correct,
		}

	case TypeBool:
		ans := q.Answer.(BoolAnswer)
		return engine.BoolQuestion{
			BaseQuestion: base,
			Correct:      ans.Correct,
		}

	case TypeText:
		ans := q.Answer.(TextAnswer)
		return engine.TextEntryQuestion{
			BaseQuestion:   base,
			ExpectedAnswer: ans.Expected,
			Keywords:       ans.Keywords,
		}

	default:
		return nil
	}
}

func (qs Questions) ToEngine() []engine.Question {
	result := make([]engine.Question, len(qs))
	for i, q := range qs {
		result[i] = q.ToEngine()
	}
	return result
}

func (qs Questions) Filter(types []Type, difficulty engine.Difficulty) Questions {
	var filtered Questions

	typeMap := make(map[Type]bool)
	for _, t := range types {
		typeMap[t] = true
	}

	for _, q := range qs {
		typeMatch := len(types) == 0 || typeMap[q.Type]
		difficultyMatch := difficulty == 0 || q.Difficulty == difficulty

		if typeMatch && difficultyMatch {
			filtered = append(filtered, q)
		}
	}

	return filtered
}

type Answer interface {
	isAnswer()
}

type ChoiceAnswer struct {
	Options []string
	Correct int
}

func (ChoiceAnswer) isAnswer() {}

type MultiAnswer struct {
	Options []string
	Correct []int
}

func (MultiAnswer) isAnswer() {}

type BoolAnswer struct {
	Correct bool
}

func (BoolAnswer) isAnswer() {}

type TextAnswer struct {
	Expected string
	Keywords []string
}

func (TextAnswer) isAnswer() {}
