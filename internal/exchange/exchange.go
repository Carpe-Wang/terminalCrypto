package exchange

import (
	"context"
	"fmt"

	"github.com/Carpe-Wang/terminalCrypto/internal/models"
)

// Exchange defines the interface for cryptocurrency exchange clients
type Exchange interface {
	// GetPrice returns the current price for a symbol
	GetPrice(ctx context.Context, symbol string) (float64, error)

	// GetTicker returns detailed market data for a symbol
	GetTicker(ctx context.Context, symbol string) (*models.Ticker, error)

	// GetCandles returns historical OHLCV data
	GetCandles(ctx context.Context, symbol, interval string, limit int) ([]models.Candle, error)

	// NormalizeSymbol converts a generic symbol (e.g., "BTC/USDT") to exchange-specific format
	NormalizeSymbol(symbol string) string

	// GetName returns the exchange name
	GetName() string
}

// Factory creates an exchange client based on the exchange name
func Factory(exchangeName, apiKey, apiSecret string) (Exchange, error) {
	switch exchangeName {
	case "binance":
		return NewBinanceClient(apiKey, apiSecret)
	case "coinbase":
		return nil, fmt.Errorf("coinbase support coming soon")
	case "okx":
		return nil, fmt.Errorf("okx support coming soon")
	default:
		return nil, fmt.Errorf("unsupported exchange: %s (supported: binance, coinbase, okx)", exchangeName)
	}
}
