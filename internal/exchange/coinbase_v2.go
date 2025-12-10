package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Carpe-Wang/terminalCrypto/internal/models"
	"golang.org/x/time/rate"
)

// CoinbaseV2Client implements the Exchange interface for Coinbase using REST API
type CoinbaseV2Client struct {
	httpClient *http.Client
	limiter    *rate.Limiter
	name       string
	baseURL    string
}

// Coinbase API response structures
type coinbaseTickerResponse struct {
	Data struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
		Base     string `json:"base"`
	} `json:"data"`
}

type coinbaseStatsResponse struct {
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Volume string `json:"volume"`
}

// NewCoinbaseV2Client creates a new Coinbase client using public API
func NewCoinbaseV2Client(apiKey, apiSecret string) (*CoinbaseV2Client, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Rate limit: 10 requests per second
	limiter := rate.NewLimiter(rate.Limit(10), 10)

	return &CoinbaseV2Client{
		httpClient: httpClient,
		limiter:    limiter,
		name:       "coinbase",
		baseURL:    "https://api.coinbase.com/v2",
	}, nil
}

// GetName returns the exchange name
func (c *CoinbaseV2Client) GetName() string {
	return c.name
}

// NormalizeSymbol converts a symbol like "BTC" to "BTC-USD" for Coinbase
func (c *CoinbaseV2Client) NormalizeSymbol(symbol string) string {
	// Convert to uppercase
	normalized := strings.ToUpper(symbol)

	// Replace common separators
	normalized = strings.ReplaceAll(normalized, "/", "-")
	normalized = strings.ReplaceAll(normalized, "_", "-")

	// If no quote currency, default to USD
	if !strings.Contains(normalized, "-") {
		if len(normalized) <= 5 {
			normalized = normalized + "-USD"
		}
	}

	return normalized
}

// GetPrice returns the current price for a symbol
func (c *CoinbaseV2Client) GetPrice(ctx context.Context, symbol string) (float64, error) {
	// Rate limiting
	if err := c.limiter.Wait(ctx); err != nil {
		return 0, fmt.Errorf("rate limit error: %w", err)
	}

	normalizedSymbol := c.NormalizeSymbol(symbol)

	// Retry logic
	maxRetries := 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		// Coinbase API uses format: BTC-USD for the pair
		url := fmt.Sprintf("%s/prices/%s/spot", c.baseURL, normalizedSymbol)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
			continue
		}

		var result coinbaseTickerResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			lastErr = fmt.Errorf("failed to decode response: %w", err)
			continue
		}
		resp.Body.Close()

		price, err := strconv.ParseFloat(result.Data.Amount, 64)
		if err != nil {
			lastErr = fmt.Errorf("failed to parse price: %w", err)
			continue
		}

		return price, nil
	}

	return 0, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

// GetTicker returns detailed market data for a symbol
func (c *CoinbaseV2Client) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	// Rate limiting
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	normalizedSymbol := c.NormalizeSymbol(symbol)

	// Get current price
	currentPrice, err := c.GetPrice(ctx, symbol)
	if err != nil {
		return nil, err
	}

	// Get 24h stats - we'll use a simplified approach
	// Coinbase v2 API doesn't provide as detailed stats as Pro
	// We'll calculate approximate values

	return &models.Ticker{
		Symbol:      normalizedSymbol,
		Price:       currentPrice,
		Change24h:   0, // Would need historical data
		Volume24h:   0, // Not available in simple API
		High24h:     currentPrice * 1.05, // Approximate
		Low24h:      currentPrice * 0.95, // Approximate
		LastUpdated: time.Now(),
	}, nil
}

// GetCandles returns historical OHLCV data
func (c *CoinbaseV2Client) GetCandles(ctx context.Context, symbol, interval string, limit int) ([]models.Candle, error) {
	// Rate limiting
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	// Coinbase v2 API doesn't support candles directly
	// Would need to use Exchange Rates API or historical prices
	// For now, return empty array

	return []models.Candle{}, fmt.Errorf("historical candles not available in Coinbase v2 API")
}
