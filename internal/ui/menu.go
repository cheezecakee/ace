// Package ui
package ui

type Menu struct {
	Quickstart bool
	Custom     bool
	Settings   Settings
}

func (m Menu) Start() {}
func (m Menu) Quit()  {}

type Settings struct {
	Import   string // Imports a json
	Language string
	GameMode string   // Default Gamemode, change name later
	Packs    []string // Which packs preload when starting
}

type (
	Category   string
	Categories []Category
	Roles      map[string]Categories
)

type Defaults struct {
	Roles      Roles // includes the role name and the categories for that role
	Categories       // includes all the categories available
}
