# cwa-weather

[![Go Report Card](https://goreportcard.com/badge/github.com/kerkerj/cwa-weather)](https://goreportcard.com/report/github.com/kerkerj/cwa-weather)
[![Go Reference](https://pkg.go.dev/badge/github.com/kerkerj/cwa-weather.svg)](https://pkg.go.dev/github.com/kerkerj/cwa-weather)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Built with Claude](https://img.shields.io/badge/Built%20with-Claude-blueviolet)](https://claude.ai)
[![Skill Scanner](https://img.shields.io/badge/Skill%20Scanner-SAFE-brightgreen)](https://github.com/cisco-ai-defense/skill-scanner)
![coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/kerkerj/6e8ceb0d7091bfa999901383f6486a29/raw/cwa-coverage.json)

[English](README.md)

中央氣象署開放資料 API 的 CLI 工具與 Go 函式庫，用於查詢臺灣天氣資料。

支援鄉鎮預報、即時觀測、36 小時概況、氣象警特報、颱風動態及海象觀測。所有輸出皆為 JSON，適合搭配 agent 或 `jq` 使用。

## 安裝

```bash
go install github.com/kerkerj/cwa-weather/cmd/cwa-weather@latest
```

或從 [GitHub Releases](https://github.com/kerkerj/cwa-weather/releases) 下載執行檔。

## 設定

至 https://opendata.cwa.gov.tw 申請 API 金鑰並匯出：

```bash
export CWA_API_KEY=your-key
```

## 使用方式

### 鄉鎮預報

```bash
# 鄉鎮層級預報
cwa-weather forecast --city 臺北市 --town 中正區

# 縣市層級預報（全部鄉鎮）
cwa-weather forecast --city 台北市    # 台→臺 自動轉換

# 篩選天氣要素
cwa-weather forecast --city 新北市 --town 板橋區 --element 溫度,天氣現象

# 篩選時間區間
cwa-weather forecast --city 臺北市 --time-from 2026-03-01T06:00:00

# 同時篩選要素與時間
cwa-weather forecast --city 臺北市 --element 降雨機率 --time-from 2026-03-01T06:00:00 --time-to 2026-03-01T18:00:00
```

> **提示**：天氣要素名稱由氣象署 API 定義，不加 `--element` 即可在 JSON 回應中查看所有可用名稱。

### 即時觀測

```bash
# 依縣市查詢
cwa-weather observe --city 新北市

# 依測站名稱查詢
cwa-weather observe --station 淡水

# 篩選天氣要素
cwa-weather observe --city 新北市 --element AirTemperature,Weather
```

### 36 小時縣市概況預報

```bash
# 縣市 36 小時概況
cwa-weather overview --city 臺北市

# 篩選天氣要素
cwa-weather overview --city 臺北市 --element Wx,PoP

# 篩選時間區間
cwa-weather overview --city 臺北市 --time-from 2026-03-01T06:00:00 --time-to 2026-03-01T18:00:00
```

### 氣象警特報

```bash
# 所有生效中的警特報
cwa-weather alert

# 特定縣市的警特報
cwa-weather alert --city 臺北市
```

### 颱風動態

```bash
# 目前熱帶氣旋資訊
cwa-weather typhoon

# 依編號與資料集篩選
cwa-weather typhoon --td-no 03 --dataset ForecastData
```

### 海象觀測

```bash
# 所有海象測站
cwa-weather sea

# 特定測站
cwa-weather sea --station 富貴角
```

### 通用查詢

```bash
# 以資料集 ID 查詢任意氣象署端點
cwa-weather query F-D0047-069 -p LocationName=板橋區
```

### 列出縣市與鄉鎮

```bash
# 列出全部 22 縣市
cwa-weather cities

# 列出特定縣市的鄉鎮
cwa-weather cities --city 臺北市
```

## 作為 Go 函式庫使用

```go
package main

import (
	"context"
	"fmt"

	"github.com/kerkerj/cwa-weather/cwa"
)

func main() {
	client := cwa.NewClient("YOUR_API_KEY")
	ctx := context.Background()

	// 鄉鎮預報
	forecast, _ := client.Forecast(ctx, "臺北市", "中正區")
	fmt.Println(forecast)

	// 搭配要素與時間篩選
	filtered, _ := client.Forecast(ctx, "臺北市", "", cwa.ForecastOption{
		Element:  "溫度,天氣現象",
		TimeFrom: "2026-03-01T06:00:00",
	})
	fmt.Println(filtered)

	// 即時觀測
	obs, _ := client.Observe(ctx, cwa.ObserveByCity("新北市"))
	fmt.Println(obs)

	// 觀測搭配要素篩選
	obsFiltered, _ := client.Observe(ctx, cwa.ObserveByCity("新北市"), cwa.ObserveWithElement("AirTemperature"))
	fmt.Println(obsFiltered)
}
```

## AI Agent 整合

本專案包含 skill 檔案，讓 AI agent 可以透過 CLI 查詢天氣資料。

需要先安裝 `cwa-weather` CLI 並設定 `CWA_API_KEY` 環境變數。

### Claude Code（Plugin）

```
/plugin marketplace add kerkerj/cwa-weather
/plugin install cwa-weather@kerkerj-cwa-weather
```

安裝後直接問 Claude：
- 「台北市今天天氣如何？」
- 「現在有颱風嗎？」
- 「氣象警特報」

### 其他 AI Agent

將 agent 指向 [`skill/`](skill/) 目錄中的 skill 檔案，即可取得指令參考與使用說明。

## 備註

- **輸出格式**：一律 JSON，可搭配 `jq` 擷取欄位。
- **支援縣市**：全臺 22 縣市。
- **台→臺 自動轉換**：`台北市` 會自動轉為 `臺北市`，以符合氣象署使用正體字的慣例。

## 授權

MIT
