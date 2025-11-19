# TerminalCrypto

A beautiful terminal-based cryptocurrency price tracker with support for multiple exchanges.

## Features

- 🚀 Real-time cryptocurrency price tracking
- 📊 Detailed 24h market data (high/low, volume, price changes)
- 🔄 Multiple exchange support (Binance, Coinbase, OKX)
- 🔐 Secure API credential storage in system keyring
- 🎨 Beautiful colored terminal UI
- ⚡ Fast and lightweight
- 🛡️ Rate limiting and error handling

## Installation

### From Source

```bash
git clone https://github.com/Carpe-Wang/terminalCrypto.git
cd terminalCrypto
go build -o terminalcrypto
sudo mv terminalcrypto /usr/local/bin/  # Optional: install globally
```

### Prerequisites

- Go 1.21 or higher
- macOS or Linux (Windows support coming soon)

## Quick Start

### 1. Setup an Exchange

Configure your preferred exchange (credentials optional for public data):

```bash
terminalcrypto setup binance
```

You can leave API credentials empty for public-only access, which is sufficient for price queries.

### 2. Get Current Prices

```bash
# Single symbol
terminalcrypto price BTC

# Multiple symbols
terminalcrypto price BTC ETH SOL

# Specify trading pair
terminalcrypto price BTC/USDT ETH/USDT
```

### 3. View Detailed Market Data

```bash
# Get 24h market statistics
terminalcrypto ticker BTC

# Multiple symbols
terminalcrypto ticker BTC ETH SOL
```

### 4. Watch Real-time Prices

```bash
# Watch prices with auto-refresh (default: 5 seconds)
terminalcrypto watch BTC ETH SOL

# Custom refresh interval
terminalcrypto watch BTC ETH --interval 3
```

Press `q` to quit the watch mode.

## Commands

### `setup`

Configure API credentials for an exchange.

```bash
terminalcrypto setup [exchange]

# Examples:
terminalcrypto setup binance
terminalcrypto setup coinbase
terminalcrypto setup okx
```

Credentials are stored securely in your system keyring:
- **macOS**: Keychain
- **Linux**: Secret Service (Gnome Keyring, KWallet)
- **Windows**: Credential Manager (coming soon)

### `price`

Get current prices for cryptocurrency symbols.

```bash
terminalcrypto price [symbols...]

# Examples:
terminalcrypto price BTC
terminalcrypto price BTC ETH SOL
terminalcrypto --exchange coinbase price BTC
```

### `ticker`

Get detailed 24-hour market data.

```bash
terminalcrypto ticker [symbols...]

# Examples:
terminalcrypto ticker BTC
terminalcrypto ticker BTC ETH
```

Shows:
- Current price
- 24h price change (amount and percentage)
- 24h high/low
- 24h trading volume

### `watch`

Watch real-time prices with auto-refresh.

```bash
terminalcrypto watch [symbols...] [flags]

# Flags:
#   -i, --interval int   refresh interval in seconds (default 5)

# Examples:
terminalcrypto watch BTC ETH
terminalcrypto watch BTC ETH SOL --interval 3
```

Price changes are color-coded:
- 🟢 Green: Price increased
- 🔴 Red: Price decreased
- ⚪ White: No change

## Configuration

Configuration file is stored at `~/.terminalcrypto/config.yaml`:

```yaml
exchange: binance
exchanges:
  binance: true
  coinbase: false
  okx: false
refresh_interval: 5
display:
  currency: USDT
  decimal_places: 2
```

You can manually edit this file or use the `--exchange` flag to override the default exchange.

## Supported Exchanges

| Exchange | Status | Public API | Authenticated API |
|----------|--------|------------|-------------------|
| Binance  | ✅ Ready | ✅ Yes | ✅ Yes |
| Coinbase | 🚧 Coming Soon | - | - |
| OKX      | 🚧 Coming Soon | - | - |

## Symbol Format

Symbols can be specified in various formats:

- `BTC` - Automatically appends USDT (becomes `BTCUSDT`)
- `BTC/USDT` - Slash separator
- `BTCUSDT` - No separator
- `BTC-USDT` - Dash separator

All formats are automatically normalized for each exchange.

## API Credentials

### Why do I need API credentials?

API credentials are **optional**. You can use the tool without them for public data (prices, tickers).

However, API credentials provide:
- Higher rate limits
- Access to private account data (future feature)
- Reduced latency (bypasses public cache)

### How to get API credentials

**Binance:**
1. Log in to [Binance](https://www.binance.com)
2. Go to API Management
3. Create a new API key
4. Save the API Key and Secret Key
5. Run `terminalcrypto setup binance` and enter your credentials

**Coinbase & OKX:** Coming soon

### Security

- Credentials are stored in your system's secure keyring
- Never committed to git (added to `.gitignore`)
- Can be deleted anytime with system keyring tools
- API keys are only used for authentication, never logged or displayed

## Development

### Project Structure

```
terminalCrypto/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command
│   ├── setup.go           # Setup command
│   ├── price.go           # Price command
│   ├── ticker.go          # Ticker command
│   └── watch.go           # Watch command
├── internal/
│   ├── config/            # Configuration management
│   ├── exchange/          # Exchange clients
│   │   ├── exchange.go    # Exchange interface
│   │   └── binance.go     # Binance implementation
│   ├── keyring/           # Credential storage
│   └── models/            # Data models
├── main.go                # Entry point
├── go.mod                 # Go module file
└── README.md              # This file
```

### Adding a New Exchange

1. Implement the `Exchange` interface in `internal/exchange/`
2. Add the exchange to the factory in `exchange.go`
3. Update documentation

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Build for current platform
go build -o terminalcrypto

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o terminalcrypto-linux
GOOS=darwin GOARCH=amd64 go build -o terminalcrypto-darwin
```

## Troubleshooting

### "Service unavailable from a restricted location"

Some exchanges (like Binance) may restrict access from certain regions. Solutions:
- Use a VPN
- Try a different exchange (Coinbase, OKX)
- Use exchange-specific endpoints (e.g., `binance.us` for US users)

### "Failed to get credentials"

Make sure you've run the setup command first:
```bash
terminalcrypto setup binance
```

If you want to use public-only access, the error can be ignored (credentials will be empty).

### "Rate limit exceeded"

The tool implements rate limiting, but if you still hit limits:
- Increase the `--interval` for watch command
- Reduce the number of symbols you're tracking
- Wait a few minutes before retrying

## Roadmap

- [x] Binance support
- [ ] Coinbase support
- [ ] OKX support
- [ ] Historical price charts (candlestick)
- [ ] Price alerts
- [ ] Portfolio tracking
- [ ] Windows support
- [ ] Configuration presets
- [ ] Export data to CSV/JSON
- [ ] WebSocket streaming for watch mode

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [go-binance](https://github.com/adshao/go-binance) - Binance API client
- [go-keyring](https://github.com/zalando/go-keyring) - Secure credential storage

## Support

If you encounter any issues or have questions:
- Open an issue on [GitHub](https://github.com/Carpe-Wang/terminalCrypto/issues)
- Check existing issues for solutions

---

Made with ❤️ by [Carpe-Wang](https://github.com/Carpe-Wang)
