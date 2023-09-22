package klse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kokweikhong/goklse/types"
)

// YahooKLCIHistoricalData is the struct for Yahoo KLCI historical data
type YahooKLCIHistoricalData struct {
	Chart struct {
		Result []struct {
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []float64 `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
			Timestamp []int64 `json:"timestamp"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

// GetBursaMarketHistoricalData returns the historical data of Bursa Malaysia market
func GetBursaMarketHistoricalData() ([]*types.OHLC, error) {
	var (
		ohlcs                   []*types.OHLC
		ohlc                    *types.OHLC
		yahooKLCIHistoricalData *YahooKLCIHistoricalData
	)

	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v8/finance/chart/%v;1=9?formatted=true&crumb=nu8VnWznYeH&lang=en-US&region=US&includeAdjustedClose=true&interval=1d&period1=1420099200&period2=%v&events=capitalGain%vdiv%vsplit&useYfid=true&corsDomain=finance.yahoo.com",
		"%5EKLSE", time.Now().Unix(), "%7C", "%7C",
	)
	fmt.Println(url)
	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/92.0.4515.159 Chrome/92.0.4515.159 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err := json.Unmarshal(body, &yahooKLCIHistoricalData); err != nil {
		return nil, err
	}

	for k, v := range yahooKLCIHistoricalData.Chart.Result[0].Timestamp {
		ohlc = new(types.OHLC)
		ohlc.Date = time.Unix(v, 0).Format("2006-01-02")
		ohlc.Open = yahooKLCIHistoricalData.Chart.Result[0].Indicators.Quote[0].Open[k]
		ohlc.High = yahooKLCIHistoricalData.Chart.Result[0].Indicators.Quote[0].High[k]
		ohlc.Low = yahooKLCIHistoricalData.Chart.Result[0].Indicators.Quote[0].Low[k]
		ohlc.Close = yahooKLCIHistoricalData.Chart.Result[0].Indicators.Quote[0].Close[k]
		ohlc.Volume = yahooKLCIHistoricalData.Chart.Result[0].Indicators.Quote[0].Volume[k]

		ohlcs = append(ohlcs, ohlc)
	}

	return ohlcs, nil
}
