package engine

import (
	"time"
)

type TimeRules struct {
	Control       TimeMode
	TotalDuration time.Duration // 0 = unlimited
	PerQuestion   time.Duration // 0 = unlimited
	Bonus         time.Duration
	Penalty       time.Duration
	Navigation    Navigation
}

type TimeOptions struct {
	TotalDuration time.Duration // 0 = unlimited
	PerQuestion   time.Duration // 0 = unlimited
	Bonus         time.Duration
	Penalty       time.Duration
}

type LifeRules struct {
	Enabled     bool
	Starting    int
	LoseOnWrong bool
}

type LifeOptions struct {
	Lives int
}

type ProgressionRules struct {
	Mode       Progression
	Difficulty Difficulty
}

type QuestionRules struct {
	CategoryFilter CategorySet
	Types          QuestionTypeSet
	Randomize      bool
	Count          QuestionCount
}

func BuildTimeRules(
	control TimeMode,
	opt TimeOptions,
) TimeRules {
	switch control {
	case Unlimited:
		return TimeRules{
			Control:    Unlimited,
			Navigation: Free,
		}
	case TotalTime:
		return TimeRules{
			Control:       TotalTime,
			TotalDuration: opt.TotalDuration,
			Navigation:    Free,
		}
	case PerQuestion:
		return TimeRules{
			Control:     PerQuestion,
			PerQuestion: opt.PerQuestion,
			Navigation:  Locked,
		}
	case PerQuestionWithBonus:
		return TimeRules{
			Control:     PerQuestionWithBonus,
			PerQuestion: opt.PerQuestion,
			Bonus:       opt.Bonus,
			Penalty:     opt.Penalty,
			Navigation:  Locked,
		}
	default:
		panic("unknown time control")
	}
}

func BuildLifeRules(
	mode LifeMode,
	opt LifeOptions,
) LifeRules {
	switch mode {
	case NoLives:
		return LifeRules{}
	case FixedLives:
		return LifeRules{
			Enabled:     true,
			Starting:    opt.Lives,
			LoseOnWrong: true,
		}
	case SuddenDeath:
		return LifeRules{
			Enabled:     true,
			Starting:    1,
			LoseOnWrong: true,
		}
	default:
		panic("unknown life mode")
	}
}
