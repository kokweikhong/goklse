package main

import (
	"context"

	"github.com/kokweikhong/goklse/internal/pkg/klsescreener"
)

func main() {
	ctx := context.Background()
	indexes, err := klsescreener.GetMarketIndexes(ctx)
	if err != nil {
		panic(err)
	}

	for _, index := range indexes {
		println(index.Name, index.Code, index.Country)
	}

}
