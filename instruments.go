package tda

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// InstrumentFundmanetals is the struct for a valid fundamentals response from TDA
type InstrumentFundamentals struct {
	// Fundamental struct containing FA indicators
	Fundamental struct {
		// Symbol for the listed security on NYSE (and/or NASDAAQ)
		Symbol string `json:"symbol"`

		// 52Week high of the stock
		High52 json.Number `json:"high52"`

		// 52Week low of the stock
		Low52 json.Number `json:"low52"`

		// DividendAmount is a numerical value (in dollars) that a company will
		// pay out to its shareholders
		DividendAmount json.Number `json:"dividendAmount"`

		// DividendYield is the annualzied dividend amount divided by the price
		// of the stock, lower share price with a high dividend means a very
		// high yield
		DividendYield json.Number `json:"dividendYield"`

		// DividendDate is the expected ex-dividend date based on TDA's
		// calendar, this is not always accurate
		DividendDate string `json:"dividendDate"`

		// Price to Earnings Ratio
		PeRatio json.Number `json:"peRatio"`

		// Price to Earnings-Growth Ratio
		PegRatio json.Number `json:"pegRatio"`

		// Price to Book Ratio
		PbRatio json.Number `json:"pbRatio"`

		// TODO: What is a PRRatio?
		PrRatio json.Number `json:"prRatio"`

		// Price to Cash-Flow Ratio
		PcfRatio json.Number `json:"pcfRatio"`

		// Gross Margin (borrowed cash) in the Trailing Twelve Months
		GrossMarginTTM json.Number `json:"grossMarginTTM"`

		// Gross Margin (borrowed cash) in the Most Recent Quarter
		GrossMarginMRQ json.Number `json:"grossMarginMRQ"`

		// Net Profit Margin (Percantage of Revenue) in Trailing Twelve Months
		NetProfitMarginTTM json.Number `json:"netProfitMarginTTM"`

		// Net Profit Margin (Percantage of Revenue) in Most Recent Quarter
		NetProfitMarginMRQ json.Number `json:"netProfitMarginMRQ"`

		// Operating Margin Trailing Twelve Months
		OperatingMarginTTM json.Number `json:"operatingMarginTTM"`

		// Operating Margin Most Recent Quarter
		OperatingMarginMRQ json.Number `json:"operatingMarginMRQ"`

		// Net Income divided by Shareholder Equity
		ReturnOnEquity json.Number `json:"returnOnEquity"`

		// Net Income divided by Net Company Assets
		ReturnOnAssets json.Number `json:"returnOnAssets"`

		// Net Income divided by Investments
		ReturnOnInvestment json.Number `json:"returnOnInvestment"`

		// Ability of a company to use near cash to extinguish long term
		// maturities
		QuickRatio json.Number `json:"quickRatio"`

		// Ability of a company to use near cash to meet short-term obligations
		CurrentRatio json.Number `json:"currentRatio"`

		// Ability of a company to pay off interest expenses on outstanding debt
		InterestCoverage json.Number `json:"interestCoverage"`

		// Company Debt compared to Current Company Capital
		TotalDebtToCapital json.Number `json:"totalDebtToCapital"`

		// Long Term Debts to Equity Ratio
		LtDebtToEquity json.Number `json:"ltDebtToEquity"`

		// All Debt to Equity Ratio
		TotalDebtToEquity json.Number `json:"totalDebtToEquity"`

		// Earnings-Per-Share Trailing Twelve Months
		EpsTTM json.Number `json:"epsTTM"`

		// Earnings-Per-Share Change % (Growth/Decline) Twailing Twelve Months
		EpsChangePercentTTM json.Number `json:"epsChangePercentTTM"`

		// Earnings-Per-Share Change over 1-Year (Annual Filing)
		EpsChangeYear json.Number `json:"epsChangeYear"`

		// Earnings-Per-Share Change QoQ (Quarterly Filing)
		EpsChange json.Number `json:"epsChange"`

		// Revenue Change YoY (Annual Filing)
		RevChangeYear json.Number `json:"revChangeYear"`

		// Revenue Change Trailing Twelve Months
		RevChangeTTM json.Number `json:"revChangeTTM"`

		// Incoming Revenue Change
		RevChangeIn json.Number `json:"revChangeIn"`

		// Number of shares currently availabile in the market
		// (NumberOfShares * SharePrice = MarketCapitalization)
		SharesOutstanding json.Number `json:"sharesOutstanding"`

		// $ Value of Available Shares Traded on Exchanges
		MarketCapFloat json.Number `json:"marketCapFloat"`

		// Market Capitalization
		MarketCap json.Number `json:"marketCap"`

		// Value of Share based on Books
		BookValuePerShare json.Number `json:"bookValuePerShare"`

		// % of Short Interest Compared to Available Shares
		ShortIntToFloat json.Number `json:"shortIntToFloat"`

		// # of Days before Company Closes out Shorted Shares
		ShortIntDayToCover json.Number `json:"shortIntDayToCover"`

		// Dividend Growth Rate over 3 Years
		DivGrowthRate3Year json.Number `json:"divGrowthRate3Year"`

		// Dividend Pay Amount (QoQ or MoM Basis)
		DividendPayAmount json.Number `json:"dividendPayAmount"`

		// Pay Date (as opposed to Ex-Div Date) for Dividend
		DividendPayDate string `json:"dividendPayDate"`

		// Beta Correlation to SPX (Broad Market)
		Beta json.Number `json:"beta"`

		// Average Volume based on 1 Day of Trading
		Vol1DayAvg json.Number `json:"vol1DayAvg"`

		// Average Volume based on 10 Days of Trading
		Vol10DayAvg json.Number `json:"vol10DayAvg"`

		// Average Volume based on 90 days of Trading
		Vol3MonthAvg json.Number `json:"vol3MonthAvg"`
	} `json:"fundamental"`
	// Committee on Uniform Securities Identification Procedures
	Cusip string `json:"cusip"`

	// Symbol for Security
	Symbol string `json:"symbol"`

	// Description of Security
	Description string `json:"description"`

	// Exchange the Security trades on
	Exchange string `json:"exchange"`

	// Type of Asset the Security is
	AssetType string `json:"assetType"`
}

// GetInstrumentFundmamentals accesses the TDAmeritrade API using an existing
// Session struct to provide fundamental data in the form of the
// InstrumentFundamentals Struct. If the payload from the TDAmeritrade API is
// empty, an error is returned (rather than only an empty struct) with the
// fundamentalsEmpty error. The only input parameter is the relevant ticker,
// which is not case-sensitive.
func (s *Session) GetInstrumentFundamentals(ticker string) (*InstrumentFundamentals, error) {
	token, err := s.GetAccessToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/instruments?symbol=%s&projection=fundamental", s.RootUrl, ticker)
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

	var fundamentalsDataT map[string]InstrumentFundamentals
	if err := json.Unmarshal(body, &fundamentalsDataT); err != nil {
		return nil, fmt.Errorf("could not parse fundamentals output: %w", err)
	}

	fundamentalsData := fundamentalsDataT[strings.ToUpper(ticker)]

	return &fundamentalsData, nil
}
