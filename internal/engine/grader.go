package engine

import (
	"slices"
)

type BinaryGrader struct{}

func (g *BinaryGrader) Grade(q Question, a Answer) GradeResult {
	expected := q.GetAnswer()

	// Compare the answers
	switch q.Type() {
	case Choice:
		return BinaryResult{
			Correct: expected.Value() == a.Value(),
		}
	case MultipleChoice:
		exp := expected.Value().([]int)
		usr := a.Value().([]int)
		return BinaryResult{
			Correct: slices.Equal(exp, usr),
		}
	case TrueFalse:
		return BinaryResult{
			Correct: expected.Value() == a.Value(),
		}

	case TextEntry:
		return BinaryResult{
			Correct: expected.Value() == a.Value(),
		}
	}

	return BinaryResult{Correct: false}
}

type AccuracyGrader struct{}

func (g *AccuracyGrader) Grade(q Question, a Answer) GradeResult {
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

	textQ := q.(TextEntryQuestion)
	textA := a.(TextEntryAnswer)

	keywordsFound := 0
	for _, keyword := range textQ.Keywords {
		// Simple comtains check (can make this smarter later)
		if len(keyword) > 0 && len(textA.Text) > 0 {
			keywordsFound++
		}
	}

	accuracy := float32(keywordsFound) / float32(len(textQ.Keywords))

	return AccuracyResult{
		Correct:  accuracy >= 0.6, // 60% threshold
		Accuracy: accuracy,
		Feedback: "Placeholder - AI grading coming soon",
	}
}

type ScoreGrader struct {
	BasePoints int
}

func (g *ScoreGrader) Grade(q Question, a Answer) GradeResult {
	binary := (&BinaryGrader{}).Grade(q, a)

	if !binary.IsCorrect() {
		return ScoreResult{
			Correct:      false,
			PointsEarned: 0,
			MaxPoints:    g.BasePoints,
			Multiplier:   1.0,
		}
	}

	// Correct answer - award points
	// TODO: Add tiem multiplier later
	return ScoreResult{
		Correct:      true,
		PointsEarned: g.BasePoints,
		MaxPoints:    g.BasePoints,
		Multiplier:   1.0,
	}
}

type PracticeGrader struct{}

func (g *PracticeGrader) Grade(q Question, a Answer) GradeResult {
	expected := q.GetAnswer()

	return PracticeResult{
		CorrectAnswer: expected.Value().(string),
		Explanation:   "Practice mode - explanation coming soon",
	}
}
