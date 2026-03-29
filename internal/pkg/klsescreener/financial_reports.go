package klsescreener

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/cdp"
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

type AnnualFinancialReport struct {
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
	nodes := []*cdp.Node{}
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`div#quarter_reports`, chromedp.ByQuery),
		chromedp.Nodes(`div#quarter_reports table tbody tr`, &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		return nil, err
	}

	// Only extract if td more than 2
	fmt.Println("Number of quarterly financial reports found:", len(nodes))
	for nodeIdx, node := range nodes {
		fmt.Printf("Processing quarterly financial report: %d/%d\n", nodeIdx+1, len(nodes))
		// how to check if td more than 2? we can check if the text of the first td colspan is 100, if it is then it means it is not a valid report
		var colspan string
		chromedp.Run(chromeCtx,
			chromedp.AttributeValue(`td:nth-child(1)`, "colspan", &colspan, nil, chromedp.ByQuery, chromedp.FromNode(node)),
		)
		if colspan == "100" {
			fmt.Println("Skipping invalid report with colspan 100")
			continue
		}
		var eps, dps, nta, pl, quarter, qDate, financialYear, announced, roe string
		if err := chromedp.Run(chromeCtx,
			chromedp.Text(`td:nth-child(1)`, &eps, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(2)`, &dps, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(3)`, &nta, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(4)`, &pl, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(5)`, &quarter, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(6)`, &qDate, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(7)`, &financialYear, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(8)`, &announced, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`td:nth-child(9)`, &roe, chromedp.ByQuery, chromedp.FromNode(node)),
		); err != nil {
			fmt.Printf("Skipping row %d: failed to read text cells: %v\n", nodeIdx+1, err)
			continue
		}

		var linkNodes []*cdp.Node
		chromedp.Run(chromeCtx,
			chromedp.Nodes(
				`td:nth-child(13) a[title="Financial Report"]`,
				&linkNodes,
				chromedp.ByQueryAll,
				chromedp.AtLeast(0),
				chromedp.FromNode(node),
			),
		)

		var reportLink string
		if len(linkNodes) > 0 {
			// Node.Attributes is a flat slice: ["href", "/path/...", "title", "Financial Report", ...]
			attrs := linkNodes[0].Attributes
			for i := 0; i < len(attrs)-1; i += 2 {
				if attrs[i] == "href" {
					reportLink = "https://www.klsescreener.com" + trimAndRemoveNewLine(attrs[i+1])
					break
				}
			}
		}
		// Convert eps, dps, nta, pl to float64
		// and remove new line and trim space for quarter, qDate, financialYear, announced, roe
		epsFloat := parseFloat(eps)
		dpsFloat := parseFloat(dps)
		ntaFloat := parseFloat(nta)
		plFloat := parseFloat(pl)
		quarter = trimAndRemoveNewLine(quarter)
		qDate = trimAndRemoveNewLine(qDate)
		financialYear = trimAndRemoveNewLine(financialYear)
		announced = trimAndRemoveNewLine(announced)
		roe = trimAndRemoveNewLine(roe)

		report := &QuarterlyFinancialReport{
			EPS:           epsFloat,
			DPS:           dpsFloat,
			NTA:           ntaFloat,
			Revenue:       "",
			PL:            plFloat,
			Quarter:       quarter,
			QDate:         qDate,
			FinancialYear: financialYear,
			Announced:     announced,
			ROE:           roe,
			ReportLink:    reportLink,
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// func parseFloat(s string) float64 {
// 	// Remove new line and trim space
// 	s = trimAndRemoveNewLine(s)
// 	// Remove commas	s = strings.ReplaceAll(s, ",", "")
// 	// Convert to float64
// 	f, err := strconv.ParseFloat(s, 64)
// 	if err != nil {
// 		return 0
// 	}
// 	return f
// }
