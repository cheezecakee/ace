// Package engine
package engine

type Game struct {
	Format Format
}

type Difficulty int

const (
	Entry Difficulty = iota + 1
	Junior
	Mid
	Senior
)

func (d Difficulty) String() string {
	switch d {
	case Entry:
		return "entry"
	case Junior:
		return "junior"
	case Mid:
		return "mid"
	case Senior:
		return "senior"
	default:
		return ""
	}
}

func (d Difficulty) IsValid() bool {
	return d >= Entry && d <= Senior
}

func ParseDifficulty(s string) Difficulty {
	switch s {
	case "entry":
		return Entry
	case "junior":
		return Junior
	case "mid":
		return Mid
	case "senior":
		return Senior
	default:
		return 0
	}
}

type (
	QuestionType    int
	QuestionTypeSet []QuestionType
)

const (
	Choice QuestionType = iota + 1
	MultipleChoice
	TextEntry
	TrueFalse
)

func (qt QuestionType) String() string {
	switch qt {
	case Choice:
		return "choice"
	case MultipleChoice:
		return "multiple choice"
	case TextEntry:
		return "text entry"
	case TrueFalse:
		return "true/false"
	default:
		return "unknown"
	}
}

type QuestionCount int

const (
	AllQuestions QuestionCount = iota
	Ten
	Thirty
	Fifty
)

func (qc QuestionCount) Int() int {
	switch qc {
	case Ten:
		return 10
	case Thirty:
		return 30
	case Fifty:
		return 50
	default:
		return 0
	}
}

type TimeMode int

const (
	Unlimited            TimeMode = iota
	PerQuestion                   // e.g 30s per question
	TotalTime                     // e.g 15 minutes total
	PerQuestionWithBonus          // e.g 1|30 style Just like chess for every right question you get extra time!
)

func (tc TimeMode) String() string {
	switch tc {
	case Unlimited:
		return "unlimited"
	case PerQuestion:
		return "per question"
	case TotalTime:
		return "total time"
	case PerQuestionWithBonus:
		return "1|30"
	default:
		return ""
	}
}

type Progression int

const (
	Fixed Progression = iota
	Scaling
)

func (dm Progression) String() string {
	// Questions don't get harder
	switch dm {
	case Fixed:
		return "fixed"
	case Scaling:
		return "scaling"
	default:
		return ""
	}
}

type Navigation int

const (
	Free Navigation = iota
	Locked
)

func (nm Navigation) String() string {
	// Whether you can navigate back to the previous question
	switch nm {
	case Free:
		return "free"
	case Locked:
		return "locked"
	default:
		return ""
	}
}

type (
	Category    string
	CategorySet []Category
)

type LifeMode int

const (
	NoLives LifeMode = iota
	FixedLives
	SuddenDeath
)
