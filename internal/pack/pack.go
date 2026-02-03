// Package pack
package pack

import (
	"time"
)

type (
	Category   string
	Categories []Category
	Role       string
	Roles      []Role
)

type Packs []Pack

type Pack struct {
	Info      Info
	Questions Questions
}

type Info struct {
	// Indentity
	ID         string
	Name       string
	Role       Role
	Categories Categories

	// Metadata
	Version   string
	Creator   string
	CreatedAt time.Time
	UpdatedAt time.Time

	// Storage
	Path  string
	Count int
}
