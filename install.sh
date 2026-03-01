#!/bin/bash

set -e

# Detect system language
LOCALE=$(defaults read -g AppleLocale 2>/dev/null || echo "en_US")
if [[ "$TERMCRYPTO_LANG" == zh* ]]; then
    USE_ZH=true
elif [[ "$TERMCRYPTO_LANG" == en* ]]; then
    USE_ZH=false
elif [[ "$LOCALE" == zh* ]]; then
    USE_ZH=true
else
    USE_ZH=false
fi

msg() {
    if $USE_ZH; then
        echo "$1"
    else
        echo "$2"
    fi
}

echo "==================================="
msg "TerminalCrypto 简化安装" "TerminalCrypto Installation"
echo "==================================="
echo ""

# Build
msg "正在构建项目..." "Building project..."
go build -o terminalcrypto

# Create simplified alias script
msg "正在创建简化命令..." "Creating shortcut commands..."

# Create crypto command
cat > crypto << 'CRYPTOEOF'
#!/bin/bash

# Detect language
LOCALE=$(defaults read -g AppleLocale 2>/dev/null || echo "en_US")
if [[ "$TERMCRYPTO_LANG" == zh* ]]; then
    USE_ZH=true
elif [[ "$TERMCRYPTO_LANG" == en* ]]; then
    USE_ZH=false
elif [[ "$LOCALE" == zh* ]]; then
    USE_ZH=true
else
    USE_ZH=false
fi

if [ $# -eq 0 ]; then
    if $USE_ZH; then
        echo "用法: crypto <币种> [币种2] [币种3] ..."
        echo ""
        echo "示例:"
        echo "  crypto BTC              # 查看 BTC 价格"
        echo "  crypto BTC ETH SOL      # 查看多个币种价格"
        echo ""
        echo "其他命令:"
        echo "  crypto watch BTC        # 实时监控 BTC"
        echo "  crypto ticker BTC       # 查看 BTC 详细数据"
    else
        echo "Usage: crypto <symbol> [symbol2] [symbol3] ..."
        echo ""
        echo "Examples:"
        echo "  crypto BTC              # Check BTC price"
        echo "  crypto BTC ETH SOL      # Check multiple prices"
        echo ""
        echo "Other commands:"
        echo "  crypto watch BTC        # Watch BTC in real-time"
        echo "  crypto ticker BTC       # Detailed BTC market data"
    fi
    exit 0
fi

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

case "$1" in
    watch)
        shift
        exec "$SCRIPT_DIR/terminalcrypto" watch "$@"
        ;;
    ticker)
        shift
        exec "$SCRIPT_DIR/terminalcrypto" ticker "$@"
        ;;
    setup)
        shift
        exec "$SCRIPT_DIR/terminalcrypto" setup "$@"
        ;;
    help|--help|-h)
        exec "$SCRIPT_DIR/terminalcrypto" --help
        ;;
    *)
        exec "$SCRIPT_DIR/terminalcrypto" price "$@"
        ;;
esac
CRYPTOEOF

chmod +x crypto

# Install to system
echo ""
msg "正在安装到 /usr/local/bin..." "Installing to /usr/local/bin..."
sudo cp terminalcrypto /usr/local/bin/
sudo cp crypto /usr/local/bin/

echo ""
echo "==================================="
msg "✅ 安装成功！" "✅ Installation successful!"
echo "==================================="
echo ""
msg "现在你可以直接使用：" "You can now use:"
echo ""
msg "  crypto BTC              # 查看 BTC 价格" "  crypto BTC              # Check BTC price"
msg "  crypto BTC ETH SOL      # 查看多个币种" "  crypto BTC ETH SOL      # Check multiple prices"
msg "  crypto watch BTC        # 实时监控" "  crypto watch BTC        # Watch in real-time"
msg "  crypto ticker BTC       # 详细数据" "  crypto ticker BTC       # Detailed market data"
echo ""
msg "首次使用建议运行（可选）：" "Recommended first-time setup (optional):"
echo "  crypto setup binance"
echo ""
msg "祝你使用愉快！🚀" "Enjoy! 🚀"
