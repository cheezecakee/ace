package pack

import (
	"slices"
	"time"
)

type Metadata struct {
	CompiledAt time.Time
	Packs      []PackInfo
	Roles      Role
	Categories Categories
	Questions  int
}

func Build(packs []*Pack) *Metadata {
	var (
		packsInfo  []PackInfo
		categories Categories
	)

	questions := 0
	roles := make(Role)

	for _, pack := range packs {
		packInfo := pack.Info()

		for _, category := range packInfo.Categories {
			if !slices.Contains(categories, category) {
				categories = append(categories, category)
			}
			// Add pack's categories for this specific role
			if !slices.Contains(roles[packInfo.Role], category) {
				roles[packInfo.Role] = append(roles[packInfo.Role], category)
			}
		}

		questions += packInfo.QuestionCount
		packsInfo = append(packsInfo, *packInfo)
	}

	return &Metadata{
		CompiledAt: time.Now(),
		Packs:      packsInfo,
		Roles:      roles,
		Categories: categories,
		Questions:  questions,
	}
}
