package klsescreener

import (
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const (
	stockListingURL = "https://www.klsescreener.com/v2"
)

type StockListing struct {
	Code string
	Name string
	LongName string
	Market string
	Sector string
	SubSector string
	Price float64
	FiftyTwoWeek string
	Volume int
	EPS float64
	DPS float64
	NTA float64
	PE float64
	DY float64
	ROE float64
	PTBV float64
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
	nodes := []*cdp.Node{}
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(stockListingURL),
		chromedp.WaitVisible(`div#by_conditions`, chromedp.ByQuery),
		// click input #submit
		chromedp.Click(`input#submit`, chromedp.ByQuery),
		chromedp.Sleep(2), // wait for the page to load
		chromedp.Nodes(`div#result table tbody tr[role="row"]`, &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		return nil, err
	}

	fmt.Println("Number of stock listings found:", len(nodes))

	for nodeIdx, node := range nodes {
		// Print info and progress
		fmt.Printf("Processing stock listing: %d/%d\n", nodeIdx+1, len(nodes))
		var name, code, longName, market, sector, subSector, fiftyTwoWeek,volume,  price, eps, dps, nta, pe, dy, roe, ptbv, marketCapInMillion string
		chromedp.Run(chromeCtx,
			chromedp.Text(`td:nth-child(2)`, &code, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(1)`, &name, chromedp.ByQuery, chromedp.FromNode(node)),
			// Extract the long name from the title attribute of the td:nth-child(1) element
			chromedp.AttributeValue(`td:nth-child(1)`, "title", &longName, nil, chromedp.ByQuery, chromedp.FromNode(node)),
			// Market is in td:nth-child(3) last small element and sector is in the first small element
			chromedp.Text(`td:nth-child(3) small:last-child`, &market, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(3) small:nth-child(1)`, &subSector, chromedp.ByQuery, chromedp.FromNode(node)),
			// Price is in td:nth-child(4) and needs to be converted to float64
			chromedp.Text(`td:nth-child(4)`, &price, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(7)`, &fiftyTwoWeek, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(8)`, &volume, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(9)`, &eps, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(10)`, &dps, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(11)`, &nta, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(12)`, &pe, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(13)`, &dy, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(14)`, &roe, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(15)`, &ptbv, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(16)`, &marketCapInMillion, chromedp.ByQuery, chromedp.FromNode(node)),
		)
		// Sector is in market text and needs to be extracted using regex
		// example market text is " Energy, Main Market"
		// we can split the market text by comma and get the first part as sector
		marketParts := splitAndTrim(market, ",")
		if len(marketParts) > 0 {
			sector = marketParts[0]
			market = marketParts[len(marketParts)-1]
		}

		// All need to remove new line and trim space
		// and convert to correct type
		name = trimAndRemoveNewLine(name)
		code = trimAndRemoveNewLine(code)
		longName = trimAndRemoveNewLine(longName)
		market = trimAndRemoveNewLine(market)
		sector = trimAndRemoveNewLine(sector)
		subSector = trimAndRemoveNewLine(subSector)
		fiftyTwoWeek = trimAndRemoveNewLine(fiftyTwoWeek)
		volume = trimAndRemoveNewLine(volume)
		price = trimAndRemoveNewLine(price)
		eps = trimAndRemoveNewLine(eps)
		dps = trimAndRemoveNewLine(dps)
		nta = trimAndRemoveNewLine(nta)
		pe = trimAndRemoveNewLine(pe)
		dy = trimAndRemoveNewLine(dy)
		roe = trimAndRemoveNewLine(roe)
		ptbv = trimAndRemoveNewLine(ptbv)
		marketCapInMillion = trimAndRemoveNewLine(marketCapInMillion)
		stockListings = append(stockListings, &StockListing{
			Code: code,
			Name: name,
			LongName: longName,
			Market: market,
			Sector: sector,
			SubSector: subSector,
			Price: parseFloat(price),
			FiftyTwoWeek: fiftyTwoWeek,
			Volume: parseInt(volume),
			EPS: parseFloat(eps),
			DPS: parseFloat(dps),
			NTA: parseFloat(nta),
			PE: parseFloat(pe),
			DY: parseFloat(dy),
			ROE: parseFloat(roe),
			PTBV: parseFloat(ptbv),
			MarketCapInMillion: parseFloat(marketCapInMillion),
		})
	}

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