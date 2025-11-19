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

var priceCmd = &cobra.Command{
	Use:   "price [symbols...]",
	Short: "Get current price for cryptocurrency symbols",
	Long: `Get the current price for one or more cryptocurrency symbols.

Examples:
  terminalcrypto price BTC
  terminalcrypto price BTC ETH SOL
  terminalcrypto price BTC/USDT ETH/USDT
  terminalcrypto --exchange binance price BTC`,
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
		symbolStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D4FF"))

		priceStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF87"))

		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0087"))

		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFF00"))

		// Print header
		fmt.Println(headerStyle.Render(fmt.Sprintf("\nPrices from %s:", strings.ToUpper(client.GetName()))))
		fmt.Println(strings.Repeat("─", 50))

		// Fetch and display prices
		for _, symbol := range args {
			price, err := client.GetPrice(ctx, symbol)
			if err != nil {
				fmt.Printf("%s: %s\n",
					symbolStyle.Render(symbol),
					errorStyle.Render(fmt.Sprintf("Error: %v", err)))
				continue
			}

			normalizedSymbol := client.NormalizeSymbol(symbol)
			fmt.Printf("%s: %s\n",
				symbolStyle.Render(normalizedSymbol),
				priceStyle.Render(fmt.Sprintf("$%.2f", price)))
		}

		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(priceCmd)
}
