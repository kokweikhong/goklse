package klse

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/kokweikhong/goklse/types"
)

// GetStockListing returns the stock listing
func GetStockListing() []*types.StockQuote {
	var (
		quotes []*types.StockQuote
		quote  *types.StockQuote
	)
	time.Sleep(2 * time.Second)
	url := "https://www.klsescreener.com/v2/screener/quote_results"
	c := colly.NewCollector()
	c.OnHTML(`tbody tr[class="list"]`, func(e *colly.HTMLElement) {
		quote = new(types.StockQuote)
		e.ForEach("td", func(tdIndex int, el *colly.HTMLElement) {
			switch tdIndex {
			case 0:
				quote.ShortName = strings.Replace(el.Text, "[s]", "", -1)
				quote.ShortName = strings.TrimSpace(quote.ShortName)

				quote.Name = el.Attr("title")
				quote.Name = strings.TrimSpace(quote.Name)
			case 1:
				quote.Code = strings.TrimSpace(el.Text)
			case 2:
				el.ForEach("small", func(spanIndex int, small *colly.HTMLElement) {
					switch spanIndex {
					case 0:
						quote.SubCategory = strings.TrimSpace(small.Text)
					case 1:
						splitCategory := strings.Split(small.Text, ",")
						if len(splitCategory) > 1 {
							quote.Category = strings.TrimSpace(splitCategory[0])
							quote.Market = strings.TrimSpace(splitCategory[1])
						}
					}
				})
			case 3:
				price := strings.TrimSpace(el.Text)
				quote.Price, _ = strconv.ParseFloat(price, 64)
			case 4:
				changePercent := strings.Replace(el.Text, "%", "", -1)
				changePercent = strings.TrimSpace(changePercent)
				quote.ChangePercent, _ = strconv.ParseFloat(changePercent, 64)
			case 5:
				fiftyTwoWeek := strings.Split(el.Text, "-")
				if len(fiftyTwoWeek) < 2 {
					break
				}
				fiftyTwoWeekLow := strings.TrimSpace(fiftyTwoWeek[0])
				fiftyTwoWeekHigh := strings.TrimSpace(fiftyTwoWeek[1])
				quote.FiftyTwoWeekLow, _ = strconv.ParseFloat(fiftyTwoWeekLow, 64)
				quote.FiftyTwoWeekHigh, _ = strconv.ParseFloat(fiftyTwoWeekHigh, 64)
			case 6:
				volume := strings.Replace(el.Text, ",", "", -1)
				volume = strings.TrimSpace(volume)
				quote.Volume, _ = strconv.ParseFloat(volume, 64)
			case 7:
				eps := strings.TrimSpace(el.Text)
				quote.EPS, _ = strconv.ParseFloat(eps, 64)
			case 8:
				dps := strings.TrimSpace(el.Text)
				quote.DPS, _ = strconv.ParseFloat(dps, 64)
			case 9:
				nta := strings.TrimSpace(el.Text)
				quote.NTA, _ = strconv.ParseFloat(nta, 64)
			case 10:
				pe := strings.TrimSpace(el.Text)
				quote.PE, _ = strconv.ParseFloat(pe, 64)
			case 11:
				dy := strings.TrimSpace(el.Text)
				quote.DY, _ = strconv.ParseFloat(dy, 64)
			case 12:
				roe := strings.TrimSpace(el.Text)
				quote.ROE, _ = strconv.ParseFloat(roe, 64)
			case 13:
				ptbv := strings.TrimSpace(el.Text)
				quote.PTBV, _ = strconv.ParseFloat(ptbv, 64)
			case 14:
				marketCap := strings.Replace(el.Text, ",", "", -1)
				marketCap = strings.TrimSpace(marketCap)
				quote.MarketCapInMillions, _ = strconv.ParseFloat(marketCap, 64)
			}

		})
		quotes = append(quotes, quote)
	})

	c.Visit(url)
	return quotes
}

// GetStockHistoricalData returns the historical data of a stock
func GetStockHistoricalData(code string) ([]*types.OHLC, error) {
	var (
		ohlc  *types.OHLC
		ohlcs []*types.OHLC
	)
	url := "https://www.klsescreener.com/v2/stocks/chart/" + code

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

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

	bodyString := string(body)

	fmt.Println(bodyString)

	// remove new line, space and tab from string
	spaceRE := regexp.MustCompile(`\s+`)
	bodyString = spaceRE.ReplaceAllString(bodyString, "")

	dataRE := regexp.MustCompile(`data=\[(.*?),\];`)

	data := dataRE.FindStringSubmatch(bodyString)

	dataRE = regexp.MustCompile(`\[(.*?)\]`)

	dataStringList := dataRE.FindAllStringSubmatch(data[1], -1)

	for _, dataString := range dataStringList {
		if len(dataString) < 2 {
			continue
		}
		splitDataString := strings.Split(dataString[1], ",")
		if len(splitDataString) < 6 {
			continue
		}
		dateInt64, err := strconv.ParseInt(splitDataString[0], 10, 64)
		if err != nil {
			continue
		}
		fmt.Println(dateInt64)
		// convert date from int64 to time.Time
		date := time.Unix(dateInt64/1000, 0)
		fmt.Println(date)
		open, _ := strconv.ParseFloat(splitDataString[1], 64)
		high, _ := strconv.ParseFloat(splitDataString[2], 64)
		low, _ := strconv.ParseFloat(splitDataString[3], 64)
		close, _ := strconv.ParseFloat(splitDataString[4], 64)
		volume, _ := strconv.ParseFloat(splitDataString[5], 64)
		ohlc = &types.OHLC{
			Date:   date.Format("2006-01-02"),
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		}
		ohlcs = append(ohlcs, ohlc)
	}
	return ohlcs, nil
}
