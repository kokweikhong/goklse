package klsescreener

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	stockPriceBaseURL = "https://www.klsescreener.com/v2/stocks/chart/"
)

type StockPrice struct {
	Date   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
}

func GetHistoricalStockPrices(ctx context.Context, stockCode string) ([]*StockPrice, error) {
	url := stockPriceBaseURL + stockCode

	prices := []*StockPrice{}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()
	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(url),
		// chromedp.Sleep(1*time.Second),
		chromedp.Reload(),
		chromedp.WaitVisible(`div#chart_div`, chromedp.ByQuery),
		chromedp.OuterHTML(`body`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	// Use goquery to parse the HTML content and extract the stock price data
	_, err = goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	// The price data is in <script> tag example data = [[1459209600000,0.8800,0.8950,0.8700,0.8950,1805900],[1459296000000,0.8900,0.9100,0.8800,0.8850,1620700],]
	// date, open, high, low, close, volume
	// use regex to extract the data
	// remove all spaces and new lines from htmlContent
	htmlContent = strings.ReplaceAll(htmlContent, " ", "")
	htmlContent = strings.ReplaceAll(htmlContent, "\n", "")
	// remove tabs from htmlContent
	htmlContent = strings.ReplaceAll(htmlContent, "\t", "")
	re := regexp.MustCompile(`data=\[(.*)\];varoriginalDrawPoint`)
	matches := re.FindStringSubmatch(htmlContent)
	if len(matches) > 1 {
		priceData := matches[1]
		fmt.Println("Extracted price data:", priceData)
		// Remove trailing comma if it exists
		priceData = strings.TrimSuffix(priceData, ",")
		// Split the price data into individual entries
		priceEntries := strings.Split(priceData, "],[")
		// for i, entry := range priceEntries {
		// 	entry = strings.Trim(entry, "[]")
		// 	priceEntries[i] = entry
		// }
		// Unmarshal the price data into a slice of StockPrice
		for _, entry := range priceEntries {
			var priceEntry []interface{}
			err = json.Unmarshal([]byte("["+entry+"]"), &priceEntry)
			if err != nil {
				fmt.Println("Error unmarshalling price entry:", err)
				continue
			}
			if len(priceEntry) != 6 {
				fmt.Println("Unexpected price entry format:", priceEntry)
				continue
			}
			prices = append(prices, &StockPrice{
				Date:   fmt.Sprintf("%v", priceEntry[0]),
				Open:   priceEntry[1].(float64),
				High:   priceEntry[2].(float64),
				Low:    priceEntry[3].(float64),
				Close:  priceEntry[4].(float64),
				Volume: int(priceEntry[5].(float64)),
			})
		}
	}

	// doc.Find(`script`).Each(func(i int, s *goquery.Selection) {
	// 	scriptContent := s.Text()
	// 	if strings.Contains(scriptContent, "data=[") {
	// 		re := regexp.MustCompile(`data=\[(.*)\];var originalDrawPoint`)
	// 		matches := re.FindStringSubmatch(scriptContent)
	// 		if len(matches) > 1 {
	// 			priceData := matches[1]
	// 			// Unmarshal the price data into a slice of StockPrice
	// 			err = json.Unmarshal([]byte(priceData), &prices)
	// 			if err != nil {
	// 				fmt.Println("Error unmarshalling price data:", err)
	// 			}
	// 		}
	// 	}
	// })

	// return prices, nil
	// }

	return prices, nil
}
