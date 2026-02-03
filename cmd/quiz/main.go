package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/ui/context"
	screen "github.com/cheezecakee/ace/internal/ui/screens"
)

func main() {
	ctx := context.NewContext()
	model := screen.NewModel(ctx)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
