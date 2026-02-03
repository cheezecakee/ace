package pack

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	packPath = "packs/"
)

var (
	ErrNotFound      = errors.New("file not found")
	ErrInvalidJSON   = errors.New("invalid JSON syntax")
	ErrMissingFields = errors.New("missing required fields")
	ErrNoQuestions   = errors.New("pack contains no questions")
	ErrInvalidData   = errors.New("invalid data")
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

func Unpack(data []byte) (*Raw, error) {
	var raw Raw
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, ErrInvalidJSON
	}

	return &raw, nil
}

func LoadAll() ([]*Pack, error) {
	files, err := os.ReadDir(packPath)
	if err != nil {
		return nil, err
	}

	var packs []*Pack
	prefix := "pack_"

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasPrefix(file.Name(), prefix) {
			continue
		}

		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(packPath, file.Name())

		// Use Load() which handles everything
		pack, err := Load(filePath)
		if err != nil {
			fmt.Printf("Warning: Failed to load %s: %v", filePath, err)
			continue
		}

		packs = append(packs, pack)
	}

	return packs, nil
}

func Load(filepath string) (*Pack, error) {
	data, err := Read(filepath)
	if err != nil {
		return nil, err
	}

	raw, err := Unpack(data)
	if err != nil {
		return nil, err
	}

	initialReport := raw.Verify()

	hasRepairableIssues := false
	for _, issue := range initialReport.Errors {
		if issue.Kind == IssueMissingID {
			hasRepairableIssues = true
			break
		}
	}

	if hasRepairableIssues {
		repairReport := raw.Repair()

		if repairReport.Repaired > 0 {
			if err := raw.Save(filepath); err != nil {
				return nil, fmt.Errorf("failed to save repaired pack: %w", err)
			}
		}

		finalReport := raw.Verify()

		if len(finalReport.Errors) > 0 {
			return nil, fmt.Errorf("pack validation failed: %d errors after repair", len(finalReport.Errors))
		}
	} else if len(initialReport.Errors) > 0 {
		return nil, fmt.Errorf("pack validation failed: %d errors", len(initialReport.Errors))
	}

	pack := raw.ToDomain(filepath)

	return pack, nil
}
