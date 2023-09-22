package main

import (
	"log"

	"github.com/kokweikhong/goklse/klse"
)

func main() {
	// klse.GetStockListing()
	resp, err := klse.GetBursaMarketHistoricalData()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range resp {
		log.Println(v)
	}

}
