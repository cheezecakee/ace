// Package pack
package pack

import (
	"fmt"
	"time"

	"github.com/cheezecakee/ace/internal/engine"
)

type (
	Category   string
	Categories []Category
)

type Role map[string]Categories

type PackInfo struct {
	ID            string
	Name          string
	Creator       string
	Version       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Role          string
	Categories    Categories
	FilePath      string
	QuestionCount int
}

type Pack struct {
	ID         string
	Name       string
	Creator    string
	Version    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Role       string
	FilePath   string
	Categories map[Category]CategoryQuestions
}

type Validator interface {
	Validate(seenIDs map[string]bool) error
	Repair(category string, index int)
}

type CategoryQuestions struct {
	Choice         []ChoiceQuestionData
	MultipleChoice []MultipleChoiceQuestionData
	TrueFalse      []TrueFalseQuestionData
	TextEntry      []TextEntryQuestionData
}

type ChoiceQuestionData struct {
	ID         string
	Difficulty string
	Prompt     string
	Options    []string
	Answer     int
}

func (q *ChoiceQuestionData) Validate(seenIDs map[string]bool) error {
	if q.ID == "" || q.Prompt == "" {
		return ErrMissingFields
	}

	if seenIDs[q.ID] {
		return ErrDuplicateID
	}
	seenIDs[q.ID] = true

	if !isValidDifficulty(q.Difficulty) {
		return ErrInvalidData
	}

	if q.Answer < 0 || q.Answer >= len(q.Options) {
		return ErrInvalidData
	}

	return nil
}

func (q *ChoiceQuestionData) Repair(category string, index int) {
	if q.ID == "" {
		q.ID = fmt.Sprintf("%s-choice-%d", category, index)
	}
}

type MultipleChoiceQuestionData struct {
	ID         string
	Difficulty string
	Prompt     string
	Options    []string
	Answer     []int
}

func (q *MultipleChoiceQuestionData) Validate(seenIDs map[string]bool) error {
	if q.ID == "" || q.Prompt == "" {
		return ErrMissingFields
	}

	if seenIDs[q.ID] {
		return ErrDuplicateID
	}
	seenIDs[q.ID] = true

	if !isValidDifficulty(q.Difficulty) {
		return ErrInvalidData
	}

	if len(q.Answer) < 1 || len(q.Answer) > len(q.Options) || len(q.Options) < 1 {
		return ErrInvalidData
	}

	return nil
}

func (q *MultipleChoiceQuestionData) Repair(category string, index int) {
	if q.ID == "" {
		q.ID = fmt.Sprintf("%s-multiple-%d", category, index)
	}
}

type TrueFalseQuestionData struct {
	ID         string
	Difficulty string
	Prompt     string
	Answer     bool
}

func (q *TrueFalseQuestionData) Validate(seenIDs map[string]bool) error {
	if q.ID == "" || q.Prompt == "" {
		return ErrMissingFields
	}

	if seenIDs[q.ID] {
		return ErrDuplicateID
	}
	seenIDs[q.ID] = true

	if !isValidDifficulty(q.Difficulty) {
		return ErrInvalidData
	}

	return nil
}

func (q *TrueFalseQuestionData) Repair(category string, index int) {
	if q.ID == "" {
		q.ID = fmt.Sprintf("%s-tf-%d", category, index)
	}
}

type TextEntryQuestionData struct {
	ID             string
	Difficulty     string
	Prompt         string
	Keywords       []string
	ExpectedAnswer string
}

func (q *TextEntryQuestionData) Validate(seenIDs map[string]bool) error {
	if q.ID == "" || q.Prompt == "" {
		return ErrMissingFields
	}

	if seenIDs[q.ID] {
		return ErrDuplicateID
	}
	seenIDs[q.ID] = true

	if !isValidDifficulty(q.Difficulty) {
		return ErrInvalidData
	}

	if len(q.ExpectedAnswer) < 3 || len(q.Keywords) < 1 {
		return ErrInvalidData
	}

	return nil
}

func (q *TextEntryQuestionData) Repair(category string, index int) {
	if q.ID == "" {
		q.ID = fmt.Sprintf("%s-text-%d", category, index)
	}
}

// -- Helper --//
func isValidDifficulty(difficulty string) bool {
	return engine.ParseDifficulty(difficulty) != 0
}
