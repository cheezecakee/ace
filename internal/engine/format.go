package engine

import "errors"

// Formats controls the game type
// --- Time ---
// 1. Time per questions: Each question has it's own time limit when time expires you can't check it again
// 2. Time duration: The whole quiz needs to be completed within a time limit, navigating through questions is allowed
// 3. Bonus on correct: Start with a specific time, and gain more time for every question answered correctly, every wrong question decreases that time
// --- Life ---
// 1. No lives: No penalties for getting a question wrong
// 2. Fixed lives: Start with a fixed number of lives, incorrect answers decreases the number of lives.
// 3. Sudden death: One life, any wrong answers results in fails
// --- ---
// Some of these modes can be mixed and matched

type Format struct {
	Time        TimeRules
	Lives       LifeRules
	Progression ProgressionRules
	Question    QuestionRules
	Description string
}

func (f Format) Validate() error {
	// --- Time ---
	switch f.Time.Control {
	case TotalTime:
		if f.Time.TotalDuration <= 0 {
			return errors.New("total time must be > 0")
		}
	case PerQuestion:
		if f.Time.PerQuestion <= 0 {
			return errors.New("per-question time must be > 0")
		}
	case PerQuestionWithBonus:
		if f.Time.PerQuestion <= 0 {
			return errors.New("per-question time must be > 0")
		}
		if f.Time.Bonus <= 0 {
			return errors.New("bonus time must be > 0")
		}
		if f.Time.Bonus >= f.Time.PerQuestion {
			return errors.New("bonus time should be less than per-question time")
		}
	}

	// --- Lives ---
	if !f.Lives.Enabled {
		if f.Lives.Starting != 0 {
			return errors.New("starting life must be 0 when lives disabled")
		}
	} else if f.Lives.Starting <= 0 {
		return errors.New("start life must be > 0")
	}

	// --- Questions ---
	if f.Question.Count < 0 {
		return errors.New("max questions cannot be negative")
	}

	if len(f.Question.Types) == 0 {
		return errors.New("at least one question type must be selected")
	}

	return nil
}
