package engine

import (
	"slices"
	"strings"
)

type BinaryGrader struct{}

func (g *BinaryGrader) Grade(q Question, a Answer) GradeResult {
	expected := q.GetAnswer()

	// Compare the answers
	switch q.Type() {
	case Choice, Bool, TextEntry:
		return BinaryResult{
			Correct: expected.Value() == a.Value(),
		}
	case MultipleChoice:
		exp := expected.Value().([]int)
		usr := a.Value().([]int)
		return BinaryResult{
			Correct: slices.Equal(exp, usr),
		}
	}

	return BinaryResult{Correct: false}
}

// AccuracyGrader - for TextEntry with partial credit
type AccuracyGrader struct{}

func (g *AccuracyGrader) Grade(q Question, a Answer) GradeResult {
	// for non-text questions, treat as binary
	if q.Type() != TextEntry {
		binary := (&BinaryGrader{}).Grade(q, a)
		isCorrect := binary.IsCorrect()

		accuracy := float32(0)
		if isCorrect {
			accuracy = 1.0
		}

		return AccuracyResult{
			Correct:  isCorrect,
			Accuracy: accuracy,
			Feedback: "",
		}
	}

	// For TextEntry, do keyword matching (placeholder until AI)
	// TextEntry grading
	textQ, ok := q.(TextEntryQuestion)
	if !ok {
		return AccuracyResult{Correct: false}
	}

	textA, ok := a.(TextEntryAnswer)
	if !ok {
		return AccuracyResult{Correct: false}
	}

	answerText := strings.ToLower(textA.Text)

	// Count how many keywords are present
	keywordsFound := 0
	for _, keyword := range textQ.Keywords {
		if strings.Contains(answerText, strings.ToLower(keyword)) {
			keywordsFound++
		}
	}

	accuracy := float32(0)
	if len(textQ.Keywords) > 0 {
		accuracy = float32(keywordsFound) / float32(len(textQ.Keywords))
	}

	return AccuracyResult{
		Correct:  accuracy == 1.0,
		Accuracy: accuracy,
		Feedback: "",
	}
}

type ScoreGrader struct{}

func (g *ScoreGrader) Grade(q Question, a Answer) GradeResult {
	expected := q.GetAnswer()

	// Compare the answers
	switch q.Type() {
	case Choice, Bool, TextEntry:
		return ScoreResult{
			Correct: expected.Value() == a.Value(),
		}
	case MultipleChoice:
		exp := expected.Value().([]int)
		usr := a.Value().([]int)
		return ScoreResult{
			Correct: slices.Equal(exp, usr),
		}
	}

	return ScoreResult{Correct: false}
}

type PracticeGrader struct{}

func (g *PracticeGrader) Grade(q Question, a Answer) GradeResult {
	return PracticeResult{
		CorrectAnswer: q.GetAnswer().Value().(string),
	}
}
