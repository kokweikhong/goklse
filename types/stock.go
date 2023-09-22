package types

// StockQuote is the struct for stock quote
type StockQuote struct {
	Name                string  `json:"name"`
	ShortName           string  `json:"shortName"`
	Code                string  `json:"code"`
	Category            string  `json:"category"`
	SubCategory         string  `json:"subCategory"`
	Market              string  `json:"market"`
	Price               float64 `json:"price"`
	ChangePercent       float64 `json:"changePercent"`
	FiftyTwoWeekLow     float64 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh    float64 `json:"fiftyTwoWeekHigh"`
	Volume              float64 `json:"volume"`
	EPS                 float64 `json:"eps"`
	DPS                 float64 `json:"dps"`
	NTA                 float64 `json:"nta"`
	PE                  float64 `json:"pe"`
	PTBV                float64 `json:"ptbv"`
	DY                  float64 `json:"dy"`
	ROE                 float64 `json:"roe"`
	MarketCapInMillions float64 `json:"marketCapInMillions"`
}

// StockQuoteWithStatistical is the struct for stock quote with statistical data
type StockQuoteWithStatistical struct {
	StockQuote
	HistoricalPrice []*OHLC `json:"historicalPrice"`
}
