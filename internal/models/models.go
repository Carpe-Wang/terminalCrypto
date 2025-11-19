package models

import "time"

// Ticker represents detailed cryptocurrency market data
type Ticker struct {
	Symbol      string    `json:"symbol"`
	Price       float64   `json:"price"`
	Change24h   float64   `json:"change_24h"`
	Volume24h   float64   `json:"volume_24h"`
	High24h     float64   `json:"high_24h"`
	Low24h      float64   `json:"low_24h"`
	LastUpdated time.Time `json:"last_updated"`
}

// Candle represents OHLCV (Open, High, Low, Close, Volume) data
type Candle struct {
	Time   time.Time `json:"time"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume float64   `json:"volume"`
}

// PriceUpdate represents a real-time price update
type PriceUpdate struct {
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}
