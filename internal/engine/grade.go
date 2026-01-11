package engine

type GradePolicy interface {
	Grade(q Question, a Answer) GradeResult
}

type GradeResult interface {
	IsCorrect() bool
	Type() GradeType
}

type GradeType int

const (
	Accuracy GradeType = iota + 1
	Binary
	Score
)

type AccuracyResult struct {
	Correct  bool
	Accuracy float32 // 0-1
	Feedback string  // What was good/missing
}

func (r AccuracyResult) IsCorrect() bool { return r.Correct }

func (r AccuracyResult) Type() GradeType { return Accuracy }

type BinaryResult struct {
	Correct bool
}

func (r BinaryResult) IsCorrect() bool { return r.Correct }

func (r BinaryResult) Type() GradeType { return Binary }

type ScoreResult struct {
	Correct      bool
	PointsEarned int
	MaxPoints    int
	Multiplier   float32 // eg., 1.5x for speed bonus
}

func (r ScoreResult) IsCorrect() bool { return r.Correct }

func (r ScoreResult) Type() GradeType { return Score }

type PracticeResult struct {
	CorrectAnswer string
	Explanation   string
}

func (r PracticeResult) IsCorrect() bool { return true } // Always "correct" in practice
func (r PracticeResult) Type() GradeType { return 0 }    // Zero value
