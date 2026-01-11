// Package session
package session

import (
	"errors"
	"time"

	"github.com/cheezecakee/ace/internal/engine"
)

var (
	ErrNotRunning = errors.New("session not running")
	ErrNoPause    = errors.New("can only pause running session")
)

type State int

const (
	NotStarted State = iota + 1
	Running
	Paused
	Completed
	Failed
)

type Session struct {
	game  engine.Game
	state State

	questions    []engine.Question
	answers      []engine.Answer
	currentIndex int
	gradeResults []engine.GradeResult

	// data fields
	score          int
	livesRemaining int

	startTime     time.Time
	pausedAt      time.Time
	timeRemaining time.Duration

	grader engine.GradePolicy
}

func (s *Session) NextQuestion() error {
	if s.state != Running {
		return ErrNotRunning
	}

	s.currentIndex++
	return nil
}

func (s *Session) Pause() error {
	if s.state != Running {
		return ErrNoPause
	}

	s.state = Paused
	return nil
}

func NewSession() *Session {
	return &Session{}
}
