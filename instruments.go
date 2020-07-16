package tda

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var FundamentalsEmpty = errors.New("tda: empty fundamentals received")

type InstrumentFundamentals struct {
	Fundamental struct {
		Symbol              string  `json:"symbol"`
		High52              float64 `json:"high52"`
		Low52               float64 `json:"low52"`
		DividendAmount      float64 `json:"dividendAmount"`
		DividendYield       float64 `json:"dividendYield"`
		DividendDate        string  `json:"dividendDate"`
		PeRatio             float64 `json:"peRatio"`
		PegRatio            float64 `json:"pegRatio"`
		PbRatio             float64 `json:"pbRatio"`
		PrRatio             float64 `json:"prRatio"`
		PcfRatio            float64 `json:"pcfRatio"`
		GrossMarginTTM      float64 `json:"grossMarginTTM"`
		GrossMarginMRQ      float64 `json:"grossMarginMRQ"`
		NetProfitMarginTTM  float64 `json:"netProfitMarginTTM"`
		NetProfitMarginMRQ  float64 `json:"netProfitMarginMRQ"`
		OperatingMarginTTM  float64 `json:"operatingMarginTTM"`
		OperatingMarginMRQ  float64 `json:"operatingMarginMRQ"`
		ReturnOnEquity      float64 `json:"returnOnEquity"`
		ReturnOnAssets      float64 `json:"returnOnAssets"`
		ReturnOnInvestment  float64 `json:"returnOnInvestment"`
		QuickRatio          float64 `json:"quickRatio"`
		CurrentRatio        float64 `json:"currentRatio"`
		InterestCoverage    float64 `json:"interestCoverage"`
		TotalDebtToCapital  float64 `json:"totalDebtToCapital"`
		LtDebtToEquity      float64 `json:"ltDebtToEquity"`
		TotalDebtToEquity   float64 `json:"totalDebtToEquity"`
		EpsTTM              float64 `json:"epsTTM"`
		EpsChangePercentTTM float64 `json:"epsChangePercentTTM"`
		EpsChangeYear       float64 `json:"epsChangeYear"`
		EpsChange           float64 `json:"epsChange"`
		RevChangeYear       float64 `json:"revChangeYear"`
		RevChangeTTM        float64 `json:"revChangeTTM"`
		RevChangeIn         float64 `json:"revChangeIn"`
		SharesOutstanding   int64   `json:"sharesOutstanding"`
		MarketCapFloat      float64 `json:"marketCapFloat"`
		MarketCap           float64 `json:"marketCap"`
		BookValuePerShare   float64 `json:"bookValuePerShare"`
		ShortIntToFloat     float64 `json:"shortIntToFloat"`
		ShortIntDayToCover  float64 `json:"shortIntDayToCover"`
		DivGrowthRate3Year  float64 `json:"divGrowthRate3Year"`
		DividendPayAmount   float64 `json:"dividendPayAmount"`
		DividendPayDate     string  `json:"dividendPayDate"`
		Beta                float64 `json:"beta"`
		Vol1DayAvg          int     `json:"vol1DayAvg"`
		Vol10DayAvg         int     `json:"vol10DayAvg"`
		Vol3MonthAvg        int     `json:"vol3MonthAvg"`
	} `json:"fundamental"`
	Cusip       string `json:"cusip"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Exchange    string `json:"exchange"`
	AssetType   string `json:"assetType"`
}

func (s *Session) GetInstrumentFundamentals(ticker string) (*InstrumentFundamentals, error) {
	token, err := s.GetAccessToken()
	if err != nil {
		return nil, &ApiError{
			Reason: "Could not authenticate with TDAmeritrade",
			Err:    errors.New("GetInstrumentFundamentals() GetAccessToken"),
		}
	}

	url := fmt.Sprintf("%s/instruments?symbol=%s&projection=fundamental", s.RootUrl, ticker)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := s.HttpClient.Do(req)

	if err = getHttpError(res); err != nil {
		return nil, &ApiError{
			Reason: "An error occured with retrieving fundamentals",
			Err:    errors.New("GetInstrumentFundamentals() getHttpError"),
		}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, &ApiError{
			Reason: "An error occured reading the respose",
			Err:    errors.New("GetInstrumentFundamentals() ioutil.ReadAll"),
		}
	}

	// check for empty body
	if string(body) == "{}" || string(body) == "" {
		return nil, FundamentalsEmpty
	}

	var fundamentalsDataT map[string]InstrumentFundamentals
	json.Unmarshal(body, &fundamentalsDataT)
	fundamentalsData := fundamentalsDataT[ticker]
	return &fundamentalsData, nil
}
