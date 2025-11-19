#!/bin/bash

# TerminalCrypto 简化安装脚本

set -e

echo "==================================="
echo "TerminalCrypto 简化安装"
echo "==================================="
echo ""

# 构建项目
echo "正在构建项目..."
go build -o terminalcrypto

# 创建简化的别名脚本
echo "正在创建简化命令..."

# 创建 crypto 命令（简化版）
cat > crypto << 'EOF'
#!/bin/bash
# 简化的加密货币价格查询命令

# 如果没有参数，显示帮助
if [ $# -eq 0 ]; then
    echo "用法: crypto <币种> [币种2] [币种3] ..."
    echo ""
    echo "示例:"
    echo "  crypto BTC              # 查看 BTC 价格"
    echo "  crypto BTC ETH SOL      # 查看多个币种价格"
    echo ""
    echo "其他命令:"
    echo "  crypto watch BTC        # 实时监控 BTC"
    echo "  crypto ticker BTC       # 查看 BTC 详细数据"
    exit 0
fi

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# 检查第一个参数是否是特殊命令
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
        # 默认就是查价格
        exec "$SCRIPT_DIR/terminalcrypto" price "$@"
        ;;
esac
EOF

chmod +x crypto

# 安装到系统
echo ""
echo "正在安装到 /usr/local/bin..."
sudo cp terminalcrypto /usr/local/bin/
sudo cp crypto /usr/local/bin/

echo ""
echo "==================================="
echo "✅ 安装成功！"
echo "==================================="
echo ""
echo "现在你可以直接使用："
echo ""
echo "  crypto BTC              # 查看 BTC 价格"
echo "  crypto BTC ETH SOL      # 查看多个币种"
echo "  crypto watch BTC        # 实时监控"
echo "  crypto ticker BTC       # 详细数据"
echo ""
echo "首次使用建议运行（可选）："
echo "  crypto setup binance"
echo ""
echo "祝你使用愉快！🚀"
