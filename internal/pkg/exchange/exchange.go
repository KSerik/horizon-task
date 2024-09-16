package exchange

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"
	"strings"
	"time"
)

type Service struct {
	cache map[string]decimal.Decimal // symbol+date -> rate
}

var coinSymbolIDs = map[string]string{
	"USDC":   "usd-coin",
	"USDC.E": "bridged-usdc-polygon-pos-bridge",
	"SFL":    "sunflower-land",
	"MATIC":  "matic-network",
	"BTC":    "bitcoin",
	"ETH":    "ethereum",
	"USDT":   "tether",
	"BNB":    "binancecoin",
	"SOL":    "solana",
	"XRP":    "ripple",
	"DOGE":   "dogecoin",
	"TON":    "the-open-network",
}

func NewService() *Service {
	return &Service{
		cache: map[string]decimal.Decimal{
			// adding historical data to speed up the process
			"SFL15-04-2024":    decimal.NewFromFloat(0.10629575068267864),
			"MATIC15-04-2024":  decimal.NewFromFloat(0.713286899222239),
			"USDC15-04-2024":   decimal.NewFromFloat(1.0003755218380892),
			"USDC.E15-04-2024": decimal.NewFromFloat(0.9999837285279117),
			"SFL01-04-2024":    decimal.NewFromFloat(0.0826953877015709),
			"MATIC01-04-2024":  decimal.NewFromFloat(1.0029078147445294),
			"USDC.E01-04-2024": decimal.NewFromFloat(1.0002272213537717),
			"USDC01-04-2024":   decimal.NewFromFloat(0.9998785322687238),
			"SFL16-04-2024":    decimal.NewFromFloat(0.11008186554416675),
			"MATIC16-04-2024":  decimal.NewFromFloat(0.7075889920386759),
			"USDC16-04-2024":   decimal.NewFromFloat(1.0007005573707544),
			"USDC.E16-04-2024": decimal.NewFromFloat(1.0003575893817702),
			"SFL02-04-2024":    decimal.NewFromFloat(0.08023586067764414),
			"USDC.E02-04-2024": decimal.NewFromFloat(0.999766574397093),
			"MATIC02-04-2024":  decimal.NewFromFloat(0.9524791332891662),
		},
	}
}

func (e *Service) GetExchangeRate(symbol, date string) (decimal.Decimal, error) {
	if rate, ok := e.cache[symbol+date]; ok {
		return rate, nil
	}

	// adding this sleep to avoid rate limiting from coingecko API
	time.Sleep(15 * time.Second)

	coinID, ok := coinSymbolIDs[strings.ToUpper(symbol)]
	if !ok {
		return decimal.Zero, fmt.Errorf("unknown coin symbol: %s", symbol)
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/history?date=%s&localization=false", coinID, date)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return decimal.Zero, fmt.Errorf("creating http request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("x-cg-api-key", "CG-qeBSXDoVzhVuWvbLpArnCbtf")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return decimal.Zero, fmt.Errorf("making http request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("response status is not OK: %s", resp.Status)
	}

	var response struct {
		MarketData struct {
			CurrentPrice struct {
				USD float64 `json:"usd"`
			} `json:"current_price"`
		} `json:"market_data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return decimal.Zero, fmt.Errorf("decoding response: %w", err)
	}

	if response.MarketData.CurrentPrice.USD == 0 {
		return decimal.Zero, fmt.Errorf("exchange rate is 0")
	}

	usdPrice := decimal.NewFromFloat(response.MarketData.CurrentPrice.USD)
	e.cache[symbol+date] = usdPrice

	return usdPrice, nil
}
