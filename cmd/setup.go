package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/Carpe-Wang/terminalCrypto/internal/config"
	"github.com/Carpe-Wang/terminalCrypto/internal/keyring"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var setupCmd = &cobra.Command{
	Use:   "setup [exchange]",
	Short: "Configure API credentials for an exchange",
	Long: `Setup allows you to configure API credentials for a cryptocurrency exchange.
Credentials are securely stored in your system's keyring (Keychain on macOS,
Credential Manager on Windows, Secret Service on Linux).

Supported exchanges: binance, coinbase, okx

Example:
  terminalcrypto setup binance

Note: For public data access (prices only), you can leave API credentials empty.
However, some endpoints may require authentication for higher rate limits.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		exchangeName := strings.ToLower(args[0])

		// Validate exchange name
		validExchanges := []string{"binance", "coinbase", "okx"}
		valid := false
		for _, ex := range validExchanges {
			if exchangeName == ex {
				valid = true
				break
			}
		}

		if !valid {
			return fmt.Errorf("invalid exchange: %s (valid options: binance, coinbase, okx)", exchangeName)
		}

		fmt.Printf("Setting up %s\n\n", exchangeName)
		fmt.Println("Enter your API credentials (leave empty for public-only access):")

		// Read API key
		fmt.Print("API Key: ")
		reader := bufio.NewReader(os.Stdin)
		apiKey, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read API key: %w", err)
		}
		apiKey = strings.TrimSpace(apiKey)

		// Read API secret (hidden)
		fmt.Print("API Secret: ")
		apiSecretBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read API secret: %w", err)
		}
		fmt.Println() // New line after password input
		apiSecret := strings.TrimSpace(string(apiSecretBytes))

		// Store credentials if provided
		if apiKey != "" || apiSecret != "" {
			if err := keyring.StoreCredentials(exchangeName, apiKey, apiSecret); err != nil {
				return fmt.Errorf("failed to store credentials: %w", err)
			}
			fmt.Println("Credentials stored securely in system keyring")
		} else {
			fmt.Println("No credentials provided. Using public-only access.")
		}

		// Set as default exchange
		if err := config.SetExchange(exchangeName); err != nil {
			return fmt.Errorf("failed to update config: %w", err)
		}

		fmt.Printf("\n%s is now your default exchange\n", exchangeName)
		fmt.Println("\nYou can now use commands like:")
		fmt.Printf("  terminalcrypto price BTC\n")
		fmt.Printf("  terminalcrypto ticker ETH\n")
		fmt.Printf("  terminalcrypto watch BTC ETH\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
