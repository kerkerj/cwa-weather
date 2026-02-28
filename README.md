# cwa-weather

CWA Open Data API CLI and Go library for Taiwan weather data.

Queries forecast and real-time observation data from Taiwan's Central Weather Administration (CWA) Open Data platform. All output is JSON, designed for agent and jq consumption.

## Installation

```bash
go install github.com/kerkerj/cwa-weather/cmd/cwa-weather@latest
```

Or download a binary from [GitHub Releases](https://github.com/kerkerj/cwa-weather/releases).

## Setup

Get an API key from https://opendata.cwa.gov.tw and export it:

```bash
export CWA_API_KEY=your-key
```

## Usage

### Forecast

```bash
# Township-level forecast
cwa-weather forecast --city 臺北市 --town 中正區

# City-level forecast (all towns)
cwa-weather forecast --city 台北市    # 台→臺 auto-converted
```

### Real-time Observation

```bash
# By city
cwa-weather observe --city 新北市

# By station name
cwa-weather observe --station 淡水
```

### Generic Query

```bash
# Query any CWA endpoint by data ID
cwa-weather query F-D0047-069 -p LocationName=板橋區
```

### List Cities and Towns

```bash
# List all 22 cities
cwa-weather cities

# List towns in a city
cwa-weather cities --city 臺北市
```

## Library Usage

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

	// Forecast
	forecast, _ := client.Forecast(ctx, "臺北市", "中正區")
	fmt.Println(forecast)

	// Observation
	obs, _ := client.Observe(ctx, cwa.ObserveByCity("新北市"))
	fmt.Println(obs)
}
```

## Notes

- **Output**: Always JSON. Pipe to `jq` for field extraction.
- **Supported cities**: All 22 Taiwan cities/counties.
- **台→臺 auto-conversion**: `台北市` is automatically converted to `臺北市` to match CWA naming.

## License

MIT
