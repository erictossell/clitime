package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	stopwatches []stopwatch.Model
	current     int
	keymap      keymap
	help        help.Model
	quitting    bool
}

type keymap struct {
	start  key.Binding
	stop   key.Binding
	reset  key.Binding
	quit   key.Binding
	new    key.Binding
	delete key.Binding
	next   key.Binding
	prev   key.Binding
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	// Note: you could further customize the time output by getting the
	// duration from m.stopwatch.Elapsed(), which returns a time.Duration, and
	// skip m.stopwatch.View() altogether.
	if m.quitting {
		return ""
	}

	if len(m.stopwatches) == 0 {
		return "No tasks. Press 'n' to create a new task.\n" + m.helpView()
	}

	var s string
	for i, sw := range m.stopwatches {
		indicator := " "
		if i == m.current {
			indicator = ">"
		}
		s += fmt.Sprintf("%s Task %d: %s\n", indicator, i+1, sw.View())
	}

	s += m.helpView()
	return s
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
		m.keymap.new,
		m.keymap.delete,
		m.keymap.next,
		m.keymap.prev,
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyHandled := false // Add this flag to track if the key has been handled

		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.new):
			m.stopwatches = append(m.stopwatches, stopwatch.NewWithInterval(time.Millisecond))
			m.current = len(m.stopwatches) - 1
			keyHandled = true

		case key.Matches(msg, m.keymap.delete):
			if len(m.stopwatches) > 0 {
				m.stopwatches = append(m.stopwatches[:m.current], m.stopwatches[m.current+1:]...)
				if m.current >= len(m.stopwatches) {
					m.current = len(m.stopwatches) - 1
				}
			}
			keyHandled = true
		
    case key.Matches(msg, m.keymap.next):
			if m.current < len(m.stopwatches)-1 {
				m.current++
			}
			keyHandled = true
		
    case key.Matches(msg, m.keymap.prev):
			if m.current > 0 {
				m.current--
			}
			keyHandled = true
		
    case key.Matches(msg, m.keymap.reset):
			if len(m.stopwatches) > 0 {
				m.stopwatches[m.current], _ = m.stopwatches[m.current].Update(stopwatch.ResetMsg{})
			}
			keyHandled = true
		
    
    case key.Matches(msg, m.keymap.start):
	    if len(m.stopwatches) > 0 {
		    if !m.stopwatches[m.current].Running() {
			    msg = StartStopMsg{ID: m.stopwatches[m.current].ID(), running: true}
		    } else {
			    fmt.Println("Current stopwatch is already running")
		    }
	    }

    case key.Matches(msg, m.keymap.stop):
	    if len(m.stopwatches) > 0 && m.stopwatches[m.current].Running() {
		    msg = StartStopMsg{ID: m.stopwatches[m.current].ID(), running: false}
	    } 

		// Only update the stopwatch if the key hasn't been handled
		if !keyHandled && len(m.stopwatches) > 0 && m.current >= 0 && m.current < len(m.stopwatches) {
			m.stopwatches[m.current], cmd = m.stopwatches[m.current].Update(msg)
		}
		return m, cmd
	}
	return m, nil
}

func main() {
	m := model{
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("p"),
				key.WithHelp("p", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "q"),
				key.WithHelp("q", "quit"),
			),
			new: key.NewBinding(
				key.WithKeys("n"),
				key.WithHelp("n", "new task"),
			),
			delete: key.NewBinding(
				key.WithKeys("d"),
				key.WithHelp("d", "delete task"),
			),
			next: key.NewBinding(
				key.WithKeys("down"),
				key.WithHelp("down", "next task"),
			),
			prev: key.NewBinding(
				key.WithKeys("up"),
				key.WithHelp("up", "previous task"),
			),
		},
		help: help.New(),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
		os.Exit(1)
	}
}
