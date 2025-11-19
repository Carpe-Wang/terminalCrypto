# 超简单使用指南 🚀

只需 2 步，就能在终端查看加密货币价格！

## 第一步：安装

```bash
cd terminalCrypto
./install.sh
```

输入你的密码确认安装即可。

## 第二步：使用

### 查看价格（最常用）

```bash
# 查看 BTC 价格
crypto BTC

# 查看多个币种
crypto BTC ETH SOL
```

就这么简单！🎉

---

## 其他功能（可选）

### 实时监控价格

```bash
crypto watch BTC
```

价格会自动刷新，按 `q` 退出。

### 查看详细数据

```bash
crypto ticker BTC
```

可以看到 24 小时最高/最低价、交易量等。

---

## 常见问题

### Q: 报错 "Service unavailable from a restricted location"

**A:** Binance 在你的地区可能受限，可以使用 VPN。

### Q: 需要 API 密钥吗？

**A:** 不需要！默认使用公开数据就够了。

如果想要更高速率限制，可以运行：
```bash
crypto setup binance
```

### Q: 支持哪些币种？

**A:** 所有 Binance 支持的币种，例如：
- BTC, ETH, BNB, SOL, AVAX, MATIC, DOT, ADA, XRP, DOGE 等

---

## 命令总结

| 命令 | 作用 |
|------|------|
| `crypto BTC` | 查看 BTC 价格 |
| `crypto BTC ETH` | 查看多个币种 |
| `crypto watch BTC` | 实时监控 BTC |
| `crypto ticker BTC` | 查看详细数据 |
| `crypto setup binance` | 配置 API（可选）|

---

就这么简单！现在输入 `crypto BTC` 试试吧！🚀
