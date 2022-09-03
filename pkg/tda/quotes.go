package tda

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

// Quote is the response for a quote.
type Quote struct {
	AssetType                          string  `json:"assetType"`
	AssetMainType                      string  `json:"assetMainType"`
	Cusip                              string  `json:"cusip"`
	AssetSubType                       string  `json:"assetSubType"`
	Symbol                             string  `json:"symbol"`
	Description                        string  `json:"description"`
	BidPrice                           float64 `json:"bidPrice"`
	BidSize                            int     `json:"bidSize"`
	BidID                              string  `json:"bidId"`
	AskPrice                           float64 `json:"askPrice"`
	AskSize                            int     `json:"askSize"`
	AskID                              string  `json:"askId"`
	LastPrice                          float64 `json:"lastPrice"`
	LastSize                           int     `json:"lastSize"`
	LastID                             string  `json:"lastId"`
	OpenPrice                          float64 `json:"openPrice"`
	HighPrice                          float64 `json:"highPrice"`
	LowPrice                           float64 `json:"lowPrice"`
	BidTick                            string  `json:"bidTick"`
	ClosePrice                         float64 `json:"closePrice"`
	NetChange                          float64 `json:"netChange"`
	TotalVolume                        int     `json:"totalVolume"`
	QuoteTimeInLong                    int64   `json:"quoteTimeInLong"`
	TradeTimeInLong                    int64   `json:"tradeTimeInLong"`
	Mark                               float64 `json:"mark"`
	Exchange                           string  `json:"exchange"`
	ExchangeName                       string  `json:"exchangeName"`
	Marginable                         bool    `json:"marginable"`
	Shortable                          bool    `json:"shortable"`
	Volatility                         float64 `json:"volatility"`
	Digits                             int     `json:"digits"`
	Five2WkHigh                        float64 `json:"52WkHigh"`
	Five2WkLow                         float64 `json:"52WkLow"`
	NAV                                float64 `json:"nAV"`
	PeRatio                            float64 `json:"peRatio"`
	DivAmount                          float64 `json:"divAmount"`
	DivYield                           float64 `json:"divYield"`
	DivDate                            string  `json:"divDate"`
	SecurityStatus                     string  `json:"securityStatus"`
	RegularMarketLastPrice             float64 `json:"regularMarketLastPrice"`
	RegularMarketLastSize              int     `json:"regularMarketLastSize"`
	RegularMarketNetChange             float64 `json:"regularMarketNetChange"`
	RegularMarketTradeTimeInLong       int64   `json:"regularMarketTradeTimeInLong"`
	NetPercentChangeInDouble           float64 `json:"netPercentChangeInDouble"`
	MarkChangeInDouble                 float64 `json:"markChangeInDouble"`
	MarkPercentChangeInDouble          float64 `json:"markPercentChangeInDouble"`
	RegularMarketPercentChangeInDouble float64 `json:"regularMarketPercentChangeInDouble"`
	Delayed                            bool    `json:"delayed"`
	RealtimeEntitled                   bool    `json:"realtimeEntitled"`
}

type QuoteRequest struct {
	Symbols string `url:"symbol"`
}

// GetQuote accesses the TDAmeritrade API using an existing Session struct to
// provide quote data for a single ticker. If you are quoting more than one
// security, use GetQuotes()
func (s *Session) GetQuote(ticker string) (*Quote, error) {
	res, err := s.GetQuotes([]string{ticker})
	if err != nil {
		return nil, err
	}

	val := (*res)[strings.ToUpper(ticker)]
	return &val, nil
}

// GetQuotes accesses the TDAmeritrade API using an existing Session struct
// to provide quote data for multiple tickers.
func (s *Session) GetQuotes(tickers []string) (*map[string]Quote, error) {
	token, err := s.GetAccessToken()
	if err != nil {
		return nil, err
	}

	v, err := query.Values(QuoteRequest{
		Symbols: strings.Join(tickers, ","),
	})
	if err != nil {
		return nil, fmt.Errorf("could not querystring tickets: %w", err)
	}

	url := fmt.Sprintf("%s/marketdata/quotes?%s", s.RootUrl, v.Encode())
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := s.httpClient.Do(req)

	if err = getHttpError(res); err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var quotes map[string]Quote
	if err := json.Unmarshal(body, &quotes); err != nil {
		return nil, fmt.Errorf("could not parse quotes output: %w", err)
	}

	return &quotes, nil
}
