// Package context
package context

import (
	"github.com/cheezecakee/ace/internal/engine"
	"github.com/cheezecakee/ace/internal/pack"
	"github.com/cheezecakee/ace/internal/session"
	"github.com/cheezecakee/ace/internal/storage"
	"github.com/cheezecakee/ace/internal/ui"
)

type Context struct {
	Keys ui.KeyMap

	Format   engine.Format
	Session  *session.Session
	User     *storage.User
	Metadata *pack.Metadata

	QuestionCache pack.QuestionIndex
	LookupCache   pack.Lookup

	Mode  engine.ModeID
	Packs map[string]bool // if they are active or not

	Styles ui.Styles
	Width  int
	Height int
}

func NewContext() *Context {
	user := storage.NewUser()
	_ = user.Load() // TODO ignore err for now

	allPacks, err := pack.LoadAll()
	if err != nil {
		panic(err)
	}

	metadata := pack.Build(allPacks)

	packs := make(map[string]bool)

	for _, p := range metadata.Packs {
		packs[p.ID] = false
	}

	// Active Packs from user settings
	for _, id := range user.Settings.ActivePacks {
		if _, ok := packs[id]; ok {
			packs[id] = true
		}
	}

	ctx := &Context{
		Keys:     ui.DefaultKeyMap(),
		Mode:     engine.StandardMode,
		User:     user,
		Metadata: metadata,
		Packs:    packs,
	}

	// Build caches
	ctx.buildCache()

	return ctx
}

type (
	SetFormatMsg  engine.Format
	SetSessionMsg *session.Session
	SetPackMsg    []pack.Pack
)

// GetActivePacks returns slice of active pack IDs
func (c *Context) GetActivePacks() []string {
	var active []string
	for id, isActive := range c.Packs {
		if isActive {
			active = append(active, id)
		}
	}
	return active
}

// IsPackActive checks if a pack is active
func (c *Context) IsPackActive(packID string) bool {
	return c.Packs[packID]
}

// TogglePack toggles a pack's active state
func (c *Context) TogglePack(packID string) {
	if _, exists := c.Packs[packID]; exists {
		c.Packs[packID] = !c.Packs[packID]
	}
}

func (c *Context) buildCache() {
	c.QuestionCache = make(pack.QuestionIndex)
	c.LookupCache = make(pack.Lookup)

	activePacks := c.GetActivePacks()

	// Try loading from disk first
	qErr := c.QuestionCache.Load()
	lErr := c.LookupCache.Load()

	if qErr != nil || lErr != nil {
		// Cache doesn't exist or is invalid, generate it
		c.QuestionCache.Generate(*c.Metadata, activePacks)
		c.LookupCache.Generate(*c.Metadata, activePacks)
	}
}

func (c *Context) RebuildCache() error {
	activePacks := c.GetActivePacks()

	if err := c.QuestionCache.Generate(*c.Metadata, activePacks); err != nil {
		return err
	}

	if err := c.LookupCache.Generate(*c.Metadata, activePacks); err != nil {
		return err
	}

	return nil
}
