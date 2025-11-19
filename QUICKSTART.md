# Quick Start Guide

Get up and running with TerminalCrypto in 5 minutes!

## Installation

```bash
# Clone the repository
git clone https://github.com/Carpe-Wang/terminalCrypto.git
cd terminalCrypto

# Build the project
make build

# (Optional) Install globally
make install
```

## First Steps

### 1. Run Your First Command

Try getting the current Bitcoin price (no setup required):

```bash
# If you installed globally:
terminalcrypto price BTC

# Or if running from build directory:
./bin/terminalcrypto price BTC
```

### 2. Setup Your Exchange (Optional)

For better rate limits and access to more features:

```bash
terminalcrypto setup binance
```

When prompted:
- **API Key**: Enter your Binance API key (or press Enter to skip)
- **API Secret**: Enter your API secret (or press Enter to skip)

**Note**: You can skip credentials for public-only access!

### 3. Try Different Commands

#### Get Current Prices
```bash
# Single cryptocurrency
terminalcrypto price BTC

# Multiple cryptocurrencies
terminalcrypto price BTC ETH SOL DOGE

# Specific trading pair
terminalcrypto price BTC/USDT
```

#### View Detailed Market Data
```bash
# 24-hour statistics for Bitcoin
terminalcrypto ticker BTC

# Multiple coins
terminalcrypto ticker BTC ETH BNB
```

#### Watch Real-time Prices
```bash
# Live price updates (refreshes every 5 seconds)
terminalcrypto watch BTC ETH

# Custom refresh interval (3 seconds)
terminalcrypto watch BTC ETH --interval 3
```

Press `q` to exit watch mode.

## Common Use Cases

### Portfolio Tracking

Watch your favorite cryptocurrencies:
```bash
terminalcrypto watch BTC ETH BNB SOL AVAX MATIC
```

### Quick Price Check

Get instant prices without authentication:
```bash
terminalcrypto price BTC ETH
```

### Market Analysis

View detailed 24h statistics:
```bash
terminalcrypto ticker BTC ETH SOL
```

## Troubleshooting

### "Service unavailable from a restricted location"

Binance may be restricted in your region. Solutions:
1. Use a VPN
2. Wait for Coinbase/OKX support (coming soon)
3. Use Binance.US if you're in the United States

### Command Not Found

If `terminalcrypto` is not found after installation:
```bash
# Make sure /usr/local/bin is in your PATH
echo $PATH

# Or run directly from build directory
./bin/terminalcrypto
```

### No Credentials Setup

You can use the tool without credentials for public data:
- Price queries work without authentication
- Ticker data works without authentication
- Watch mode works without authentication

Credentials are only needed for:
- Higher rate limits
- Private account data (future feature)

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Check out the [configuration options](config.example.yaml)
- Star the project on GitHub if you find it useful!

## Getting Help

- Open an issue: https://github.com/Carpe-Wang/terminalCrypto/issues
- Read the docs: [README.md](README.md)

---

Happy crypto tracking! 🚀
