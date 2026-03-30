package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/kokweikhong/goklse/internal/pkg/klsescreener"
)

func main() {
	// Get stock listings and company summaries
	listingsWithSummaries, err := combineStockListingsAndCompanySummaries()
	if err != nil {
		panic(err)
	}

	// Write the combined stock listings with company summaries to a CSV file
	filename := fmt.Sprintf("stock_listings_with_summaries_%s.csv", time.Now().Format("20060102_150405"))
	// use os.Create to create the file and then use csv.NewWriter to write the data to the file
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()
	// Write the header
	writer.Write([]string{"Code", "Name", "LongName", "Market", "Sector", "SubSector", "Price", "FiftyTwoWeek", "Volume", "EPS", "DPS", "NTA", "PE", "DY", "ROE", "PTBV", "MarketCapInMillion", "Summary", "Website"})
	// Write the combined stock listings with company summaries	for _, item := range listingsWithSummaries {
	for _, item := range listingsWithSummaries {
		writer.Write([]string{
			item.Code,
			item.Name,
			item.LongName,
			item.Market,
			item.Sector,
			item.SubSector,
			fmt.Sprintf("%f", item.Price),
			item.FiftyTwoWeek,
			fmt.Sprintf("%d", item.Volume),
			fmt.Sprintf("%f", item.EPS),
			fmt.Sprintf("%f", item.DPS),
			fmt.Sprintf("%f", item.NTA),
			fmt.Sprintf("%f", item.PE),
			fmt.Sprintf("%f", item.DY),
			fmt.Sprintf("%f", item.ROE),
			fmt.Sprintf("%f", item.PTBV),
			fmt.Sprintf("%f", item.MarketCapInMillion),
			item.Summary,
			item.Website,
		})
	}
	fmt.Printf("Exported combined stock listings with company summaries to %s\n", filename)
}

type StockListingWithSummary struct {
	*klsescreener.StockListing
	*klsescreener.CompanySummary
}

func combineStockListingsAndCompanySummaries() ([]*StockListingWithSummary, error) {
	// Get stock listings
	ctx := context.Background()
	stockListings, err := klsescreener.GetStockListings(ctx)
	if err != nil {
		return nil, err
	}
	// Get company summaries for each stock listing
	var stockListingsWithSummaries []*StockListingWithSummary
	for _, listing := range stockListings {
		summary, err := klsescreener.GetCompanySummary(ctx, listing.Code)
		if err != nil {
			// If there is an error getting the company summary, we can still add the stock listing without the summary
			stockListingsWithSummaries = append(stockListingsWithSummaries, &StockListingWithSummary{
				StockListing:   listing,
				CompanySummary: &klsescreener.CompanySummary{},
			})
			continue
		}
		stockListingsWithSummaries = append(stockListingsWithSummaries, &StockListingWithSummary{
			StockListing:   listing,
			CompanySummary: summary,
		})
	}
	// Print combined stock listings with company summaries
	for _, item := range stockListingsWithSummaries {
		println(item.Code, item.Name, item.Market, item.Sector, item.Summary, item.Website)
	}
	return stockListingsWithSummaries, nil
}
