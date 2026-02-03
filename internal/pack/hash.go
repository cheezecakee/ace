package pack

import (
	"fmt"
	"hash/fnv"
)

type Identifier interface {
	Hash() string
	ID() string
}

type PackHash struct {
	Name    string
	Creator string
	Version string
}

type QuestionHash struct {
	PackHash   string
	Prompt     string
	Difficulty string
	Type       Type
	Category   string
	Index      int
}

func NewPackHash(name, creator, version string) Identifier {
	return &PackHash{
		Name:    name,
		Creator: creator,
		Version: version,
	}
}

func (h *PackHash) Hash() string {
	hash := fnv.New64a()
	hash.Write([]byte(h.Name))
	hash.Write([]byte(h.Creator))
	hash.Write([]byte(h.Version))
	sum := fmt.Sprintf("%x", hash.Sum64())
	return sum[:6]
}

func (h *PackHash) ID() string {
	return "pack-" + h.Hash()
}

func NewQuestionHash(packHash, prompt, difficulty, category string, qtype Type, index int) Identifier {
	return &QuestionHash{
		PackHash:   packHash,
		Prompt:     prompt,
		Difficulty: difficulty,
		Type:       qtype,
		Category:   category,
		Index:      index,
	}
}

func (h *QuestionHash) Hash() string {
	hash := fnv.New64a()
	hash.Write([]byte(h.Prompt))
	hash.Write([]byte(h.Difficulty))
	hash.Write([]byte(h.Category))
	sum := fmt.Sprintf("%x", hash.Sum64())
	return sum[:6]
}

func (h *QuestionHash) ID() string {
	return fmt.Sprintf("q-%s-%s-%02d-%s", h.PackHash, h.Type.String(), h.Index, h.Hash())
}
