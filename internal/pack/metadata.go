package pack

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	metadataPath = "packs/metadata.json"
)

type (
	Catalog    map[Role][]string // role -> pack IDs
	PacksIndex map[string]Info   // pack ID -> pack info
)

type Metadata struct {
	Questions  int        `json:"questions"`  // Total number of questions across all packs
	Categories Categories `json:"categories"` // All unique categories
	Roles      Roles      `json:"roles"`      // All unique roles

	Catalog Catalog `json:"catalog"`

	Packs PacksIndex `json:"packs"`
}

func Build(packs []*Pack) *Metadata {
	m := &Metadata{
		Questions:  0,
		Categories: Categories{},
		Roles:      Roles{},
		Catalog:    make(Catalog),
		Packs:      make(PacksIndex),
	}

	categorySet := make(map[Category]bool)
	roleSet := make(map[Role]bool)

	for _, pack := range packs {
		// Add pack info
		m.Packs[pack.Info.ID] = pack.Info

		// Count questions
		m.Questions += pack.Info.Count

		// Track role
		roleSet[pack.Info.Role] = true

		// Add pack to catalog for this role
		m.Catalog[pack.Info.Role] = append(m.Catalog[pack.Info.Role], pack.Info.ID)

		// Track categories
		for _, cat := range pack.Info.Categories {
			categorySet[cat] = true
		}
	}

	// Convert sets to slices
	for cat := range categorySet {
		m.Categories = append(m.Categories, cat)
	}
	for role := range roleSet {
		m.Roles = append(m.Roles, role)
	}

	// Save locally
	m.Save()

	return m
}

func (m *Metadata) LoadPack(packID string) (*Pack, error) {
	info, exists := m.Packs[packID]
	if !exists {
		return nil, fmt.Errorf("pack %s not found in metadata", packID)
	}

	// Load pack from file
	pack, err := Load(info.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to load pack %s: %w", packID, err)
	}

	return pack, nil
}

func (m *Metadata) Save() error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = os.WriteFile(metadataPath, data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

func (m *Metadata) Load() error {
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	err = json.Unmarshal(data, m)
	if err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return nil
}

func (m *Metadata) GetRoles() Roles {
	return m.Roles
}

func (m *Metadata) GetCategories() Categories {
	return m.Categories
}

func (m *Metadata) GetPacksByRole(role Role) []Info {
	packIDs, exists := m.Catalog[role]
	if !exists {
		return []Info{}
	}

	infos := make([]Info, 0, len(packIDs))
	for _, id := range packIDs {
		if info, exists := m.Packs[id]; exists {
			infos = append(infos, info)
		}
	}

	return infos
}

func (m *Metadata) GetCategoriesByRole(role Role) []string {
	packs := m.GetPacksByRole(role)

	categorySet := make(map[string]bool)
	for _, packInfo := range packs {
		for _, cat := range packInfo.Categories {
			categorySet[string(cat)] = true
		}
	}

	categories := make([]string, 0, len(categorySet))
	for cat := range categorySet {
		categories = append(categories, cat)
	}

	return categories
}

func (m *Metadata) PackIDs() []string {
	var packIDs []string
	for _, pack := range m.Packs {
		packIDs = append(packIDs, pack.ID)
	}

	return packIDs
}

// ActivePacks TODO Comeback to this later
func (m *Metadata) ActivePacks(active []Info) []string {
	packs := make(map[string]bool)
	for _, id := range m.PackIDs() {
		packs[id] = false
	}

	return nil
}
