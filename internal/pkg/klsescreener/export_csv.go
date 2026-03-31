package klsescreener

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func ExportStockListingsToCSV(ctx context.Context, listings []*StockListing, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write the header

	writer.Write([]string{"Code", "Name", "LongName", "Market", "Sector", "SubSector", "Price", "FiftyTwoWeek", "Volume", "EPS", "DPS", "NTA", "PE", "DY", "ROE", "PTBV", "MarketCapInMillion"})
	// Write the stock listings
	for _, listing := range listings {
		writer.Write([]string{
			listing.Code,
			listing.Name,
			listing.LongName,
			listing.Market,
			listing.Sector,
			listing.SubSector,
			fmt.Sprintf("%f", listing.Price),
			listing.FiftyTwoWeek,
			fmt.Sprintf("%d", listing.Volume),
			fmt.Sprintf("%f", listing.EPS),
			fmt.Sprintf("%f", listing.DPS),
			fmt.Sprintf("%f", listing.NTA),
			fmt.Sprintf("%f", listing.PE),
			fmt.Sprintf("%f", listing.DY),
			fmt.Sprintf("%f", listing.ROE),
			fmt.Sprintf("%f", listing.PTBV),
			fmt.Sprintf("%f", listing.MarketCapInMillion),
		})
	}
	return nil
}

func ExportHistoricalStockPricesToCSV(ctx context.Context, prices []*StockPrice, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write the header
	writer.Write([]string{"Timestamp", "Date", "Open", "High", "Low", "Close", "Volume"})
	// Write the stock prices
	for _, price := range prices {
		date := time.Unix(price.Timestamp/1000, 0).Format("2006-01-02")
		writer.Write([]string{
			fmt.Sprintf("%d", price.Timestamp),
			date,
			fmt.Sprintf("%f", price.Open),
			fmt.Sprintf("%f", price.High),
			fmt.Sprintf("%f", price.Low),
			fmt.Sprintf("%f", price.Close),
			fmt.Sprintf("%d", price.Volume),
		})
	}
	return nil
}

func ExportQuarterlyFinancialReportsToCSV(ctx context.Context, reports []*QuarterlyFinancialReport, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write the header
	writer.Write([]string{"EPS", "DPS", "NTA", "Revenue", "PL", "Quarter", "QDate", "FinancialYear", "Announced", "ROE", "ReportLink"})
	// Write the financial reports
	for _, report := range reports {
		writer.Write([]string{
			fmt.Sprintf("%f", report.EPS),
			fmt.Sprintf("%f", report.DPS),
			fmt.Sprintf("%f", report.NTA),
			report.Revenue,
			fmt.Sprintf("%f", report.PL),
			report.Quarter,
			report.QDate,
			report.FinancialYear,
			report.Announced,
			report.ROE,
			report.ReportLink,
		})
	}
	return nil
}

func ExportAnnualFinancialReportsToCSV(ctx context.Context, reports []*AnnualFinancialReport, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write the header
	writer.Write([]string{"FinancialYear", "Revenue", "Net", "EPS", "DP", "NetPercent", "ReportLink"})
	// Write the financial reports
	for _, report := range reports {
		writer.Write([]string{
			report.FinancialYear,
			fmt.Sprintf("%f", report.Revenue),
			fmt.Sprintf("%f", report.Net),
			fmt.Sprintf("%f", report.EPS),
			report.DP,
			report.NetPercent,
			report.ReportLink,
		})
	}
	return nil
}

func ExportKLCIComponentStocksToCSV(ctx context.Context, stocks []*ComponentStock, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write the header
	writer.Write([]string{"Name", "Code", "MarketCap", "Category"})
	// Write the component stocks
	for _, stock := range stocks {
		writer.Write([]string{
			stock.Name,
			stock.Code,
			stock.MarketCap,
			stock.Category,
		})
	}
	return nil
}
