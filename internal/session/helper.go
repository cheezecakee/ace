package session

import (
	"time"

	"github.com/cheezecakee/ace/internal/engine"
)

func (s *Session) checkGameOver() bool {
	// Out of lives
	if s.format.Lives.Enabled && s.livesRemaining <= 0 {
		s.state = Failed
		s.endTime = time.Now()
		return true
	}

	// All questions answered
	if s.isLastQuestion() && s.answers[s.currentIndex] != nil {
		s.state = Completed
		s.endTime = time.Now()
		return true
	}

	return false
}

func (s *Session) advance() {
	if s.currentIndex < len(s.questions)-1 {
		s.currentIndex++
		s.resetQuestionTimer()
	}
}

func (s *Session) isLastQuestion() bool {
	return s.currentIndex == len(s.questions)-1
}

func (s *Session) resetQuestionTimer() {
	switch s.format.Time.Control {
	case engine.PerQuestion, engine.PerQuestionWithBonus:
		s.timeRemaining = s.format.Time.PerQuestion
	}
}

func (s *Session) applyTimeBonus() {
	if s.format.Time.Control != engine.PerQuestionWithBonus {
		return
	}

	s.timeRemaining += s.format.Time.Bonus
}

func (s *Session) applyTimePenalty() {
	if s.format.Time.Control != engine.PerQuestionWithBonus {
		return
	}

	s.timeRemaining -= s.format.Time.Penalty
	if s.timeRemaining < 0 {
		s.timeRemaining = 0
	}
}
