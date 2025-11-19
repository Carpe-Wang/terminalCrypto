package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/Carpe-Wang/terminalCrypto/internal/exchange"
	"github.com/Carpe-Wang/terminalCrypto/internal/keyring"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var tickerCmd = &cobra.Command{
	Use:   "ticker [symbols...]",
	Short: "Get detailed market data for cryptocurrency symbols",
	Long: `Get detailed 24-hour market data including price, volume, high/low, and price change.

Examples:
  terminalcrypto ticker BTC
  terminalcrypto ticker BTC ETH SOL
  terminalcrypto ticker BTC/USDT
  terminalcrypto --exchange binance ticker BTC`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

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

		// Define styles
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFF00"))

		symbolStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D4FF"))

		labelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

		positiveStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF87"))

		negativeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0087"))

		valueStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

		// Print header
		fmt.Println(headerStyle.Render(fmt.Sprintf("\n24h Market Data from %s:", strings.ToUpper(client.GetName()))))
		fmt.Println(strings.Repeat("═", 60))

		// Fetch and display ticker data
		for i, symbol := range args {
			if i > 0 {
				fmt.Println(strings.Repeat("─", 60))
			}

			ticker, err := client.GetTicker(ctx, symbol)
			if err != nil {
				fmt.Printf("\n%s: Error: %v\n", symbol, err)
				continue
			}

			// Display ticker information
			fmt.Printf("\n%s\n", symbolStyle.Render(ticker.Symbol))

			// Price
			fmt.Printf("  %s  %s\n",
				labelStyle.Render("Price:      "),
				valueStyle.Render(fmt.Sprintf("$%.2f", ticker.Price)))

			// 24h Change
			changePercent := (ticker.Change24h / (ticker.Price - ticker.Change24h)) * 100
			var changeStr string
			if ticker.Change24h >= 0 {
				changeStr = positiveStyle.Render(fmt.Sprintf("+$%.2f (+%.2f%%)", ticker.Change24h, changePercent))
			} else {
				changeStr = negativeStyle.Render(fmt.Sprintf("$%.2f (%.2f%%)", ticker.Change24h, changePercent))
			}
			fmt.Printf("  %s  %s\n",
				labelStyle.Render("24h Change:"),
				changeStr)

			// 24h High
			fmt.Printf("  %s  %s\n",
				labelStyle.Render("24h High:  "),
				valueStyle.Render(fmt.Sprintf("$%.2f", ticker.High24h)))

			// 24h Low
			fmt.Printf("  %s  %s\n",
				labelStyle.Render("24h Low:   "),
				valueStyle.Render(fmt.Sprintf("$%.2f", ticker.Low24h)))

			// 24h Volume
			fmt.Printf("  %s  %s\n",
				labelStyle.Render("24h Volume:"),
				valueStyle.Render(fmt.Sprintf("%.2f", ticker.Volume24h)))
		}

		fmt.Println(strings.Repeat("═", 60))
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tickerCmd)
}
