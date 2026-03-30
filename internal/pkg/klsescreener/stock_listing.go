package klsescreener

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	stockListingURL = "https://www.klsescreener.com/v2"
)

type StockListing struct {
	Code               string
	Name               string
	LongName           string
	Market             string
	Sector             string
	SubSector          string
	Price              float64
	FiftyTwoWeek       string
	Volume             int
	EPS                float64
	DPS                float64
	NTA                float64
	PE                 float64
	DY                 float64
	ROE                float64
	PTBV               float64
	MarketCapInMillion float64
}

func GetStockListings(ctx context.Context) ([]*StockListing, error) {
	var stockListings []*StockListing
	// Pass the options to chromedp eg headless, window size, etc
	// example opts := append(chromedp.DefaultExecAllocatorOptions[:],
	//chromedp.ExecPath("C:/Program Files/Google/Chrome/Application/chrome.exe"),        // Windows
	// 	chromedp.ExecPath("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"), // MacOS
	// 	chromedp.Flag("headless", false),
	// )
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()
	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(stockListingURL),
		chromedp.WaitVisible(`div#by_conditions`, chromedp.ByQuery),
		// click input #submit
		chromedp.Click(`input#submit`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second), // wait for the page to load
		chromedp.WaitVisible(`div#result`, chromedp.ByQuery),
		chromedp.OuterHTML(`div#result`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	doc.Find(`div#result table tbody tr[role="row"]`).Each(func(i int, s *goquery.Selection) {
		code := trimAndRemoveNewLine(s.Find(`td:nth-child(2)`).Text())
		name := trimAndRemoveNewLine(s.Find(`td:nth-child(1)`).Text())
		// PPB                                                [s], if contains "[s]" then remove [s] and trim space
		name = strings.ReplaceAll(name, "[s]", "")
		name = trimAndRemoveNewLine(name)
		longName, _ := s.Find(`td:nth-child(1)`).Attr("title")
		market := trimAndRemoveNewLine(s.Find(`td:nth-child(3) small:last-child`).Text())
		subSector := trimAndRemoveNewLine(s.Find(`td:nth-child(3) small:nth-child(1)`).Text())
		price := parseFloat(removeCommaForNumber(trimAndRemoveNewLine(s.Find(`td:nth-child(4)`).Text())))
		fiftyTwoWeek := trimAndRemoveNewLine(s.Find(`td:nth-child(7)`).Text())
		volume := parseInt(removeCommaForNumber(trimAndRemoveNewLine(s.Find(`td:nth-child(8)`).Text())))
		eps := parseFloat(removeCommaForNumber(trimAndRemoveNewLine(s.Find(`td:nth-child(9)`).Text())))
		dps := parseFloat(removeCommaForNumber(trimAndRemoveNewLine(s.Find(`td:nth-child(10)`).Text())))
		nta := parseFloat(removeCommaForNumber(trimAndRemoveNewLine(s.Find(`td:nth-child(11)`).Text())))
		pe := parseFloat(removeCommaForNumber(trimAndRemoveNewLine(s.Find(`td:nth-child(12)`).Text())))
		dy := parseFloat(removeCommaForNumber(trimAndRemoveNewLine(s.Find(`td:nth-child(13)`).Text())))
		roe := parseFloat(removeCommaForNumber(trimAndRemoveNewLine(s.Find(`td:nth-child(14)`).Text())))
		ptbv := parseFloat(removeCommaForNumber(trimAndRemoveNewLine(s.Find(`td:nth-child(15)`).Text())))
		marketCapInMillion := parseFloat(removeCommaForNumber(trimAndRemoveNewLine(s.Find(`td:nth-child(16)`).Text())))
		// Sector is in market text and needs to be extracted using regex
		// example market text is " Energy, Main Market"
		// we can split the market text by comma and get the first part as sector
		sector := ""
		marketParts := splitAndTrim(market, ",")
		if len(marketParts) > 0 {
			sector = marketParts[0]
			market = marketParts[len(marketParts)-1]
		}
		stockListings = append(stockListings, &StockListing{
			Code:               code,
			Name:               name,
			LongName:           longName,
			Market:             market,
			Sector:             sector,
			SubSector:          subSector,
			Price:              price,
			FiftyTwoWeek:       fiftyTwoWeek,
			Volume:             volume,
			EPS:                eps,
			DPS:                dps,
			NTA:                nta,
			PE:                 pe,
			DY:                 dy,
			ROE:                roe,
			PTBV:               ptbv,
			MarketCapInMillion: marketCapInMillion,
		})
	})

	fmt.Printf("Extracted %d stock listings\n", len(stockListings))

	return stockListings, nil

}

func splitAndTrim(s string, sep string) []string {
	parts := []string{}
	for _, part := range strings.Split(s, sep) {
		parts = append(parts, trimAndRemoveNewLine(part))
	}
	return parts
}

func trimAndRemoveNewLine(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimSpace(s)
	return s
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
