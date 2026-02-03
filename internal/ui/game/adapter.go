// Package game
package game

import (
	"github.com/cheezecakee/ace/internal/engine"
	ctx "github.com/cheezecakee/ace/internal/ui/context"
)

func NewQuestionUI(
	q engine.Question,
	ctx *ctx.Context,
) QuestionUI {
	switch qt := q.(type) {
	case engine.ChoiceQuestion:
		return NewChoiceQuestionUI(qt, ctx)

	case engine.MultipleChoiceQuestion:
		return NewMultiChoiceQuestionUI(qt, ctx)

	case engine.BoolQuestion:
		return NewBoolQuestionUI(qt, ctx)

	case engine.TextEntryQuestion:
		return NewTextQuestionUI(qt, ctx)

	default:
		panic("unsupported question type")
	}
}
