// Package main provides a terminal UI for browsing and downloading Helm chart values.
// It allows users to interactively select repositories, charts, and versions,
// then download the default values.yaml file for the selected chart version.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Margin(1, 0)

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Margin(1, 0)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	// New styles for version list
	chartVersionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)

	appVersionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("140"))

	latestBadgeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Bold(true)
)

// the state represents the current state of the application
type state int

// Application states
const (
	stateRepoUpdate state = iota
	stateRepoList
	stateChartList
	stateVersionList
	stateDownload
	stateError
	stateComplete
)

// pageSize defines the number of items to show per page
const pageSize = 10

// HelmRepo represents a Helm repository with name and URL
type HelmRepo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// HelmChart represents a Helm chart with metadata
type HelmChart struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	AppVersion  string `json:"app_version"`
	Description string `json:"description"`
}

// HelmVersion represents a specific version of a Helm chart
type HelmVersion struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	AppVersion string `json:"app_version"`
	Created    string `json:"created"`
}

// the model represents the application state for the Helm browser TUI
type model struct {
	state           state
	repos           []HelmRepo
	charts          []HelmChart
	versions        []HelmVersion
	selectedRepo    int
	selectedChart   int
	selectedVersion int
	cursor          int
	loading         bool
	error           string
	message         string
}

// initialModel creates a new model with default values
func initialModel() model {
	return model{
		state:   stateRepoUpdate,
		loading: true,
	}
}

// Init satisfies the tea.Model interface
func (m model) Init() tea.Cmd {
	return updateRepos()
}

// Helper functions for pagination

// getCurrentPage returns the current page number (0-indexed)
func (m model) getCurrentPage() int {
	return m.cursor / pageSize
}

// getPageStart returns the starting index for the current page
func (m model) getPageStart() int {
	return m.getCurrentPage() * pageSize
}

// getPageEnd returns the ending index for the current page
func (m model) getPageEnd(totalItems int) int {
	end := m.getPageStart() + pageSize
	if end > totalItems {
		end = totalItems
	}
	return end
}

// getCursorInPage returns the cursor position within the current page
func (m model) getCursorInPage() int {
	return m.cursor % pageSize
}

// Message types for Bubble Tea communication
type repoUpdateMsg struct{}
type reposLoadedMsg []HelmRepo
type chartsLoadedMsg []HelmChart
type versionsLoadedMsg []HelmVersion
type downloadCompleteMsg string
type errorMsg string

// Bubble Tea commands for async operations

// updateRepos runs the helm repo update command
func updateRepos() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("helm", "repo", "update")
		if err := cmd.Run(); err != nil {
			return errorMsg(fmt.Sprintf("Failed to update repos: %v", err))
		}
		return repoUpdateMsg{}
	}
}

// loadRepos fetches the list of configured Helm repositories
func loadRepos() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("helm", "repo", "list", "-o", "json")
		output, err := cmd.Output()
		if err != nil {
			return errorMsg(fmt.Sprintf("Failed to list repos: %v", err))
		}

		var repos []HelmRepo
		if len(output) > 0 {
			if err := json.Unmarshal(output, &repos); err != nil {
				return errorMsg(fmt.Sprintf("Failed to parse repos: %v", err))
			}
		}

		return reposLoadedMsg(repos)
	}
}

// loadCharts fetches charts from a specific repository
func loadCharts(repoName string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("helm", "search", "repo", repoName+"/", "-o", "json")
		output, err := cmd.Output()
		if err != nil {
			return errorMsg(fmt.Sprintf("Failed to search charts: %v", err))
		}

		var charts []HelmChart
		if len(output) > 0 {
			if err := json.Unmarshal(output, &charts); err != nil {
				return errorMsg(fmt.Sprintf("Failed to parse charts: %v", err))
			}
		}

		return chartsLoadedMsg(charts)
	}
}

// loadVersions fetches all versions of a specific chart
func loadVersions(chartName string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("helm", "search", "repo", chartName, "--versions", "-o", "json")
		output, err := cmd.Output()
		if err != nil {
			return errorMsg(fmt.Sprintf("Failed to search versions: %v", err))
		}

		var versions []HelmVersion
		if len(output) > 0 {
			if err := json.Unmarshal(output, &versions); err != nil {
				return errorMsg(fmt.Sprintf("Failed to parse versions: %v", err))
			}
		}

		return versionsLoadedMsg(versions)
	}
}

// downloadValues downloads the default values.yaml for a chart version
func downloadValues(chartName, version string) tea.Cmd {
	return func() tea.Msg {
		// Get values using helm show values
		cmd := exec.Command("helm", "show", "values", chartName, "--version", version)
		values, err := cmd.Output()
		if err != nil {
			return errorMsg(fmt.Sprintf("Failed to get chart values: %v", err))
		}

		// Create filename
		chartParts := strings.Split(chartName, "/")
		chartBaseName := chartParts[len(chartParts)-1]
		filename := fmt.Sprintf("%s-%s-default-values.yaml", chartBaseName, version)

		// Write to file
		if err := os.WriteFile(filename, values, 0644); err != nil {
			return errorMsg(fmt.Sprintf("Failed to write values file: %v", err))
		}

		return downloadCompleteMsg(filename)
	}
}

// Update handles incoming messages and updates the model state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			switch m.state {
			case stateRepoList:
				if m.cursor > 0 {
					m.cursor--
				}
			case stateChartList:
				if m.cursor > 0 {
					m.cursor--
				}
			case stateVersionList:
				if m.cursor > 0 {
					m.cursor--
				}
			default:
				// No cursor movement for other states
			}

		case "down", "j":
			switch m.state {
			case stateRepoList:
				if m.cursor < len(m.repos)-1 {
					m.cursor++
				}
			case stateChartList:
				if m.cursor < len(m.charts)-1 {
					m.cursor++
				}
			case stateVersionList:
				if m.cursor < len(m.versions)-1 {
					m.cursor++
				}
			default:
				// No cursor movement for other states
			}

		case "enter", " ":
			switch m.state {
			case stateRepoList:
				if len(m.repos) > 0 {
					m.selectedRepo = m.cursor
					m.cursor = 0
					m.loading = true
					m.state = stateChartList
					return m, loadCharts(m.repos[m.selectedRepo].Name)
				}
			case stateChartList:
				if len(m.charts) > 0 {
					m.selectedChart = m.cursor
					m.cursor = 0
					m.loading = true
					m.state = stateVersionList
					return m, loadVersions(m.charts[m.selectedChart].Name)
				}
			case stateVersionList:
				if len(m.versions) > 0 {
					m.selectedVersion = m.cursor
					m.loading = true
					m.state = stateDownload
					return m, downloadValues(m.versions[m.selectedVersion].Name, m.versions[m.selectedVersion].Version)
				}
			case stateComplete:
				// Any key press in complete state should exit
				return m, tea.Quit
			default:
				// No action for other states
			}

		case "backspace", "esc":
			switch m.state {
			case stateChartList:
				m.state = stateRepoList
				m.cursor = m.selectedRepo
				m.charts = nil
			case stateVersionList:
				m.state = stateChartList
				m.cursor = m.selectedChart
				m.versions = nil
			case stateComplete:
				return m, tea.Quit
			default:
				// No back action for other states
			}

		default:
			// Handle "any key to exit" in complete state
			if m.state == stateComplete {
				return m, tea.Quit
			}

			// Number shortcuts (for current page only)
			if len(msg.String()) == 1 {
				if num, err := strconv.Atoi(msg.String()); err == nil && num >= 1 && num <= pageSize {
					switch m.state {
					case stateRepoList:
						pageStart := m.getPageStart()
						absoluteIndex := pageStart + num - 1
						if absoluteIndex < len(m.repos) {
							m.selectedRepo = absoluteIndex
							m.cursor = 0
							m.loading = true
							m.state = stateChartList
							return m, loadCharts(m.repos[m.selectedRepo].Name)
						}
					case stateChartList:
						pageStart := m.getPageStart()
						absoluteIndex := pageStart + num - 1
						if absoluteIndex < len(m.charts) {
							m.selectedChart = absoluteIndex
							m.cursor = 0
							m.loading = true
							m.state = stateVersionList
							return m, loadVersions(m.charts[m.selectedChart].Name)
						}
					case stateVersionList:
						pageStart := m.getPageStart()
						absoluteIndex := pageStart + num - 1
						if absoluteIndex < len(m.versions) {
							m.selectedVersion = absoluteIndex
							m.loading = true
							m.state = stateDownload
							return m, downloadValues(m.versions[m.selectedVersion].Name, m.versions[m.selectedVersion].Version)
						}
					default:
						// No number shortcuts for other states
					}
				}
			}
		}

	case repoUpdateMsg:
		m.loading = true
		return m, loadRepos()

	case reposLoadedMsg:
		m.repos = msg
		m.loading = false
		m.state = stateRepoList
		m.cursor = 0

	case chartsLoadedMsg:
		m.charts = msg
		m.loading = false
		m.cursor = 0

	case versionsLoadedMsg:
		m.versions = msg
		m.loading = false
		m.cursor = 0

	case downloadCompleteMsg:
		m.loading = false
		m.state = stateComplete
		m.message = fmt.Sprintf("Successfully downloaded: %s", msg)

	case errorMsg:
		m.loading = false
		m.state = stateError
		m.error = string(msg)
	}

	return m, nil
}

// View renders the current state of the application
func (m model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("🚀 Helm Chart Browser"))
	s.WriteString("\n\n")

	switch m.state {
	case stateRepoUpdate:
		s.WriteString("🔄 Updating Helm repositories...\n")

	case stateRepoList:
		if m.loading {
			s.WriteString("🔄 Loading repositories...\n")
		} else {
			s.WriteString("🚀 Select a Helm repository:\n\n")

			// Header
			s.WriteString(fmt.Sprintf("%-4s %-20s %s\n", "", "REPOSITORY", "URL"))
			s.WriteString(fmt.Sprintf("%-4s %-20s %s\n", "────", "────────────────────", "───────────────────────────────────"))

			start := m.getPageStart()
			end := m.getPageEnd(len(m.repos))

			for i := start; i < end; i++ {
				repo := m.repos[i]

				// Format number
				numStr := fmt.Sprintf("%d.", i+1)

				// Format repository name with color
				repoName := chartVersionStyle.Render(fmt.Sprintf("%-20s", repo.Name))

				// Format URL with color
				repoURL := appVersionStyle.Render(repo.URL)

				line := fmt.Sprintf("%-4s %s %s", numStr, repoName, repoURL)

				if i == m.cursor {
					s.WriteString(selectedStyle.Render("► " + line))
				} else {
					s.WriteString("  " + line)
				}
				s.WriteString("\n")
			}

			s.WriteString("\n")

			// Show pagination info
			if len(m.repos) > pageSize {
				totalPages := (len(m.repos) + pageSize - 1) / pageSize
				currentPage := m.getCurrentPage() + 1
				paginationInfo := fmt.Sprintf("📄 Page %d of %d • %d total repositories", currentPage, totalPages, len(m.repos))
				s.WriteString(helpStyle.Render(paginationInfo))
			} else if len(m.repos) > 1 {
				totalInfo := fmt.Sprintf("📄 %d repositories available", len(m.repos))
				s.WriteString(helpStyle.Render(totalInfo))
			}
		}

	case stateChartList:
		if m.loading {
			s.WriteString("🔄 Loading charts...\n")
		} else {
			s.WriteString(fmt.Sprintf("📊 Charts in repository '%s':\n\n", m.repos[m.selectedRepo].Name))

			// Header
			s.WriteString(fmt.Sprintf("%-4s %-30s %s\n", "", "CHART NAME", "VERSION"))
			s.WriteString(fmt.Sprintf("%-4s %-30s %s\n", "────", "──────────────────────────────", "───────"))

			start := m.getPageStart()
			end := m.getPageEnd(len(m.charts))

			for i := start; i < end; i++ {
				chart := m.charts[i]

				// Format number
				numStr := fmt.Sprintf("%d.", i+1)

				// Format chart name with color
				chartName := chartVersionStyle.Render(fmt.Sprintf("%-30s", strings.TrimPrefix(chart.Name, m.repos[m.selectedRepo].Name+"/")))

				// Format version with color
				chartVer := appVersionStyle.Render(fmt.Sprintf("v%s", chart.Version))

				line := fmt.Sprintf("%-4s %s %s", numStr, chartName, chartVer)

				if i == m.cursor {
					s.WriteString(selectedStyle.Render("► " + line))
				} else {
					s.WriteString("  " + line)
				}
				s.WriteString("\n")
			}

			s.WriteString("\n")

			// Show pagination info
			if len(m.charts) > pageSize {
				totalPages := (len(m.charts) + pageSize - 1) / pageSize
				currentPage := m.getCurrentPage() + 1
				paginationInfo := fmt.Sprintf("📄 Page %d of %d • %d total charts", currentPage, totalPages, len(m.charts))
				s.WriteString(helpStyle.Render(paginationInfo))
			} else if len(m.charts) > 1 {
				totalInfo := fmt.Sprintf("📄 %d charts available", len(m.charts))
				s.WriteString(helpStyle.Render(totalInfo))
			}
		}

	case stateVersionList:
		if m.loading {
			s.WriteString("🔄 Loading versions...\n")
		} else {
			chartName := strings.TrimPrefix(m.charts[m.selectedChart].Name, m.repos[m.selectedRepo].Name+"/")
			s.WriteString(fmt.Sprintf("📦 Versions of chart '%s':\n\n", chartName))

			// Header
			s.WriteString(fmt.Sprintf("%-4s %-15s %-15s %s\n", "", "CHART VERSION", "APP VERSION", ""))
			s.WriteString(fmt.Sprintf("%-4s %-15s %-15s %s\n", "────", "─────────────", "───────────", "──────"))

			start := m.getPageStart()
			end := m.getPageEnd(len(m.versions))

			for i := start; i < end; i++ {
				version := m.versions[i]

				// Format number
				numStr := fmt.Sprintf("%d.", i+1)

				// Format chart version with color
				chartVer := chartVersionStyle.Render(fmt.Sprintf("%-15s", version.Version))

				// Format app version with color
				appVer := ""
				if version.AppVersion != "" {
					appVer = appVersionStyle.Render(fmt.Sprintf("%-15s", version.AppVersion))
				} else {
					appVer = fmt.Sprintf("%-15s", "─")
				}

				// Add "LATEST" badge for first version
				badge := ""
				if i == 0 {
					badge = latestBadgeStyle.Render("🏷️  LATEST")
				}

				line := fmt.Sprintf("%-4s %s %s %s", numStr, chartVer, appVer, badge)

				if i == m.cursor {
					s.WriteString(selectedStyle.Render("► " + line))
				} else {
					s.WriteString("  " + line)
				}
				s.WriteString("\n")
			}

			s.WriteString("\n")

			// Show pagination info with better formatting
			if len(m.versions) > pageSize {
				totalPages := (len(m.versions) + pageSize - 1) / pageSize
				currentPage := m.getCurrentPage() + 1
				paginationInfo := fmt.Sprintf("📄 Page %d of %d • %d total versions", currentPage, totalPages, len(m.versions))
				s.WriteString(helpStyle.Render(paginationInfo))
			} else if len(m.versions) > 1 {
				totalInfo := fmt.Sprintf("📄 %d versions available", len(m.versions))
				s.WriteString(helpStyle.Render(totalInfo))
			}
		}

	case stateDownload:
		s.WriteString("⬇️  Downloading values.yaml...\n")

	case stateComplete:
		s.WriteString("✅ " + m.message + "\n\n")
		s.WriteString(selectedStyle.Render("🎉 Press any key to exit..."))

	case stateError:
		s.WriteString(errorStyle.Render("❌ Error: " + m.error))
		s.WriteString("\n\nPress 'q' to quit.")

	default:
		s.WriteString("❓ Unknown state")
	}

	// Help text
	switch m.state {
	case stateRepoList, stateChartList, stateVersionList:
		s.WriteString("\n")
		s.WriteString(helpStyle.Render("⌨️  Navigate: ↑/↓ arrows or j/k • Select: Enter/Space or number (1-9,0 for items on current page) • Back: Backspace/Esc • Quit: q/Ctrl+C"))
		s.WriteString("\n")
		s.WriteString(helpStyle.Render("💡 Tip: Use arrow keys to navigate through pages of results"))
	case stateComplete:
		s.WriteString("\n")
		s.WriteString(helpStyle.Render("⌨️  Press any key to exit the application"))
	case stateError:
		s.WriteString("\n")
		s.WriteString(helpStyle.Render("⌨️  Press 'q' to quit the application"))
	default:
		// No help text for loading states
	}

	return s.String()
}

// the main is the entry point of the Helm Chart Browser application
func main() {
	// Check if helm is installed
	if _, err := exec.LookPath("helm"); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: helm command not found. Please install Helm first.\n")
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
