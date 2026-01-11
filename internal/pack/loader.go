package pack

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrNotFound      = errors.New("file not found")
	ErrInvalidJSON   = errors.New("invalid JSON syntax")
	ErrMissingFields = errors.New("missing required fields")
	ErrNoQuestions   = errors.New("pack contains no questions")
	ErrInvalidData   = errors.New("invalid data")
	ErrDuplicateID   = errors.New("id already exists")
)

func Read(filepath string) ([]byte, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return data, nil
}

func Unpack(data []byte) (*Pack, error) {
	var pack Pack
	err := json.Unmarshal(data, &pack)
	if err != nil {
		return nil, ErrInvalidJSON
	}

	return &pack, nil
}

func Load(filepath string) (*Pack, error) {
	data, err := Read(filepath)
	if err != nil {
		return nil, err
	}

	pack, err := Unpack(data)
	if err != nil {
		return nil, err
	}

	pack.Repair()

	err = pack.Validate()
	if err != nil {
		return nil, err
	}

	pack.FilePath = filepath

	return pack, nil
}

func (p *Pack) Validate() error {
	if p.ID == "" || p.Name == "" || p.Role == "" || p.Creator == "" {
		return ErrMissingFields
	}

	if p.Count() <= 0 {
		return ErrNoQuestions
	}

	seenIDs := make(map[string]bool)
	for _, questions := range p.Categories {
		for i := range questions.Choice {
			if err := questions.Choice[i].Validate(seenIDs); err != nil {
				return err
			}
		}
		for i := range questions.MultipleChoice {
			if err := questions.MultipleChoice[i].Validate(seenIDs); err != nil {
				return err
			}
		}
		for i := range questions.TrueFalse {
			if err := questions.TrueFalse[i].Validate(seenIDs); err != nil {
				return err
			}
		}
		for i := range questions.TextEntry {
			if err := questions.TextEntry[i].Validate(seenIDs); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Pack) Count() int {
	c := 0
	for _, i := range p.Categories {
		c += len(i.Choice)
		c += len(i.MultipleChoice)
		c += len(i.TextEntry)
		c += len(i.TrueFalse)
	}

	return c
}

func (p *Pack) Repair() {
	if p.ID == "" {
		p.generateID()
	}

	for category, questions := range p.Categories {
		for i := range questions.Choice {
			questions.Choice[i].Repair(string(category), i)
		}

		for i := range questions.MultipleChoice {
			questions.MultipleChoice[i].Repair(string(category), i)
		}
		for i := range questions.TrueFalse {
			questions.TrueFalse[i].Repair(string(category), i)
		}
		for i := range questions.TextEntry {
			questions.TextEntry[i].Repair(string(category), i)
		}
	}
}

func (p *Pack) generateID() string {
	return fmt.Sprintf("pack-%s-%s", p.Role, sanitize(p.Creator))
}

func sanitize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")

	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func (p *Pack) Info() *PackInfo {
	return &PackInfo{
		ID:            p.ID,
		Name:          p.Name,
		Creator:       p.Creator,
		Version:       p.Version,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
		Role:          p.Role,
		Categories:    p.getCategories(),
		FilePath:      p.FilePath,
		QuestionCount: p.Count(),
	}
}

func (p *Pack) getCategories() Categories {
	var categories Categories
	for category := range p.Categories {
		categories = append(categories, category)
	}

	return categories
}
