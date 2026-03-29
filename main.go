package main

import (
	"context"

	"github.com/kokweikhong/goklse/internal/pkg/klsescreener"
)

func main() {
	ctx := context.Background()
	// reports, err := klsescreener.GetQuarterlyFinancialReports(ctx, "5169")
	// if err != nil {
	// 	panic(err)
	// }

	// for _, report := range reports {
	// 	fmt.Println(report.Announced, report.Quarter, report.FinancialYear, report.EPS, report.DPS, report.NTA, report.PL, report.ROE, report.ReportLink)
	// }

	// data, err := klsescreener.GetStockListings(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// // Write the data to a file or print it to the console
	// for _, listing := range data {
	// 	println(listing.Code, listing.Name, listing.Market, listing.Sector)
	// }

	// file, err := os.Create("stock_listings.csv")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()

	// writer := csv.NewWriter(file)
	// defer writer.Flush()
	// // Write the header
	// writer.Write([]string{"Code", "Name", "LongName", "Market", "Sector", "SubSector", "Price", "FiftyTwoWeek", "Volume", "EPS", "DPS", "NTA", "PE", "DY", "ROE", "PTBV", "MarketCapInMillion"})
	// // Write the stock listings
	// for _, listing := range data {
	// 	writer.Write([]string{
	// 		listing.Code,
	// 		listing.Name,
	// 		listing.LongName,
	// 		listing.Market,
	// 		listing.Sector,
	// 		listing.SubSector,
	// 		fmt.Sprintf("%f", listing.Price),
	// 		listing.FiftyTwoWeek,
	// 		fmt.Sprintf("%d", listing.Volume),
	// 		fmt.Sprintf("%f", listing.EPS),
	// 		fmt.Sprintf("%f", listing.DPS),
	// 		fmt.Sprintf("%f", listing.NTA),
	// 		fmt.Sprintf("%f", listing.PE),
	// 		fmt.Sprintf("%f", listing.DY),
	// 		fmt.Sprintf("%f", listing.ROE),
	// 		fmt.Sprintf("%f", listing.PTBV),
	// 		fmt.Sprintf("%f", listing.MarketCapInMillion),
	// 	})
	// }

	// annualReports, err := klsescreener.GetAnnualFinancialReports(ctx, "5169")
	// if err != nil {
	// 	panic(err)
	// }

	// for _, report := range annualReports {
	// 	println(report.FinancialYear, report.Revenue, report.Net, report.EPS, report.DP, report.NetPercent, report.ReportLink)
	// }

	prices, err := klsescreener.GetHistoricalStockPrices(ctx, "5169")
	if err != nil {
		panic(err)
	}
	for _, price := range prices {
		println(price.Date, price.Open, price.High, price.Low, price.Close, price.Volume)
	}
}
