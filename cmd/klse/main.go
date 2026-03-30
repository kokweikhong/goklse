package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kokweikhong/goklse/internal/pkg/klsescreener"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "klse",
		Usage: "KLSE CLI tools",
		Commands: []*cli.Command{
			// Get stock listings from KLSE Screener
			{
				Name:      "get-stock-listings",
				Aliases:   []string{"gsl"},
				Usage:     "Get stock listings from KLSE Screener",
				UsageText: "klse get-stock-listings [--export-csv | -csv]",
				ArgsUsage: "Use --export-csv or -csv flag to export the stock listings to stock_listings.csv",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "export-csv",
						Aliases: []string{"csv"},
						Usage:   "Export stock listings to stock_listings.csv",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					listings, err := klsescreener.GetStockListings(ctx)
					if err != nil {
						return err
					}

					for _, listing := range listings {
						fmt.Println(listing.Code, listing.Name, listing.Market, listing.Sector)
					}

					// filename := "stock_listings_timestamp.csv"
					filename := fmt.Sprintf("stock_listings_%s.csv", time.Now().Format("20060102_150405"))
					if cmd.Bool("export-csv") || cmd.Bool("csv") {
						if err := klsescreener.ExportStockListingsToCSV(ctx, listings, filename); err != nil {
							return err
						}
						fmt.Printf("Exported stock listings to %s\n", filename)
					}

					return nil
				},
			},

			// Get historical stock prices from KLSE Screener
			{
				Name:      "get-historical-stock-prices",
				Aliases:   []string{"ghsp"},
				Usage:     "Get historical stock prices from KLSE Screener",
				UsageText: "klse get-historical-stock-prices [stock-code]",
				ArgsUsage: "Example: klse get-historical-stock-prices 5169",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "code",
						Usage:    "Stock code to get historical stock prices for. Example: 5169",
						Required: true,
						Aliases:  []string{"c"},
					},
					&cli.BoolFlag{
						Name:    "export-csv",
						Aliases: []string{"csv"},
						Usage:   "Export historical stock prices to historical_stock_prices.csv",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					stockCode := cmd.String("code")
					prices, err := klsescreener.GetHistoricalStockPrices(ctx, stockCode)
					if err != nil {
						return err
					}
					for _, price := range prices {
						// Convert the date from timestamp to string
						date := time.Unix(price.Timestamp/1000, 0).Format("2006-01-02")
						fmt.Println(date, price.Open, price.High, price.Low, price.Close, price.Volume)
					}

					filename := fmt.Sprintf("%s_historical_stock_prices_%s.csv", stockCode, time.Now().Format("20060102_150405"))
					if cmd.Bool("export-csv") || cmd.Bool("csv") {
						if err := klsescreener.ExportHistoricalStockPricesToCSV(ctx, prices, filename); err != nil {
							return err
						}
						fmt.Printf("Exported historical stock prices to %s\n", filename)
					}

					return nil
				},
			},

			// Get financial reports from KLSE Screener
			{
				Name:      "get-financial-reports",
				Aliases:   []string{"gfr"},
				Usage:     "Get financial reports from KLSE Screener",
				UsageText: "klse get-financial-reports -type [annual | quarterly] -code [stock-code]",
				ArgsUsage: "Example: klse get-financial-reports -type annual -code 5169",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "type",
						Usage:    "Type of financial report to get. Example: annual or quarterly",
						Required: true,
						Aliases:  []string{"t"},
					},
					&cli.StringFlag{
						Name:     "code",
						Usage:    "Stock code to get financial reports for. Example: 5169",
						Required: true,
						Aliases:  []string{"c"},
					},
					&cli.BoolFlag{
						Name:    "export-csv",
						Aliases: []string{"csv"},
						Usage:   "Export financial reports to financial_reports.csv",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					reportType := cmd.String("type")
					stockCode := cmd.String("code")
					if reportType == "annual" || reportType == "a" {
						annualReports, err := klsescreener.GetAnnualFinancialReports(ctx, stockCode)
						if err != nil {
							return err
						}
						for _, report := range annualReports {
							fmt.Println(*report)
						}
						if cmd.Bool("export-csv") || cmd.Bool("csv") {
							filename := fmt.Sprintf("%s_annual_financial_reports_%s.csv", stockCode, time.Now().Format("20060102_150405"))
							if err := klsescreener.ExportAnnualFinancialReportsToCSV(ctx, annualReports, filename); err != nil {
								return err
							}
							fmt.Printf("Exported annual financial reports to %s\n", filename)
						}
					} else if reportType == "quarterly" || reportType == "q" {
						quarterlyReports, err := klsescreener.GetQuarterlyFinancialReports(ctx, stockCode)
						if err != nil {
							return err
						}
						for _, report := range quarterlyReports {
							fmt.Println(*report)
						}
						if cmd.Bool("export-csv") || cmd.Bool("csv") {
							filename := fmt.Sprintf("%s_quarterly_financial_reports_%s.csv", stockCode, time.Now().Format("20060102_150405"))
							if err := klsescreener.ExportQuarterlyFinancialReportsToCSV(ctx, quarterlyReports, filename); err != nil {
								return err
							}
							fmt.Printf("Exported quarterly financial reports to %s\n", filename)
						}
					} else {
						return fmt.Errorf("invalid report type: %s. Valid types are annual (a) and quarterly (q)", reportType)
					}

					return nil
				},
			},

			// Get market indexes from KLSE Screener
			{
				Name:    "get-indexes",
				Aliases: []string{"gmi"},
				Usage:   "Get indexes from KLSE Screener",
				// Flags for market or bursa
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "category",
						Usage:    "Category to get indexes for. Example: bursa or market",
						Required: true,
						Aliases:  []string{"c"},
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					category := cmd.String("category")
					if category == "bursa" || category == "b" {
						indexes, err := klsescreener.GetBursaIndexes(ctx)
						if err != nil {
							return err
						}
						for _, index := range indexes {
							fmt.Println(index.Name, index.Code)
						}
						fmt.Println("Note: To get historical stock prices for the indexes, klse get-historical-stock-prices -code [index-code]")
						fmt.Println("Example: klse get-historical-stock-prices -code 0063I")
					} else if category == "market" || category == "m" {
						indexes, err := klsescreener.GetMarketIndexes(ctx)
						if err != nil {
							return err
						}
						for _, index := range indexes {
							fmt.Println(index.Name, index.Code, index.Country)
						}
						fmt.Println("Note: To get historical stock prices for the indexes, klse get-historical-stock-prices -code [index-code]")
						fmt.Println("Example: klse get-historical-stock-prices -code FBMKLCI")
					} else {
						return fmt.Errorf("invalid category: %s. Valid categories are bursa (b) and market (m)", category)
					}
					return nil
				},
			},

			// Get stock overview from KLSE Screener
			{
				Name:    "get-stock-overview",
				Aliases: []string{"gso"},
				Usage:   "Get stock overview from KLSE Screener",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "category",
						Usage:    "Category to get stock overview for. Example: summary or financials",
						Required: true,
						Aliases:  []string{"s", "f"},
					},
					&cli.StringFlag{
						Name:     "code",
						Usage:    "Stock code to get stock overview for. Example: 5169",
						Required: true,
						Aliases:  []string{"c"},
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					category := cmd.String("category")
					stockCode := cmd.String("code")
					if category == "summary" || category == "s" {
						summary, err := klsescreener.GetCompanySummary(ctx, stockCode)
						if err != nil {
							return err
						}
						fmt.Println(summary.Summary)
						fmt.Println("Company website:", summary.Website)
					} else {
						return fmt.Errorf("invalid category: %s. Valid categories are summary (s) and financials (f)", category)
					}
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
