package klsescreener

import (
	"context"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	financialReportBaseURL = "https://www.klsescreener.com/v2/stocks/view/"
)

type QuarterlyFinancialReport struct {
	EPS           float64
	DPS           float64
	NTA           float64
	Revenue       string
	PL            float64
	Quarter       string
	QDate         string
	FinancialYear string
	Announced     string
	ROE           string
	ReportLink    string
}

func GetQuarterlyFinancialReports(ctx context.Context, stockCode string) ([]*QuarterlyFinancialReport, error) {
	url := financialReportBaseURL + stockCode

	reports := []*QuarterlyFinancialReport{}

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
		chromedp.WaitVisible(`div#quarter_reports`, chromedp.ByQuery),
		chromedp.OuterHTML(`div#quarter_reports`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	// Use goquery to parse the HTML content and extract the financial report data
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	doc.Find(`div#quarter_reports table tbody tr`).Each(func(i int, s *goquery.Selection) {
		// Check if the first td has colspan 100, if yes then skip
		if colspan, exists := s.Find(`td:nth-child(1)`).Attr("colspan"); exists && colspan == "100" {
			return
		}
		eps := parseFloat(removeCommaForNumber(s.Find(`td:nth-child(1)`).Text()))
		dps := parseFloat(removeCommaForNumber(s.Find(`td:nth-child(2)`).Text()))
		nta := parseFloat(removeCommaForNumber(s.Find(`td:nth-child(3)`).Text()))
		pl := parseFloat(removeCommaForNumber(s.Find(`td:nth-child(4)`).Text()))
		quarter := trimAndRemoveNewLine(s.Find(`td:nth-child(5)`).Text())
		qDate := trimAndRemoveNewLine(s.Find(`td:nth-child(6)`).Text())
		financialYear := trimAndRemoveNewLine(s.Find(`td:nth-child(7)`).Text())
		announced := trimAndRemoveNewLine(s.Find(`td:nth-child(8)`).Text())
		roe := trimAndRemoveNewLine(s.Find(`td:nth-child(9)`).Text())
		reportLink, _ := s.Find(`td:nth-child(13) a[title="Financial Report"]`).Attr("href")
		if reportLink != "" {
			reportLink = "https://www.klsescreener.com" + reportLink
		}
		report := &QuarterlyFinancialReport{
			EPS:           eps,
			DPS:           dps,
			NTA:           nta,
			Revenue:       "",
			PL:            pl,
			Quarter:       quarter,
			QDate:         qDate,
			FinancialYear: financialYear,
			Announced:     announced,
			ROE:           roe,
			ReportLink:    reportLink,
		}
		reports = append(reports, report)
	})

	return reports, nil
}

type AnnualFinancialReport struct {
	FinancialYear string
	Revenue       float64
	Net           float64
	EPS           float64
	DP            string
	NetPercent    string
	ReportLink    string
}

func GetAnnualFinancialReports(ctx context.Context, stockCode string) ([]*AnnualFinancialReport, error) {
	// annual
	url := financialReportBaseURL + stockCode
	reports := []*AnnualFinancialReport{}

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
		chromedp.Sleep(2*time.Second),
		chromedp.Click(`a[href="#annual"]`, chromedp.ByQuery),
		chromedp.WaitVisible(`div#annual`, chromedp.ByQuery),
		chromedp.OuterHTML(`div#annual`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}
	// Use goquery to parse the HTML content and extract the financial report data
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	doc.Find(`div#annual table tbody tr`).Each(func(i int, s *goquery.Selection) {
		// Skip the first row
		if i == 0 {
			return
		}
		financialYear := trimAndRemoveNewLine(s.Find(`td:nth-child(1)`).Text())
		// all float need to try to remove comma before parse
		revenue := parseFloat(removeCommaForNumber(s.Find(`td:nth-child(2)`).Text()))
		net := parseFloat(removeCommaForNumber(s.Find(`td:nth-child(3)`).Text()))
		eps := parseFloat(removeCommaForNumber(s.Find(`td:nth-child(4)`).Text()))
		dp := trimAndRemoveNewLine(s.Find(`.number:nth-child(5)`).Text())
		netPercent := trimAndRemoveNewLine(s.Find(`.number:nth-child(6)`).Text())
		reportLink, _ := s.Find(`td:nth-child(7) a`).Attr("href")
		if reportLink != "" {
			reportLink = "https://www.klsescreener.com" + reportLink
		}
		report := &AnnualFinancialReport{
			FinancialYear: financialYear,
			Revenue:       revenue,
			Net:           net,
			EPS:           eps,
			DP:            dp,
			NetPercent:    netPercent,
			ReportLink:    reportLink,
		}
		reports = append(reports, report)
	})
	return reports, nil
}

func removeCommaForNumber(s string) string {
	return strings.ReplaceAll(s, ",", "")
}
