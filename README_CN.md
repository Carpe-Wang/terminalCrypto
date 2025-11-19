# TerminalCrypto

一个漂亮的终端加密货币价格追踪工具，支持多个交易所。

[English](README.md) | 简体中文

> 💡 **想要快速开始？** 查看 [超简单使用指南](SIMPLE_GUIDE_CN.md) - 只需输入 `crypto BTC` 就能看价格！

## 功能特性

- 🚀 实时加密货币价格追踪
- 📊 详细的 24 小时市场数据（最高/最低价、交易量、价格变化）
- 🔄 支持多个交易所（Binance、Coinbase、OKX）
- 🔐 安全的 API 凭证存储在系统钥匙串中
- 🎨 美观的彩色终端界面
- ⚡ 快速且轻量级
- 🛡️ 速率限制和错误处理

## 安装

### 从源码安装

```bash
git clone https://github.com/Carpe-Wang/terminalCrypto.git
cd terminalCrypto
go build -o terminalcrypto
sudo mv terminalcrypto /usr/local/bin/  # 可选：全局安装
```

### 前置要求

- Go 1.21 或更高版本
- macOS 或 Linux（Windows 支持即将推出）

## 快速开始

### 1. 设置交易所

配置你喜欢的交易所（凭证是可选的，用于公开数据访问）：

```bash
terminalcrypto setup binance
```

你可以将 API 凭证留空，仅使用公开访问权限，这对于价格查询已经足够。

### 2. 获取当前价格

```bash
# 单个币种
terminalcrypto price BTC

# 多个币种
terminalcrypto price BTC ETH SOL

# 指定交易对
terminalcrypto price BTC/USDT ETH/USDT
```

### 3. 查看详细市场数据

```bash
# 获取 24 小时市场统计
terminalcrypto ticker BTC

# 多个币种
terminalcrypto ticker BTC ETH SOL
```

### 4. 监控实时价格

```bash
# 自动刷新监控价格（默认：5 秒）
terminalcrypto watch BTC ETH SOL

# 自定义刷新间隔
terminalcrypto watch BTC ETH --interval 3
```

按 `q` 退出监控模式。

## 命令说明

### `setup`

配置交易所的 API 凭证。

```bash
terminalcrypto setup [交易所名称]

# 示例：
terminalcrypto setup binance
terminalcrypto setup coinbase
terminalcrypto setup okx
```

凭证会安全地存储在系统钥匙串中：
- **macOS**: Keychain（钥匙串）
- **Linux**: Secret Service（Gnome Keyring、KWallet）
- **Windows**: Credential Manager（即将支持）

### `price`

获取加密货币的当前价格。

```bash
terminalcrypto price [币种...]

# 示例：
terminalcrypto price BTC
terminalcrypto price BTC ETH SOL
terminalcrypto --exchange coinbase price BTC
```

### `ticker`

获取详细的 24 小时市场数据。

```bash
terminalcrypto ticker [币种...]

# 示例：
terminalcrypto ticker BTC
terminalcrypto ticker BTC ETH
```

显示内容：
- 当前价格
- 24 小时价格变化（金额和百分比）
- 24 小时最高/最低价
- 24 小时交易量

### `watch`

实时监控价格并自动刷新。

```bash
terminalcrypto watch [币种...] [选项]

# 选项：
#   -i, --interval int   刷新间隔（秒），默认 5

# 示例：
terminalcrypto watch BTC ETH
terminalcrypto watch BTC ETH SOL --interval 3
```

价格变化通过颜色标识：
- 🟢 绿色：价格上涨
- 🔴 红色：价格下跌
- ⚪ 白色：无变化

## 配置

配置文件存储在 `~/.terminalcrypto/config.yaml`：

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

你可以手动编辑此文件，或使用 `--exchange` 选项覆盖默认交易所。

## 支持的交易所

| 交易所 | 状态 | 公开 API | 认证 API |
|--------|------|----------|----------|
| Binance  | ✅ 可用 | ✅ 是 | ✅ 是 |
| Coinbase | 🚧 即将推出 | - | - |
| OKX      | 🚧 即将推出 | - | - |

## 币种格式

币种可以用多种格式指定：

- `BTC` - 自动添加 USDT（变为 `BTCUSDT`）
- `BTC/USDT` - 斜杠分隔符
- `BTCUSDT` - 无分隔符
- `BTC-USDT` - 短横线分隔符

所有格式都会自动为每个交易所规范化。

## API 凭证

### 为什么需要 API 凭证？

API 凭证是**可选的**。你可以在没有凭证的情况下使用该工具获取公开数据（价格、行情）。

但是，API 凭证可以提供：
- 更高的速率限制
- 访问私有账户数据（未来功能）
- 降低延迟（绕过公开缓存）

### 如何获取 API 凭证

**Binance：**
1. 登录 [Binance](https://www.binance.com)
2. 进入 API 管理
3. 创建新的 API 密钥
4. 保存 API Key 和 Secret Key
5. 运行 `terminalcrypto setup binance` 并输入你的凭证

**Coinbase 和 OKX：** 即将推出

### 安全性

- 凭证存储在系统的安全钥匙串中
- 永远不会提交到 git（已添加到 `.gitignore`）
- 可随时通过系统钥匙串工具删除
- API 密钥仅用于身份验证，不会被记录或显示

## 开发

### 项目结构

```
terminalCrypto/
├── cmd/                    # CLI 命令
│   ├── root.go            # 根命令
│   ├── setup.go           # 设置命令
│   ├── price.go           # 价格命令
│   ├── ticker.go          # 行情命令
│   └── watch.go           # 监控命令
├── internal/
│   ├── config/            # 配置管理
│   ├── exchange/          # 交易所客户端
│   │   ├── exchange.go    # 交易所接口
│   │   └── binance.go     # Binance 实现
│   ├── keyring/           # 凭证存储
│   └── models/            # 数据模型
├── main.go                # 入口文件
├── go.mod                 # Go 模块文件
└── README.md              # 本文件
```

### 添加新交易所

1. 在 `internal/exchange/` 中实现 `Exchange` 接口
2. 在 `exchange.go` 的工厂函数中添加该交易所
3. 更新文档

### 运行测试

```bash
go test ./...
```

### 构建

```bash
# 为当前平台构建
go build -o terminalcrypto

# 为特定平台构建
GOOS=linux GOARCH=amd64 go build -o terminalcrypto-linux
GOOS=darwin GOARCH=amd64 go build -o terminalcrypto-darwin
```

## 常见问题

### "Service unavailable from a restricted location"（服务在受限地区不可用）

某些交易所（如 Binance）可能会限制某些地区的访问。解决方案：
- 使用 VPN
- 尝试其他交易所（Coinbase、OKX）
- 使用特定地区的端点（例如美国用户使用 `binance.us`）

### "Failed to get credentials"（获取凭证失败）

确保你已经先运行了设置命令：
```bash
terminalcrypto setup binance
```

如果你想使用仅公开访问权限，可以忽略此错误（凭证将为空）。

### "Rate limit exceeded"（超出速率限制）

该工具实现了速率限制，但如果仍然达到限制：
- 增加 watch 命令的 `--interval` 间隔
- 减少你追踪的币种数量
- 等待几分钟后重试

## 路线图

- [x] Binance 支持
- [ ] Coinbase 支持
- [ ] OKX 支持
- [ ] 历史价格图表（K线图）
- [ ] 价格提醒
- [ ] 投资组合追踪
- [ ] Windows 支持
- [ ] 配置预设
- [ ] 导出数据为 CSV/JSON
- [ ] watch 模式的 WebSocket 流式传输

## 贡献

欢迎贡献！请随时提交 Pull Request。

1. Fork 本仓库
2. 创建你的功能分支（`git checkout -b feature/amazing-feature`）
3. 提交你的更改（`git commit -m 'Add some amazing feature'`）
4. 推送到分支（`git push origin feature/amazing-feature`）
5. 开启一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 LICENSE 文件。

## 致谢

- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [Viper](https://github.com/spf13/viper) - 配置管理
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - 终端 UI 框架
- [go-binance](https://github.com/adshao/go-binance) - Binance API 客户端
- [go-keyring](https://github.com/zalando/go-keyring) - 安全凭证存储

## 支持

如果你遇到任何问题或有疑问：
- 在 [GitHub](https://github.com/Carpe-Wang/terminalCrypto/issues) 上开启一个 issue
- 查看现有 issues 寻找解决方案

---

使用 ❤️ 制作，作者：[Carpe-Wang](https://github.com/Carpe-Wang)
