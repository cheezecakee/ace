package session

import (
	"time"

	"github.com/cheezecakee/ace/internal/engine"
)

func (s *Session) HasStarted() bool {
	return s.GetState() != NotStarted
}

func (s *Session) GetCurrentQuestion() engine.Question {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.questions[s.currentIndex]
}

func (s *Session) GetCurrentIndex() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentIndex
}

func (s *Session) GetTimeRemaining() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.timeRemaining
}

func (s *Session) GetLivesRemaining() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.livesRemaining
}

func (s *Session) GetScore() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.score
}

func (s *Session) GetState() State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

func (s *Session) GetFormat() engine.Format {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.format
}

func (s *Session) IsCompleted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state.IsTerminal()
}

func (s *Session) CanNavigateBack() bool {
	return s.format.Time.Navigation == engine.Free
}

func (s *Session) IsAnswered(index int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if index < 0 || index >= len(s.answers) {
		return false
	}

	return s.answers[index] != nil
}

// GetElapsedTime returns total time since session started
func (s *Session) GetElapsedTime() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.state == NotStarted {
		return 0
	}

	if s.endTime.IsZero() {
		return time.Since(s.startTime)
	}

	return s.endTime.Sub(s.startTime)
}

func (s *Session) GetResults() Result {
	s.mu.RLock()
	defer s.mu.RUnlock()

	correct := 0
	for _, result := range s.gradeResults {
		if result != nil && result.IsCorrect() {
			correct++
		}
	}

	timeTaken := s.endTime.Sub(s.startTime)
	if s.endTime.IsZero() {
		timeTaken = time.Since(s.startTime)
	}

	return Result{
		TotalQuestions: len(s.questions),
		Correct:        correct,
		Incorrect:      len(s.questions) - correct,
		Score:          s.score,
		TimeTaken:      timeTaken,
		State:          s.state,
		Answers:        s.answers,
		GradeResults:   s.gradeResults,
	}
}
