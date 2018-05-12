package coinmarketcap

import (
	"fmt"

	"github.com/cleanunicorn/poacher/pkg/util"
)

const tickerURL = "https://api.coinmarketcap.com/v2/ticker/?start=%d&limit=%d&convert=%s"

const (
	// TickerStartDefault default start index
	TickerStartDefault = 0
	// TickerLimitMax default limit
	TickerLimitMax = 100
	// TickerConvertDefault default currency conversion
	TickerConvertDefault = "USD"
)

// TickerItem is an item in the retrieved Ticker
type TickerItem struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	Symbol            string  `json:"symbol"`
	WebsiteSlug       string  `json:"website_slug"`
	Rank              int     `json:"rank"`
	CirculatingSupply float64 `json:"circulating_supply"`
	TotalSupply       float64 `json:"total_supply"`
	MaxSupply         float64 `json:"max_supply"`
	Quotes            map[string]struct {
		Price            float64 `json:"price"`
		Volume24H        float64 `json:"volume_24h"`
		MarketCap        float64 `json:"market_cap"`
		PercentChange1H  float64 `json:"percent_change_1h"`
		PercentChange24H float64 `json:"percent_change_24h"`
		PercentChange7D  float64 `json:"percent_change_7d"`
	} `json:"quotes"`
	LastUpdated int `json:"last_updated"`
}

// ResponseTickers struct
type ResponseTickers struct {
	Data     map[string]TickerItem `json:"data"`
	Metadata struct {
		Timestamp           int64       `json:"timestamp"`
		NumCryptocurrencies int         `json:"num_cryptocurrencies"`
		Error               interface{} `json:"error"`
	} `json:"metadata"`
}

// Ticker retrieves ticker
func Ticker(start, limit int, convert string) (ResponseTickers, error) {
	var ts ResponseTickers
	err := util.GetJSON(
		fmt.Sprintf(tickerURL, start, limit, convert),
		&ts,
	)
	if err != nil {
		return ResponseTickers{}, err
	}

	return ts, nil
}
