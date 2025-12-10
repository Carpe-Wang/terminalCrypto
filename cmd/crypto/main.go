package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Carpe-Wang/terminalCrypto/internal/config"
	"github.com/Carpe-Wang/terminalCrypto/internal/exchange"
	"github.com/Carpe-Wang/terminalCrypto/internal/keyring"
	"github.com/Carpe-Wang/terminalCrypto/internal/models"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	// 获取命令名称（btc, eth, sol 等）
	cmdName := strings.ToLower(filepath.Base(os.Args[0]))
	symbol := strings.ToUpper(cmdName)

	// 初始化配置
	if err := config.InitConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "配置初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 获取交易所名称
	exchangeName := config.GetExchange()

	// 获取凭据
	var apiKey, apiSecret string
	creds, err := keyring.GetCredentials(exchangeName)
	if err == nil {
		apiKey = creds.APIKey
		apiSecret = creds.APISecret
	}

	// 创建交易所客户端
	client, err := exchange.Factory(exchangeName, apiKey, apiSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建交易所客户端失败: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// 定义样式
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		Background(lipgloss.Color("#1a1a2e")).
		Padding(0, 1)

	symbolStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D4FF"))

	priceStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF"))

	bigPriceStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF"))

	positiveStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF87"))

	negativeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6B6B"))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(1, 2)

	// 获取 Ticker 数据
	ticker, err := client.GetTicker(ctx, symbol)
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取 %s 数据失败: %v\n", symbol, err)
		os.Exit(1)
	}

	// 计算涨跌幅
	changePercent := 0.0
	prevPrice := ticker.Price - ticker.Change24h
	if prevPrice != 0 {
		changePercent = (ticker.Change24h / prevPrice) * 100
	}

	// 构建输出
	var sb strings.Builder

	// 标题
	sb.WriteString(titleStyle.Render(fmt.Sprintf(" %s 实时行情 ", symbol)))
	sb.WriteString("\n\n")

	// 当前价格（大字体效果）
	sb.WriteString(symbolStyle.Render(ticker.Symbol))
	sb.WriteString("\n")

	// 根据价格大小调整显示格式
	priceStr := formatPrice(ticker.Price)
	sb.WriteString(bigPriceStyle.Render(priceStr))

	// 涨跌指示
	if ticker.Change24h >= 0 {
		sb.WriteString(positiveStyle.Render(fmt.Sprintf("  ▲ +$%.2f (+%.2f%%)", ticker.Change24h, changePercent)))
	} else {
		sb.WriteString(negativeStyle.Render(fmt.Sprintf("  ▼ $%.2f (%.2f%%)", ticker.Change24h, changePercent)))
	}
	sb.WriteString("\n\n")

	// 24小时数据
	sb.WriteString(labelStyle.Render("━━━ 24小时数据 ━━━\n"))
	sb.WriteString(labelStyle.Render("最高: ") + priceStyle.Render(formatPrice(ticker.High24h)) + "  ")
	sb.WriteString(labelStyle.Render("最低: ") + priceStyle.Render(formatPrice(ticker.Low24h)) + "\n")

	// 计算振幅
	if ticker.Low24h > 0 {
		amplitude := ((ticker.High24h - ticker.Low24h) / ticker.Low24h) * 100
		sb.WriteString(labelStyle.Render("振幅: ") + priceStyle.Render(fmt.Sprintf("%.2f%%", amplitude)) + "\n")
	}

	// 成交量（如果有）
	if ticker.Volume24h > 0 {
		sb.WriteString(labelStyle.Render("成交量: ") + priceStyle.Render(formatVolume(ticker.Volume24h)) + "\n")
	}

	// 获取 K 线数据绘制走势图
	candles, err := client.GetCandles(ctx, symbol, "1h", 24)
	if err == nil && len(candles) > 0 {
		sb.WriteString("\n")
		sb.WriteString(labelStyle.Render("━━━ 24小时走势 ━━━\n"))
		sb.WriteString(renderMiniChart(candles))
	}

	// 底部信息
	sb.WriteString("\n")
	sb.WriteString(labelStyle.Render(fmt.Sprintf("来源: %s | %s",
		strings.ToUpper(client.GetName()),
		time.Now().Format("2006-01-02 15:04:05"))))

	// 输出
	fmt.Println(boxStyle.Render(sb.String()))
}

// formatPrice 根据价格大小智能格式化
func formatPrice(price float64) string {
	if price >= 10000 {
		return fmt.Sprintf("$%.2f", price)
	} else if price >= 100 {
		return fmt.Sprintf("$%.2f", price)
	} else if price >= 1 {
		return fmt.Sprintf("$%.4f", price)
	} else if price >= 0.01 {
		return fmt.Sprintf("$%.6f", price)
	}
	return fmt.Sprintf("$%.8f", price)
}

// formatVolume 格式化成交量
func formatVolume(vol float64) string {
	if vol >= 1000000000 {
		return fmt.Sprintf("%.2fB", vol/1000000000)
	} else if vol >= 1000000 {
		return fmt.Sprintf("%.2fM", vol/1000000)
	} else if vol >= 1000 {
		return fmt.Sprintf("%.2fK", vol/1000)
	}
	return fmt.Sprintf("%.2f", vol)
}

// renderMiniChart 渲染简单的终端走势图
func renderMiniChart(candles []models.Candle) string {
	if len(candles) == 0 {
		return ""
	}

	// 找出最高和最低价
	minPrice := candles[0].Low
	maxPrice := candles[0].High
	for _, c := range candles {
		if c.Low < minPrice {
			minPrice = c.Low
		}
		if c.High > maxPrice {
			maxPrice = c.High
		}
	}

	priceRange := maxPrice - minPrice
	if priceRange == 0 {
		priceRange = 1
	}

	// 图表高度
	height := 8
	width := len(candles)

	// 创建图表矩阵
	chart := make([][]rune, height)
	for i := range chart {
		chart[i] = make([]rune, width)
		for j := range chart[i] {
			chart[i][j] = ' '
		}
	}

	// 使用收盘价绘制
	for x, candle := range candles {
		// 计算 y 位置 (0 = 顶部, height-1 = 底部)
		normalizedPrice := (candle.Close - minPrice) / priceRange
		y := height - 1 - int(normalizedPrice*float64(height-1))
		if y < 0 {
			y = 0
		}
		if y >= height {
			y = height - 1
		}

		// 根据涨跌选择字符
		if candle.Close >= candle.Open {
			chart[y][x] = '█'
		} else {
			chart[y][x] = '▄'
		}
	}

	// 转换为字符串
	var result strings.Builder

	// 添加涨跌颜色
	upStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87"))
	downStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			char := chart[y][x]
			if char != ' ' {
				if candles[x].Close >= candles[x].Open {
					result.WriteString(upStyle.Render(string(char)))
				} else {
					result.WriteString(downStyle.Render(string(char)))
				}
			} else {
				result.WriteRune(' ')
			}
		}
		result.WriteRune('\n')
	}

	return result.String()
}