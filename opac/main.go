package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var welcome = `
~~~~ Welcome to OPAC 2022 ~~~~

Results provided by Fatcat | Perpetual Access to Millions of Open Research
Publications From Around The World | https://fatcat.wiki/
`

var start = &Model{
	Results: []string{},
}

// Model is a basic query and response interaction.
type Model struct {
	Query   textinput.Model
	Results []string
}

func initialModel() *Model {
	ti := textinput.New()
	ti.Placeholder = "text user interfaces"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40
	return &Model{
		Query: ti,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			vs := url.Values{}
			vs.Set("q", m.Query.Value())
			resp, err := http.Get(fmt.Sprintf("https://search.fatcat.wiki/fatcat_release/_search?%s", vs.Encode()))
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			var rr ReleaseResponse
			if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
				log.Fatal(err)
			}
			m.Results = rr.Summary()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	m.Query, cmd = m.Query.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n\n%s",
		welcome,
		m.Query.View(),
		strings.Join(m.Results, "\n"),
	) + "\n"
}

type ReleaseResponse struct {
	Hits struct {
		Hits []struct {
			Source struct {
				Title string `json:"title"`
				DOI   string `json:"doi"`
				Ident string `json:"ident"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (r *ReleaseResponse) Summary() (result []string) {
	var (
		styleTitle = lipgloss.NewStyle().Bold(true)
		styleDim   = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))
		styleDOI   = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	)
	for _, h := range r.Hits.Hits {
		v := strings.TrimSpace(h.Source.Title)
		if len(v) < 5 {
			continue
		}
		result = append(result, fmt.Sprintf("%s [%s]\n    %s",
			styleTitle.Render(h.Source.Title),
			styleDOI.Render(h.Source.DOI),
			styleDim.Render("https://fatcat.wiki/release/"+h.Source.Ident)))
	}
	return
}

func (r *ReleaseResponse) Titles() (result []string) {
	for _, h := range r.Hits.Hits {
		v := strings.TrimSpace(h.Source.Title)
		if len(v) < 5 {
			continue
		}
		result = append(result, h.Source.Title)
	}
	return
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		log.Fatalf("could not start: %v", err)
	}
}
