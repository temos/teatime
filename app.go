package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var durations = []time.Duration{
	time.Second * 3,
	time.Second * 30,
	time.Second * 45,
	time.Second * 60,
	time.Second * (60 + 30),
	time.Second * 60 * 2,
	time.Second * 60 * 3,
}

type state int

const (
	pickingDuration state = iota
	customDuration
	running
	done
)

type Model struct {
	state       state
	durationIdx int
	remaining   time.Duration
	input       textinput.Model
}

type TickMsg struct{}

func TickIn(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

type DoneMsg struct{}

func (m *Model) Init() tea.Cmd {
	m.state = pickingDuration

	m.input = textinput.New()
	m.input.Prompt = "Custom duration in GO format: "
	m.input.Placeholder = "1m30s"
	m.input.Focus()

	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	commands := make([]tea.Cmd, 0, 1)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "down":
			if m.state == pickingDuration {
				m.durationIdx = (m.durationIdx + 1) % (len(durations) + 1)
			}

		case "up":
			if m.state == pickingDuration {
				m.durationIdx = m.durationIdx - 1
				if m.durationIdx < 0 {
					m.durationIdx = len(durations)
				}
			}

		case "enter":
			if m.state == pickingDuration {
				if m.durationIdx == len(durations) {
					m.state = customDuration
				} else {
					m.state = running
					m.remaining = durations[m.durationIdx]

					commands = append(commands, TickIn(time.Second))
				}
			}

			if m.state == customDuration {
				dur, err := time.ParseDuration(m.input.Value())
				if err == nil {
					m.state = running
					m.remaining = dur.Round(time.Second)
					commands = append(commands, TickIn(time.Second))
				}
			}
		}

	case TickMsg:
		m.remaining -= time.Second
		if m.remaining == 0 {
			doneCmd := func() tea.Msg { return DoneMsg{} }
			commands = append(commands, doneCmd)
		} else {
			commands = append(commands, TickIn(time.Second))
		}

	case DoneMsg:
		m.state = done
		return m, tea.Quit
	}

	if m.state == customDuration {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		commands = append(commands, cmd)
	}

	return m, tea.Batch(commands...)
}

func (m *Model) View() string {
	switch m.state {
	case pickingDuration:
		var b strings.Builder

		for i, duration := range durations {
			if i == m.durationIdx {
				b.WriteString("[*] ")
			} else {
				b.WriteString("[ ] ")
			}

			b.WriteString(duration.String())
			b.WriteRune('\n')
		}

		if m.durationIdx == len(durations) {
			b.WriteString("[*] Custom duration\n")
		} else {
			b.WriteString("[ ] Custom duration\n")
		}

		return b.String()

	case customDuration:
		return m.input.View()

	case running:
		return m.remaining.String()

	case done:
		return ""

	default:
		panic("unknown state: " + strconv.FormatInt(int64(m.state), 10))
	}
}
