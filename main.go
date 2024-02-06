package main

import (
	"data_generation/cmd"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the application state
type Model struct {
	Choices        []string
	Selected       map[int]struct{}
	cursorPosition int
	NumRows        int
	OutputFilename string
}

func (m Model) Init() tea.Cmd {
	// Select all columns by default
	for i := range m.Choices {
		m.Selected[i] = struct{}{}
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j", "down":
			if m.cursorPosition < len(m.Choices)-1 {
				m.cursorPosition++
			}
		case "k", "up":
			if m.cursorPosition > 0 {
				m.cursorPosition--
			}
		case "enter":
			switch m.cursorPosition {
			case 0: // Welcome message
				m.cursorPosition++
			case 1: // Column selection
				m.toggleSelection() // Toggle selection of the current choice
			case 2: // Number of rows input
				// Here should be logic to handle the input for the number of rows
				// For simplicity, let's assume input is received elsewhere and stored in m.NumRows
				m.NumRows = 10000
				m.cursorPosition++
			case 3: // Output filename input
				// Here should be logic to handle the input for the output filename
				// For simplicity, let's assume input is received elsewhere and stored in m.OutputFilename
				m.OutputFilename = "output.csv"
				m.cursorPosition++
				// Convert selected map keys to a slice of strings
				selected := make([]string, 0, len(m.Selected))
				for key := range m.Selected {
					selected = append(selected, m.Choices[key])
				}
				// Call GenerateData function with proper arguments
				cmd.GenerateData(m.NumRows, m.OutputFilename, selected)
				return m, tea.Quit // Exit after generating data
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch m.cursorPosition {
	case 0: // Welcome message
		return "Welcome to Data Generation Tool!\n\nIt has been created for the purpose of playing around with fake datasets\nThe main advantage of it is the fact that you can create huge sizes datasets quite quickly.\n\n[Press Enter to continue]"
	case 1: // Column selection
		var view strings.Builder
		view.WriteString("Select columns (use ↑ and ↓ to navigate, press Enter to select):\n\n")
		for i, choice := range m.Choices {
			cursor := " " // no cursor
			if m.cursorPosition == i+1 {
				cursor = ">" // cursor!
			}
			selected := " " // not selected
			if _, ok := m.Selected[i]; ok {
				selected = "x" // selected!
			}
			if i == m.cursorPosition-1 {
				view.WriteString(fmt.Sprintf("[%s] %s\n", selected, choice)) // Highlight the cursor position
			} else {
				view.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, selected, choice))
			}
		}
		view.WriteString("[Press Enter to continue]")
		return view.String()
	case 2: // Number of rows input
		return "Enter the number of rows (press Enter when done):\n[Press Enter to continue]"
	case 3: // Output filename input
		return "Enter the output file name (press Enter when done):\n[Press Enter to continue]"
	default:
		return ""
	}
}

// toggleSelection toggles the selection of the current choice
func (m *Model) toggleSelection() {
	_, ok := m.Selected[m.cursorPosition-1]
	if ok {
		delete(m.Selected, m.cursorPosition-1)
	} else {
		m.Selected[m.cursorPosition-1] = struct{}{}
	}
}

func main() {
	initialModel := Model{
		Choices:        []string{"ID", "Timestamp", "ProductName", "Company", "Price", "Quantity", "Discount", "TotalPrice", "CustomerID", "FirstName", "LastName", "Email", "Address", "City", "State", "Zip", "Country"},
		Selected:       make(map[int]struct{}),
		cursorPosition: 0,
	}

	// Start the Bubble Tea program
	p := tea.NewProgram(initialModel)
	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
