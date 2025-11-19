package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Carpe-Wang/terminalCrypto/internal/exchange"
	"github.com/Carpe-Wang/terminalCrypto/internal/keyring"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	refreshInterval int
)

type priceData struct {
	symbol    string
	price     float64
	lastPrice float64
	err       error
}

type tickMsg time.Time

type model struct {
	client   exchange.Exchange
	symbols  []string
	prices   map[string]*priceData
	quitting bool
	err      error
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		fetchPrices(m.client, m.symbols),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Duration(refreshInterval)*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchPrices(client exchange.Exchange, symbols []string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		results := make(map[string]*priceData)

		for _, symbol := range symbols {
			price, err := client.GetPrice(ctx, symbol)
			normalizedSymbol := client.NormalizeSymbol(symbol)

			results[normalizedSymbol] = &priceData{
				symbol: normalizedSymbol,
				price:  price,
				err:    err,
			}
		}

		return results
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}

	case tickMsg:
		return m, tea.Batch(
			tickCmd(),
			fetchPrices(m.client, m.symbols),
		)

	case map[string]*priceData:
		// Update prices and track previous values
		for symbol, newData := range msg {
			if oldData, exists := m.prices[symbol]; exists {
				newData.lastPrice = oldData.price
			}
			m.prices[symbol] = newData
		}
		return m, nil

	case error:
		m.err = msg
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	// Define styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFF00")).
		Background(lipgloss.Color("#333333")).
		Padding(0, 1)

	symbolStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D4FF")).
		Width(15)

	priceUpStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF87"))

	priceDownStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF0087"))

	priceNeutralStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF"))

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0087"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)

	// Build the view
	var s strings.Builder

	// Title
	s.WriteString(titleStyle.Render(fmt.Sprintf(" Real-time Prices (%s) ", strings.ToUpper(m.client.GetName()))))
	s.WriteString("\n")
	s.WriteString(strings.Repeat("═", 50))
	s.WriteString("\n\n")

	// Price table
	if len(m.prices) == 0 {
		s.WriteString("Loading prices...\n")
	} else {
		for _, symbol := range m.symbols {
			normalizedSymbol := m.client.NormalizeSymbol(symbol)
			data, exists := m.prices[normalizedSymbol]

			if !exists {
				continue
			}

			if data.err != nil {
				s.WriteString(fmt.Sprintf("%s %s\n",
					symbolStyle.Render(normalizedSymbol),
					errorStyle.Render(fmt.Sprintf("Error: %v", data.err))))
				continue
			}

			// Determine price change indicator
			var priceStyle lipgloss.Style
			var indicator string

			if data.lastPrice > 0 {
				if data.price > data.lastPrice {
					priceStyle = priceUpStyle
					indicator = "↑"
				} else if data.price < data.lastPrice {
					priceStyle = priceDownStyle
					indicator = "↓"
				} else {
					priceStyle = priceNeutralStyle
					indicator = "─"
				}
			} else {
				priceStyle = priceNeutralStyle
				indicator = "─"
			}

			s.WriteString(fmt.Sprintf("%s %s %s\n",
				symbolStyle.Render(normalizedSymbol),
				priceStyle.Render(fmt.Sprintf("$%.2f", data.price)),
				indicator))
		}
	}

	// Footer
	s.WriteString("\n")
	s.WriteString(strings.Repeat("═", 50))
	s.WriteString("\n")
	s.WriteString(helpStyle.Render(fmt.Sprintf("Refreshing every %d seconds • Press 'q' to quit", refreshInterval)))
	s.WriteString("\n")

	return s.String()
}

var watchCmd = &cobra.Command{
	Use:   "watch [symbols...]",
	Short: "Watch real-time prices for cryptocurrency symbols",
	Long: `Watch real-time cryptocurrency prices with auto-refresh.
Prices are color-coded to show increases (green) and decreases (red).

Examples:
  terminalcrypto watch BTC
  terminalcrypto watch BTC ETH SOL
  terminalcrypto watch BTC/USDT ETH/USDT --interval 3
  terminalcrypto --exchange binance watch BTC`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get credentials (may be empty for public access)
		var apiKey, apiSecret string
		creds, err := keyring.GetCredentials(exchangeName)
		if err == nil {
			apiKey = creds.APIKey
			apiSecret = creds.APISecret
		}

		// Create exchange client
		client, err := exchange.Factory(exchangeName, apiKey, apiSecret)
		if err != nil {
			return err
		}

		// Create the model
		m := model{
			client:  client,
			symbols: args,
			prices:  make(map[string]*priceData),
		}

		// Run the Bubble Tea program
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running watch: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().IntVarP(&refreshInterval, "interval", "i", 5, "refresh interval in seconds")
}
