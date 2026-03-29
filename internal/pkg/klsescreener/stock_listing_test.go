package klsescreener_test

import (
	"context"
	"testing"

	"github.com/kokweikhong/goklse/internal/pkg/klsescreener"
	"github.com/stretchr/testify/assert"
)

func TestGetStockListings(t *testing.T) {
	ctx := context.Background()
	stockListings, err := klsescreener.GetStockListings(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, stockListings)
	// Print the first 5 stock listings for verification
	for i, listing := range stockListings {
		if i >= 5 {
			break
		}
		t.Logf("Stock %d: Code=%s, Name=%s, Market=%s, Sector=%s", i+1, listing.Code, listing.Name, listing.Market, listing.Sector)
	}
}