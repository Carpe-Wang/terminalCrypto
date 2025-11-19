package exchange

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Carpe-Wang/terminalCrypto/internal/models"
	binance "github.com/adshao/go-binance/v2"
	"golang.org/x/time/rate"
)

// BinanceClient implements the Exchange interface for Binance
type BinanceClient struct {
	client  *binance.Client
	limiter *rate.Limiter
	name    string
}

// NewBinanceClient creates a new Binance client
func NewBinanceClient(apiKey, apiSecret string) (*BinanceClient, error) {
	// For public endpoints, we can use empty credentials
	client := binance.NewClient(apiKey, apiSecret)

	// Rate limit: 10 requests per second (stay well below Binance's 6000 weight/minute)
	limiter := rate.NewLimiter(rate.Limit(10), 10)

	return &BinanceClient{
		client:  client,
		limiter: limiter,
		name:    "binance",
	}, nil
}

// GetName returns the exchange name
func (b *BinanceClient) GetName() string {
	return b.name
}

// NormalizeSymbol converts a symbol like "BTC/USDT" to "BTCUSDT" for Binance
func (b *BinanceClient) NormalizeSymbol(symbol string) string {
	// Remove common separators and convert to uppercase
	normalized := strings.ToUpper(symbol)
	normalized = strings.ReplaceAll(normalized, "/", "")
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, "_", "")

	// If no quote currency specified, default to USDT
	if len(normalized) <= 4 && !strings.HasSuffix(normalized, "USDT") {
		normalized = normalized + "USDT"
	}

	return normalized
}

// GetPrice returns the current price for a symbol
func (b *BinanceClient) GetPrice(ctx context.Context, symbol string) (float64, error) {
	// Rate limiting
	if err := b.limiter.Wait(ctx); err != nil {
		return 0, fmt.Errorf("rate limit error: %w", err)
	}

	normalizedSymbol := b.NormalizeSymbol(symbol)

	prices, err := b.client.NewListPricesService().Symbol(normalizedSymbol).Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get price from Binance: %w", err)
	}

	if len(prices) == 0 {
		return 0, fmt.Errorf("no price data returned for symbol: %s", normalizedSymbol)
	}

	price, err := strconv.ParseFloat(prices[0].Price, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}

	return price, nil
}

// GetTicker returns detailed market data for a symbol
func (b *BinanceClient) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	// Rate limiting
	if err := b.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	normalizedSymbol := b.NormalizeSymbol(symbol)

	ticker, err := b.client.NewListPriceChangeStatsService().Symbol(normalizedSymbol).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticker from Binance: %w", err)
	}

	if len(ticker) == 0 {
		return nil, fmt.Errorf("no ticker data returned for symbol: %s", normalizedSymbol)
	}

	t := ticker[0]

	price, _ := strconv.ParseFloat(t.LastPrice, 64)
	priceChange, _ := strconv.ParseFloat(t.PriceChange, 64)
	volume, _ := strconv.ParseFloat(t.Volume, 64)
	high, _ := strconv.ParseFloat(t.HighPrice, 64)
	low, _ := strconv.ParseFloat(t.LowPrice, 64)

	return &models.Ticker{
		Symbol:      normalizedSymbol,
		Price:       price,
		Change24h:   priceChange,
		Volume24h:   volume,
		High24h:     high,
		Low24h:      low,
		LastUpdated: time.Now(),
	}, nil
}

// GetCandles returns historical OHLCV data
func (b *BinanceClient) GetCandles(ctx context.Context, symbol, interval string, limit int) ([]models.Candle, error) {
	// Rate limiting
	if err := b.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	normalizedSymbol := b.NormalizeSymbol(symbol)

	klines, err := b.client.NewKlinesService().
		Symbol(normalizedSymbol).
		Interval(interval).
		Limit(limit).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get candles from Binance: %w", err)
	}

	candles := make([]models.Candle, len(klines))
	for i, k := range klines {
		open, _ := strconv.ParseFloat(k.Open, 64)
		high, _ := strconv.ParseFloat(k.High, 64)
		low, _ := strconv.ParseFloat(k.Low, 64)
		closePrice, _ := strconv.ParseFloat(k.Close, 64)
		volume, _ := strconv.ParseFloat(k.Volume, 64)

		candles[i] = models.Candle{
			Time:   time.Unix(k.OpenTime/1000, 0),
			Open:   open,
			High:   high,
			Low:    low,
			Close:  closePrice,
			Volume: volume,
		}
	}

	return candles, nil
}
