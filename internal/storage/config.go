package storage

import (
	"encoding/json"
	"os"
)

const (
	settingsFile = "savedata/preferences.json"
	// statsFile    = "savedata/stats.json"
	// sessionsFile = "savedata/sessions.json"
)

type User struct {
	Settings Settings `json:"settings"`
}

type Settings struct {
	Language    string   `json:"language"`
	ActivePacks []string `json:"active_packs"`
}

func NewUser() *User {
	u := &User{}
	u.defaultSettings()

	return u
}

func (u *User) Load() error {
	// Load other stuff here later on
	data, err := os.ReadFile(settingsFile)
	if err != nil {
		if os.IsNotExist(err) {
			u.defaultSettings()
			return nil
		}
		return err
	}

	if err := json.Unmarshal(data, u); err != nil {
		return err
	}

	return nil
}

func (u *User) Save() error {
	data, err := json.MarshalIndent(u, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsFile, data, 0o644)
}

func (u *User) defaultSettings() {
	u.Settings.Language = "en"
	u.Settings.ActivePacks = []string{}
}
