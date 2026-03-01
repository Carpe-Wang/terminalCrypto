package i18n

import (
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Lang string

const (
	LangZH Lang = "zh"
	LangEN Lang = "en"
)

var (
	currentLang Lang
	once        sync.Once
)

var zhStrings = map[string]string{
	"config_init_failed":      "配置初始化失败: %v\n",
	"exchange_client_failed":  "创建交易所客户端失败: %v\n",
	"fetch_data_failed":       "获取 %s 数据失败: %v\n",
	"realtime_quote":          " %s 实时行情 ",
	"24h_data_header":         "━━━ 24小时数据 ━━━\n",
	"high":                    "最高: ",
	"low":                     "最低: ",
	"amplitude":               "振幅: ",
	"volume":                  "成交量: ",
	"24h_trend_header":        "━━━ 24小时走势 ━━━\n",
	"source":                  "来源: %s | %s",
}

var enStrings = map[string]string{
	"config_init_failed":      "Config init failed: %v\n",
	"exchange_client_failed":  "Failed to create exchange client: %v\n",
	"fetch_data_failed":       "Failed to fetch %s data: %v\n",
	"realtime_quote":          " %s Live ",
	"24h_data_header":         "━━━ 24h Data ━━━\n",
	"high":                    "High: ",
	"low":                     "Low:  ",
	"amplitude":               "Range: ",
	"volume":                  "Volume: ",
	"24h_trend_header":        "━━━ 24h Trend ━━━\n",
	"source":                  "Source: %s | %s",
}

// DetectLang detects the system language.
// Priority: TERMCRYPTO_LANG env > macOS defaults > LANG env > fallback to English.
func DetectLang() Lang {
	// 1. Check explicit override
	if lang := os.Getenv("TERMCRYPTO_LANG"); lang != "" {
		if strings.HasPrefix(strings.ToLower(lang), "zh") {
			return LangZH
		}
		return LangEN
	}

	// 2. macOS: defaults read -g AppleLocale
	out, err := exec.Command("defaults", "read", "-g", "AppleLocale").Output()
	if err == nil {
		locale := strings.TrimSpace(string(out))
		if strings.HasPrefix(locale, "zh") {
			return LangZH
		}
		return LangEN
	}

	// 3. Fallback: LANG / LC_ALL env
	for _, envKey := range []string{"LC_ALL", "LANG", "LC_MESSAGES"} {
		if val := os.Getenv(envKey); val != "" {
			if strings.HasPrefix(strings.ToLower(val), "zh") {
				return LangZH
			}
			return LangEN
		}
	}

	return LangEN
}

// Init initializes the language. Called automatically on first T() call.
func Init() {
	once.Do(func() {
		currentLang = DetectLang()
	})
}

// GetLang returns the current language.
func GetLang() Lang {
	Init()
	return currentLang
}

// T returns the translated string for the given key.
func T(key string) string {
	Init()
	var m map[string]string
	if currentLang == LangZH {
		m = zhStrings
	} else {
		m = enStrings
	}
	if s, ok := m[key]; ok {
		return s
	}
	return key
}
