package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gen2brain/beeep"
)

func main() {
	model := &Model{}
	app := tea.NewProgram(model, tea.WithAltScreen())

	var err error
	_, err = app.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if model.state == done {
		beeep.Notify("TeaTime", "Tea timer done!", "")
		beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
	}
}
