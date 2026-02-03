package pack

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cheezecakee/ace/internal/engine"
)

const (
	cachePath         = "cache/"
	questionCachePath = cachePath + "question.json"
	lookupCachePath   = cachePath + "lookup.json"
)

type (
	TypeIndex map[Type][]string // Type -> QuestionIDs
	RoleIndex map[Role]TypeIndex
)

type (
	QuestionIndex map[string]engine.Question // QuestionID -> Question
	Lookup        map[engine.Difficulty]RoleIndex
)

type Cache interface {
	Load() error
	Save() error
	Generate(m Metadata, packIDs []string) error
}

func (c QuestionIndex) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal question index: %w", err)
	}

	err = os.WriteFile(questionCachePath, data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write question index: %w", err)
	}

	return nil
}

func (c *QuestionIndex) Load() error {
	data, err := os.ReadFile(questionCachePath)
	if err != nil {
		return fmt.Errorf("failed to read question index: %w", err)
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal question index: %w", err)
	}

	return nil
}

func (c *QuestionIndex) Generate(m Metadata, packIDs []string) error {
	// Clear existing cache
	*c = make(QuestionIndex)

	// Only load specified packs
	for _, packID := range packIDs {
		pack, err := m.LoadPack(packID)
		if err != nil {
			return fmt.Errorf("failed to load pack %s: %w", packID, err)
		}

		// Add all questions from this pack
		for _, q := range pack.Questions {
			(*c)[q.ID] = q.ToEngine()
		}
	}

	return c.Save()
}

func (c QuestionIndex) Fetch(ids []string) engine.Questions {
	var questions engine.Questions

	for _, id := range ids {
		questions = append(questions, c[id])
	}

	return questions
}

func (c Lookup) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal lookup index: %w", err)
	}

	err = os.WriteFile(lookupCachePath, data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write lookup index: %w", err)
	}

	return nil
}

func (c *Lookup) Load() error {
	data, err := os.ReadFile(lookupCachePath)
	if err != nil {
		return fmt.Errorf("failed to read lookup index: %w", err)
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal lookup index: %w", err)
	}

	return nil
}

func (c *Lookup) Generate(m Metadata, packIDs []string) error {
	// Clear existing lookup
	*c = make(Lookup)

	// Only load specified packs
	for _, packID := range packIDs {
		pack, err := m.LoadPack(packID)
		if err != nil {
			return fmt.Errorf("failed to load pack %s: %w", packID, err)
		}

		role := pack.Info.Role

		// Process each question in the pack
		for _, q := range pack.Questions {
			// Initialize nested maps if they don't exist
			if (*c)[q.Difficulty] == nil {
				(*c)[q.Difficulty] = make(RoleIndex)
			}
			if (*c)[q.Difficulty][role] == nil {
				(*c)[q.Difficulty][role] = make(TypeIndex)
			}

			// Append question ID to the appropriate bucket
			(*c)[q.Difficulty][role][q.Type] = append(
				(*c)[q.Difficulty][role][q.Type],
				q.ID,
			)
		}
	}

	return c.Save()
}

func (c Lookup) GetQuestionIDs(
	difficulty engine.Difficulty,
	role Role,
	types []Type,
) []string {
	var ids []string

	// Check if difficulty and role exist
	roleIdx, ok := c[difficulty]
	if !ok {
		return ids
	}

	typeIdx, ok := roleIdx[role]
	if !ok {
		return ids
	}

	// Collect IDs for all requested types
	for _, t := range types {
		if questionIDs, ok := typeIdx[t]; ok {
			ids = append(ids, questionIDs...)
		}
	}

	return ids
}

// GetAvailableRoles returns all roles that have questions available
// for the given difficulty and question types
func (c Lookup) GetAvailableRoles(
	difficulty engine.Difficulty,
	types []Type,
) []Role {
	var roles []Role

	// Check if difficulty exists in lookup
	roleIdx, ok := c[difficulty]
	if !ok {
		return roles
	}

	// Iterate through each role and check if it has questions for the requested types
	for role, typeIdx := range roleIdx {
		hasQuestions := false

		// Check if this role has any questions for the requested types
		for _, t := range types {
			if questionIDs, exists := typeIdx[t]; exists && len(questionIDs) > 0 {
				hasQuestions = true
				break
			}
		}

		// If this role has questions, add it to the result
		if hasQuestions {
			roles = append(roles, role)
		}
	}

	return roles
}
