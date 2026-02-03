package engine

import "time"

type Mode interface {
	Format(difficulty Difficulty) Format
}

func GetGameMode(mode ModeID) Mode {
	switch mode {
	case CustomMode:
		return newCustomMode()
	case RapidMode:
		return newRapidMode()
	case QuickMode:
		return newQuickMode()
	case HardcoreMode:
		return newHardcoreMode()
	case StandardMode:
		return newStandardMode()
	default:
		return nil
	}
}

type Standard struct{}

func newStandardMode() *Standard {
	return &Standard{}
}

func (gm *Standard) Format(difficulty Difficulty) Format {
	return Format{
		Time: BuildTimeRules(
			Unlimited,
			TimeOptions{},
		),
		Lives: BuildLifeRules(
			NoLives,
			LifeOptions{},
		),
		Progression: ProgressionRules{
			Mode:       Scaling,
			Difficulty: difficulty,
		},
		Question: QuestionRules{
			Types:     QuestionTypeSet{TextEntry},
			Randomize: true,
		},
		Description: "Interview-style practice with grading",
	}
}

type Quick struct{}

func newQuickMode() *Quick {
	return &Quick{}
}

func (gm *Quick) Format(difficulty Difficulty) Format {
	return Format{
		Time: BuildTimeRules(
			TotalTime,
			TimeOptions{
				TotalDuration: 300 * time.Second, // 5 minutes
			},
		),
		Lives: BuildLifeRules(
			NoLives,
			LifeOptions{},
		),
		Progression: ProgressionRules{
			Mode:       Fixed,
			Difficulty: difficulty,
		},
		Question: QuestionRules{
			Types:     QuestionTypeSet{Choice, MultipleChoice, Bool},
			Randomize: true,
		},
		Description: "Quick-fire questions to warm up before interviews",
	}
}

type Rapid struct{}

func newRapidMode() *Rapid {
	return &Rapid{}
}

func (gm *Rapid) Format(difficulty Difficulty) Format {
	return Format{
		Time: BuildTimeRules(
			PerQuestion,
			TimeOptions{
				PerQuestion: 40 * time.Second,
				Penalty:     5 * time.Second,
			},
		),
		Lives: BuildLifeRules(
			NoLives,
			LifeOptions{},
		),
		Progression: ProgressionRules{
			Mode:       Fixed,
			Difficulty: difficulty,
		},
		Question: QuestionRules{
			Types:     QuestionTypeSet{Choice, MultipleChoice, Bool},
			Randomize: true,
		},
		Description: "Fast-paced reaction training",
	}
}

type Hardcore struct{}

func newHardcoreMode() *Hardcore {
	return &Hardcore{}
}

func (gm *Hardcore) Format(difficulty Difficulty) Format {
	return Format{
		Time: BuildTimeRules(
			PerQuestionWithBonus,
			TimeOptions{
				PerQuestion: 30 * time.Second,
				Bonus:       10 * time.Second,
			},
		),
		Lives: BuildLifeRules(
			SuddenDeath,
			LifeOptions{},
		),
		Progression: ProgressionRules{
			Mode:       Scaling,
			Difficulty: Entry,
		},
		Question: QuestionRules{
			Types:     QuestionTypeSet{Choice, MultipleChoice, Bool},
			Randomize: true,
		},
		Description: "High-pressure survival mode",
	}
}

type Custom struct {
	Control TimeMode // default Unlimited

	// Time settings
	TotalDuration int // seconds
	PerQuestion   int // seconds
	Bonus         int // seconds
	Penalty       int // seconds

	// Life settings
	LifeMode LifeMode // default NoLives
	Lives    int

	// Progression settings
	Progression Progression // default Fixed

	// Question filters
	Categories    CategorySet     // default all
	Types         QuestionTypeSet // default all
	Randomize     bool
	QuestionCount QuestionCount // 0 = all
}

func newCustomMode() *Custom {
	return &Custom{
		Control:       Unlimited,
		LifeMode:      NoLives,
		Progression:   Fixed,
		Randomize:     false,
		QuestionCount: 0,
	}
}

func (gm *Custom) Format(difficulty Difficulty) Format {
	return Format{
		Time: BuildTimeRules(
			gm.Control,
			TimeOptions{
				TotalDuration: time.Duration(gm.TotalDuration) * time.Second,
				PerQuestion:   time.Duration(gm.PerQuestion) * time.Second,
				Bonus:         time.Duration(gm.Bonus) * time.Second,
				Penalty:       time.Duration(gm.Penalty) * time.Second,
			},
		),
		Lives: BuildLifeRules(
			gm.LifeMode,
			LifeOptions{
				Lives: gm.Lives,
			},
		),
		Progression: ProgressionRules{
			Mode:       gm.Progression,
			Difficulty: difficulty,
		},
		Question: QuestionRules{
			CategoryFilter: gm.Categories,
			Types:          gm.Types,
			Randomize:      gm.Randomize,
			Count:          gm.QuestionCount,
		},
		Description: "Custom practice mode",
	}
}
