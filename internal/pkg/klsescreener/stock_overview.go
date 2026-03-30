package klsescreener

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	StockOverviewBaseURL = "https://www.klsescreener.com/v2/stocks/view/"
)

type CompanySummary struct {
	Summary string
	Website string
}

func GetCompanySummary(ctx context.Context, stockCode string) (*CompanySummary, error) {
	url := StockOverviewBaseURL + stockCode

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
		// chromedp.WaitVisible(`div[data-target="#company_summary"]`, chromedp.ByQuery),
		chromedp.Sleep(time.Second),
		chromedp.OuterHTML(`body`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	// Use goquery to parse the HTML content and extract the company summary data
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	companySummary := &CompanySummary{}
	doc.Find(`div#company_summary div.modal-body`).Each(func(i int, s *goquery.Selection) {
		companySummary.Summary = s.Text()
		// Remove new line and trim the company summary and also tabs
		companySummary.Summary = strings.ReplaceAll(companySummary.Summary, "\n", "")
		companySummary.Summary = strings.ReplaceAll(companySummary.Summary, "\t", "")
		companySummary.Summary = strings.TrimSpace(companySummary.Summary)

		// The website is in Website: https://...from summary, we can use regex to extract the website
		re := regexp.MustCompile(`Website:\s*(https?://[^\s]+)`)
		matches := re.FindStringSubmatch(companySummary.Summary)
		if len(matches) > 1 {
			companySummary.Website = matches[1]
		}
		// Remove the website from the summary
		companySummary.Summary = re.ReplaceAllString(companySummary.Summary, "")
		// Remove new line and trim the company summary and also tabs again after removing website
		companySummary.Summary = strings.ReplaceAll(companySummary.Summary, "\n", "")
		companySummary.Summary = strings.ReplaceAll(companySummary.Summary, "\t", "")
		companySummary.Summary = strings.TrimSpace(companySummary.Summary)
		// Remove "Website:" from the summary and website if it is still there
		companySummary.Summary = strings.ReplaceAll(companySummary.Summary, "Website:", "")
		companySummary.Summary = strings.TrimSpace(companySummary.Summary)
		// Also in website
		companySummary.Website = strings.ReplaceAll(companySummary.Website, "Website:", "")
		companySummary.Website = strings.TrimSpace(companySummary.Website)

	})

	return companySummary, nil
}
