package cmd

import (
	"fmt"
	"os"

	"github.com/Carpe-Wang/terminalCrypto/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile      string
	exchangeName string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "terminalcrypto",
	Short: "A terminal-based cryptocurrency price tracker",
	Long: `TerminalCrypto is a CLI application that allows you to track cryptocurrency
prices and market data from various exchanges (Binance, Coinbase, OKX) directly
from your terminal.

Features:
  - Real-time price tracking
  - Detailed market data (24h high/low, volume, price changes)
  - Support for multiple exchanges
  - Secure API credential storage in system keyring
  - Beautiful terminal UI`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config
		if err := config.InitConfig(); err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}

		// If exchange flag is set, use it; otherwise use config default
		if exchangeName == "" {
			exchangeName = config.GetExchange()
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.terminalcrypto/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&exchangeName, "exchange", "e", "", "exchange to use (binance, coinbase, okx)")
}
