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
		Symbol              string      `json:"symbol"`
		High52              json.Number `json:"high52"`
		Low52               json.Number `json:"low52"`
		DividendAmount      json.Number `json:"dividendAmount"`
		DividendYield       json.Number `json:"dividendYield"`
		DividendDate        string      `json:"dividendDate"`
		PeRatio             json.Number `json:"peRatio"`
		PegRatio            json.Number `json:"pegRatio"`
		PbRatio             json.Number `json:"pbRatio"`
		PrRatio             json.Number `json:"prRatio"`
		PcfRatio            json.Number `json:"pcfRatio"`
		GrossMarginTTM      json.Number `json:"grossMarginTTM"`
		GrossMarginMRQ      json.Number `json:"grossMarginMRQ"`
		NetProfitMarginTTM  json.Number `json:"netProfitMarginTTM"`
		NetProfitMarginMRQ  json.Number `json:"netProfitMarginMRQ"`
		OperatingMarginTTM  json.Number `json:"operatingMarginTTM"`
		OperatingMarginMRQ  json.Number `json:"operatingMarginMRQ"`
		ReturnOnEquity      json.Number `json:"returnOnEquity"`
		ReturnOnAssets      json.Number `json:"returnOnAssets"`
		ReturnOnInvestment  json.Number `json:"returnOnInvestment"`
		QuickRatio          json.Number `json:"quickRatio"`
		CurrentRatio        json.Number `json:"currentRatio"`
		InterestCoverage    json.Number `json:"interestCoverage"`
		TotalDebtToCapital  json.Number `json:"totalDebtToCapital"`
		LtDebtToEquity      json.Number `json:"ltDebtToEquity"`
		TotalDebtToEquity   json.Number `json:"totalDebtToEquity"`
		EpsTTM              json.Number `json:"epsTTM"`
		EpsChangePercentTTM json.Number `json:"epsChangePercentTTM"`
		EpsChangeYear       json.Number `json:"epsChangeYear"`
		EpsChange           json.Number `json:"epsChange"`
		RevChangeYear       json.Number `json:"revChangeYear"`
		RevChangeTTM        json.Number `json:"revChangeTTM"`
		RevChangeIn         json.Number `json:"revChangeIn"`
		SharesOutstanding   json.Number `json:"sharesOutstanding"`
		MarketCapFloat      json.Number `json:"marketCapFloat"`
		MarketCap           json.Number `json:"marketCap"`
		BookValuePerShare   json.Number `json:"bookValuePerShare"`
		ShortIntToFloat     json.Number `json:"shortIntToFloat"`
		ShortIntDayToCover  json.Number `json:"shortIntDayToCover"`
		DivGrowthRate3Year  json.Number `json:"divGrowthRate3Year"`
		DividendPayAmount   json.Number `json:"dividendPayAmount"`
		DividendPayDate     string      `json:"dividendPayDate"`
		Beta                json.Number `json:"beta"`
		Vol1DayAvg          json.Number `json:"vol1DayAvg"`
		Vol10DayAvg         json.Number `json:"vol10DayAvg"`
		Vol3MonthAvg        json.Number `json:"vol3MonthAvg"`
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
