// Package session
package session

import (
	"errors"
	"sync"
	"time"

	"github.com/cheezecakee/ace/internal/engine"
)

var (
	ErrNotRunning   = errors.New("session not running")
	ErrAlreadyEnded = errors.New("session already ended")
	ErrNoNavigation = errors.New("navigation not allowed in this mode")
	ErrInvalidIndex = errors.New("invalid question index")
)

type Session struct {
	mu sync.RWMutex

	format engine.Format
	state  State

	questions    engine.Questions
	answers      []engine.Answer
	currentIndex int
	gradeResults []engine.GradeResult

	score          int
	livesRemaining int

	startTime     time.Time
	endTime       time.Time
	timeRemaining time.Duration

	grader engine.GradePolicy
}

type Result struct {
	TotalQuestions int
	Correct        int
	Incorrect      int
	Score          int
	TimeTaken      time.Duration
	State          State
	Answers        []engine.Answer
	GradeResults   []engine.GradeResult
}

func NewSession(format engine.Format, questions engine.Questions, grader engine.GradePolicy) *Session {
	return &Session{
		format:         format,
		questions:      questions,
		answers:        make([]engine.Answer, len(questions)),
		gradeResults:   make([]engine.GradeResult, len(questions)),
		state:          NotStarted,
		livesRemaining: format.Lives.Starting,
		grader:         grader,
	}
}

// Begin starts the session (no goroutines, just state initialization)
func (s *Session) Begin() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != NotStarted {
		return ErrAlreadyEnded
	}

	s.state = Running
	s.startTime = time.Now()

	// Initialize time remaining based on mode
	switch s.format.Time.Control {
	case engine.TotalTime:
		s.timeRemaining = s.format.Time.TotalDuration
	case engine.PerQuestion, engine.PerQuestionWithBonus:
		s.timeRemaining = s.format.Time.PerQuestion
	case engine.Unlimited:
		s.timeRemaining = 0
	}

	return nil
}

// SubmitAnswer handles answer submission synchronously
func (s *Session) SubmitAnswer(answer engine.Answer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != Running {
		return ErrNotRunning
	}

	// Store answer
	s.answers[s.currentIndex] = answer

	// Grade answer
	question := s.questions[s.currentIndex]
	result := s.grader.Grade(question, answer)
	s.gradeResults[s.currentIndex] = result

	// Update score and lives
	if result.IsCorrect() {
		s.score++
		s.applyTimeBonus()
	} else {
		if s.format.Lives.Enabled && s.format.Lives.LoseOnWrong {
			s.livesRemaining--
		}
		s.applyTimePenalty()
	}

	// Check game over conditions
	if s.checkGameOver() {
		return nil
	}

	// Auto-advance for locked navigation
	if s.format.Time.Navigation == engine.Locked {
		s.advance()
	}

	return nil
}

// NextQuestion moves to next question (manual navigation)
func (s *Session) NextQuestion() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != Running {
		return ErrNotRunning
	}

	if s.format.Time.Navigation != engine.Free {
		return ErrNoNavigation
	}

	if s.currentIndex >= len(s.questions)-1 {
		return ErrInvalidIndex
	}

	s.currentIndex++
	return nil
}

// PrevQuestion moves to previous question (manual navigation)
func (s *Session) PrevQuestion() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != Running {
		return ErrNotRunning
	}

	if s.format.Time.Navigation != engine.Free {
		return ErrNoNavigation
	}

	if s.currentIndex <= 0 {
		return ErrInvalidIndex
	}

	s.currentIndex--
	return nil
}

// Tick updates the game timer (called by UI on timer ticks)
func (s *Session) Tick(elapsed time.Duration) (timeExpired bool, newState State) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != Running {
		return false, s.state
	}

	if s.format.Time.Control != engine.Unlimited {
		s.timeRemaining -= elapsed

		if s.timeRemaining <= 0 {
			s.timeRemaining = 0
			s.state = TimeExpired
			s.endTime = time.Now()
			return true, TimeExpired
		}
	}

	return false, Running
}
