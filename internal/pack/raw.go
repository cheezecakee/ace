// Package pack - raw converts raw data to domain data and includes features that verifies if the format of the pack is correct, if there are any missing ids it generates ones, and if there are any missing fields.
package pack

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cheezecakee/ace/internal/engine"
)

type Raw struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Version   string    `json:"version"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Categories map[string]RawCategory `json:"categories"` // category name -> questions
}

type RawCategory struct {
	Choice         []RawChoiceQuestion `json:"choice"`
	MultipleChoice []RawMultiQuestion  `json:"multiple_choice"`
	Bool      []RawBoolQuestion   `json:"bool"`
	TextEntry      []RawTextQuestion   `json:"text_entry"`
}

type RawChoiceQuestion struct {
	ID         string   `json:"id"`
	Difficulty string   `json:"difficulty"`
	Prompt     string   `json:"prompt"`
	Options    []string `json:"options"`
	Answer     int      `json:"answer"`
}

type RawMultiQuestion struct {
	ID         string   `json:"id"`
	Difficulty string   `json:"difficulty"`
	Prompt     string   `json:"prompt"`
	Options    []string `json:"options"`
	Answer     []int    `json:"answer"`
}

type RawBoolQuestion struct {
	ID         string `json:"id"`
	Difficulty string `json:"difficulty"`
	Prompt     string `json:"prompt"`
	Answer     bool   `json:"answer"`
}

type RawTextQuestion struct {
	ID         string   `json:"id"`
	Difficulty string   `json:"difficulty"`
	Prompt     string   `json:"prompt"`
	Expected   string   `json:"expected"` // Renamed this, was expected_answer before
	Keywords   []string `json:"keywords"`
}

func (r *Raw) Save(filepath string) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal pack: %w", err)
	}

	err = os.WriteFile(filepath, data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write pack: %w", err)
	}

	return nil
}

/** VERIFY **/

func (r *Raw) Verify() Report {
	var report Report

	// Verify pack-level fields
	if r.ID == "" {
		report.Errors = append(report.Errors, NewError(
			IssueMissingID,
			"Missing pack ID",
			"pack",
			"",
		))
	}

	if r.Name == "" {
		report.Errors = append(report.Errors, NewError(
			IssueMissingField,
			"Missing pack name",
			"pack",
			r.ID,
		))
	}

	if r.Role == "" {
		report.Errors = append(report.Errors, NewError(
			IssueMissingField,
			"Missing pack role",
			"pack",
			r.ID,
		))
	}

	if r.Creator == "" {
		report.Errors = append(report.Errors, NewError(
			IssueMissingField,
			"Missing pack creator",
			"pack",
			r.ID,
		))
	}

	if len(r.Categories) == 0 {
		report.Errors = append(report.Errors, NewError(
			IssueMissingField,
			"Pack has no categories",
			"pack",
			r.ID,
		))
	}

	// Check for duplicate IDs across all questions
	seenIDs := make(map[string]bool)

	// Verify all questions in all categories
	for categoryName, category := range r.Categories {
		// Verify Choice questions
		for i, q := range category.Choice {
			qReport := q.Verify()
			report.merge(qReport)

			if q.ID != "" {
				if seenIDs[q.ID] {
					report.Errors = append(report.Errors, NewError(
						IssueDuplicateID,
						fmt.Sprintf("Duplicate question ID: %s", q.ID),
						fmt.Sprintf("categories.%s.choice[%d]", categoryName, i),
						q.ID,
					))
				}
				seenIDs[q.ID] = true
			}
		}

		// Verify Multi questions
		for i, q := range category.MultipleChoice {
			qReport := q.Verify()
			report.merge(qReport)

			if q.ID != "" {
				if seenIDs[q.ID] {
					report.Errors = append(report.Errors, NewError(
						IssueDuplicateID,
						fmt.Sprintf("Duplicate question ID: %s", q.ID),
						fmt.Sprintf("categories.%s.multiple_choice[%d]", categoryName, i),
						q.ID,
					))
				}
				seenIDs[q.ID] = true
			}
		}

		// Verify Bool questions
		for i, q := range category.Bool {
			qReport := q.Verify()
			report.merge(qReport)

			if q.ID != "" {
				if seenIDs[q.ID] {
					report.Errors = append(report.Errors, NewError(
						IssueDuplicateID,
						fmt.Sprintf("Duplicate question ID: %s", q.ID),
						fmt.Sprintf("categories.%s.bool[%d]", categoryName, i),
						q.ID,
					))
				}
				seenIDs[q.ID] = true
			}
		}

		// Verify Text questions
		for i, q := range category.TextEntry {
			qReport := q.Verify()
			report.merge(qReport)

			if q.ID != "" {
				if seenIDs[q.ID] {
					report.Errors = append(report.Errors, NewError(
						IssueDuplicateID,
						fmt.Sprintf("Duplicate question ID: %s", q.ID),
						fmt.Sprintf("categories.%s.text_entry[%d]", categoryName, i),
						q.ID,
					))
				}
				seenIDs[q.ID] = true
			}
		}
	}

	return report
}

/** REPAIR **/

// Repair fixes all repairable issues
func (r *Raw) Repair() Report {
	var report Report

	// Generate pack ID if missing
	if r.ID == "" {
		packHash := NewPackHash(r.Name, r.Creator, r.Version)
		r.ID = packHash.ID()
		report.Repaired++
	}

	// Repair all question IDs using hash system
	for categoryName, category := range r.Categories {
		// Repair Choice questions
		for i := range category.Choice {
			q := &category.Choice[i]
			if q.ID == "" {
				qHash := NewQuestionHash(
					r.ID,
					q.Prompt,
					q.Difficulty,
					categoryName,
					TypeChoice,
					i,
				)
				q.ID = qHash.ID()
				report.Repaired++
			}
		}

		// Repair Multi questions
		for i := range category.MultipleChoice {
			q := &category.MultipleChoice[i]
			if q.ID == "" {
				qHash := NewQuestionHash(
					r.ID,
					q.Prompt,
					q.Difficulty,
					categoryName,
					TypeMulti,
					i,
				)
				q.ID = qHash.ID()
				report.Repaired++
			}
		}

		// Repair Bool questions
		for i := range category.Bool {
			q := &category.Bool[i]
			if q.ID == "" {
				qHash := NewQuestionHash(
					r.ID,
					q.Prompt,
					q.Difficulty,
					categoryName,
					TypeBool,
					i,
				)
				q.ID = qHash.ID()
				report.Repaired++
			}
		}

		// Repair Text questions
		for i := range category.TextEntry {
			q := &category.TextEntry[i]
			if q.ID == "" {
				qHash := NewQuestionHash(
					r.ID,
					q.Prompt,
					q.Difficulty,
					categoryName,
					TypeText,
					i,
				)
				q.ID = qHash.ID()
				report.Repaired++
			}
		}
	}

	return report
}

func (r *Raw) ToDomain(filepath string) *Pack {
	var questions []Question

	// Convert all questions from all categories
	for categoryName, category := range r.Categories {
		// Convert Choice questions
		for _, rawQ := range category.Choice {
			questions = append(questions, Question{
				ID:         rawQ.ID,
				Difficulty: engine.ParseDifficulty(rawQ.Difficulty),
				Category:   categoryName,
				Type:       TypeChoice,
				Prompt:     rawQ.Prompt,
				Answer: ChoiceAnswer{
					Options: rawQ.Options,
					Correct: rawQ.Answer,
				},
			})
		}

		// Convert MultipleChoice questions
		for _, rawQ := range category.MultipleChoice {
			questions = append(questions, Question{
				ID:         rawQ.ID,
				Difficulty: engine.ParseDifficulty(rawQ.Difficulty),
				Category:   categoryName,
				Type:       TypeMulti,
				Prompt:     rawQ.Prompt,
				Answer: MultiAnswer{
					Options: rawQ.Options,
					Correct: rawQ.Answer,
				},
			})
		}

		// Convert Bool questions
		for _, rawQ := range category.Bool {
			questions = append(questions, Question{
				ID:         rawQ.ID,
				Difficulty: engine.ParseDifficulty(rawQ.Difficulty),
				Category:   categoryName,
				Type:       TypeBool,
				Prompt:     rawQ.Prompt,
				Answer: BoolAnswer{
					Correct: rawQ.Answer,
				},
			})
		}

		// Convert TextEntry questions
		for _, rawQ := range category.TextEntry {
			questions = append(questions, Question{
				ID:         rawQ.ID,
				Difficulty: engine.ParseDifficulty(rawQ.Difficulty),
				Category:   categoryName,
				Type:       TypeText,
				Prompt:     rawQ.Prompt,
				Answer: TextAnswer{
					Expected: rawQ.Expected,
					Keywords: rawQ.Keywords,
				},
			})
		}
	}

	// Extract unique categories
	categorySet := make(map[string]bool)
	for cat := range r.Categories {
		categorySet[cat] = true
	}
	categories := make(Categories, 0, len(categorySet))
	for cat := range categorySet {
		categories = append(categories, Category(cat))
	}

	return &Pack{
		Info: Info{
			ID:         r.ID,
			Name:       r.Name,
			Role:       Role(r.Role),
			Categories: categories,
			Version:    r.Version,
			Creator:    r.Creator,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
			Path:       filepath,
			Count:      len(questions),
		},
		Questions: questions,
	}
}
