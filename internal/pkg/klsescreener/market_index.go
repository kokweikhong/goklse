package klsescreener

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	MarketIndexBaseURL = "https://www.klsescreener.com/v2/markets"
)

type BursaIndex struct {
	Name string
	Code string
}

func GetBursaIndexes(ctx context.Context) ([]*BursaIndex, error) {
	indexes := []*BursaIndex{}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()
	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	var htmlContent string
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(MarketIndexBaseURL),
		chromedp.WaitVisible(`div#content`, chromedp.ByQuery),
		chromedp.OuterHTML(`body`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	// Parse the HTML content using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	// Extract the index data from the HTML content
	// The index data is in div#content and div contains class row and equal and the last div contains the index data
	// document.querySelector("#content > div.wrapper-disabled > div > div.row > div > div:nth-child(31)")
	doc.Find(`#content > div.wrapper-disabled > div > div.row > div > div:nth-child(31) a`).Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		code, _ := s.Attr("href")
		// code example /v2/stocks/view/0063I, we want to extract 0063I
		code = strings.TrimPrefix(code, "/v2/stocks/view/")
		indexes = append(indexes, &BursaIndex{
			Name: name,
			Code: code,
		})
	})

	return indexes, nil
}

type MarketIndex struct {
	Name    string
	Code    string
	Country string
}

func GetMarketIndexes(ctx context.Context) ([]*MarketIndex, error) {
	indexes := []*MarketIndex{}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()
	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	var htmlContent string
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(MarketIndexBaseURL),
		chromedp.WaitVisible(`div#content`, chromedp.ByQuery),
		chromedp.OuterHTML(`body`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	// Parse the HTML content using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	// Extract the index data from the HTML content
	// document.querySelector("#content > div.wrapper-disabled > div > div.row > div > div:nth-child(2)")
	// The country is document.querySelector("#content > div.wrapper-disabled > div > div.row > div > div:nth-child(2) > div:nth-child(3) > div > div.col-sm-7 > div:nth-child(2)")
	doc.Find(`#content > div.wrapper-disabled > div > div.row > div > div:nth-child(2) a`).Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		code, _ := s.Attr("href")
		// code example /v2/stocks/view/0063I, we want to extract 0063I
		// Get the last part of the code after the last slash instead of the first part after /v2/stocks/view/
		code = code[strings.LastIndex(code, "/")+1:]
		country := s.Parent().Parent().Parent().Find(`div.col-sm-7 div:nth-child(2)`).Text()
		indexes = append(indexes, &MarketIndex{
			Name:    trimAndRemoveNewLine(name),
			Code:    code,
			Country: trimAndRemoveNewLine(country),
		})
	})
	return indexes, nil
}

type ComponentStock struct {
	Name      string
	Code      string
	MarketCap string
	Category  string
}

const (
	KLCIComponentStocksURL   = "https://www.klsescreener.com/v2/stocks/view/0200I"
	KLCI70ComponentStocksURL = "https://www.klsescreener.com/v2/stocks/view/0863I"
)

func GetKLCIComponentStocks(ctx context.Context) ([]*ComponentStock, error) {

	KLCIstocks := []*ComponentStock{}
	KLCI70Stocks := []*ComponentStock{}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()
	chromeCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	var htmlContent string
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(KLCIComponentStocksURL),
		chromedp.WaitVisible(`div#page`, chromedp.ByQuery),
		chromedp.OuterHTML(`body`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	// Parse the HTML content using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	// Extract the component stock data from the HTML content
	doc.Find(`#page table tbody tr`).Each(func(i int, s *goquery.Selection) {
		name := s.Find(`td:nth-child(2) a`).Text()
		code := s.Find(`td:nth-child(2) a`).AttrOr("href", "")
		code = code[strings.LastIndex(code, "/")+1:]
		marketCap := s.Find(`td:nth-child(5)`).Text()
		// If name and code is empty, skip
		if strings.TrimSpace(name) == "" || strings.TrimSpace(code) == "" {
			return
		}
		KLCIstocks = append(KLCIstocks, &ComponentStock{
			Name:      trimAndRemoveNewLine(name),
			Code:      code,
			MarketCap: trimAndRemoveNewLine(marketCap),
			Category:  "KLCI",
		})
	})

	err = chromedp.Run(chromeCtx,
		chromedp.Navigate(KLCI70ComponentStocksURL),
		chromedp.WaitVisible(`div#page`, chromedp.ByQuery),
		chromedp.OuterHTML(`body`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}
	// Parse the HTML content using goquery
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	// Extract the component stock data from the HTML content
	doc.Find(`#page table tbody tr`).Each(func(i int, s *goquery.Selection) {
		name := s.Find(`td:nth-child(2) a`).Text()
		code := s.Find(`td:nth-child(2) a`).AttrOr("href", "")
		code = code[strings.LastIndex(code, "/")+1:]
		marketCap := s.Find(`td:nth-child(5)`).Text()
		// If name and code is empty, skip
		if strings.TrimSpace(name) == "" || strings.TrimSpace(code) == "" {
			return
		}
		KLCI70Stocks = append(KLCI70Stocks, &ComponentStock{
			Name:      trimAndRemoveNewLine(name),
			Code:      code,
			MarketCap: trimAndRemoveNewLine(marketCap),
			Category:  "KLCI70",
		})
	})
	return append(KLCIstocks, KLCI70Stocks...), nil
}
