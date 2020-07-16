package tda

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

type ApiError struct {
	Reason string
	Err    error
}

func (r *ApiError) Error() string {
	return fmt.Sprintf("reason %s: err %v", r.Reason, r.Err)
}

type Session struct {
	Refresh     string
	ConsumerKey string
	RootUrl     string
	HttpClient  http.Client
}

type AccessTokenStruct struct {
	GrantType    string `url:"grant_type"`
	RefreshToken string `url:"refresh_token"`
	ClientID     string `url:"client_id"`
	RedirectUri  string `url:"redirect_uri"`
}

func (s *Session) InitSession() {
	s.HttpClient = http.Client{Timeout: time.Second * 10}
}

func getHttpError(res *http.Response) error {

	switch res.StatusCode {
	case 200:
		return nil
	case 400:
		return &ApiError{
			Reason: "Bad request was made, are you sure your keys are properly typed?",
			Err:    errors.New("400"),
		}
	case 401:
		return &ApiError{
			Reason: "Invalid refresh token or consumer key, try again",
			Err:    errors.New("401"),
		}

	default:
		return &ApiError{
			Reason: "Server error handling your request",
			Err:    errors.New(fmt.Sprintf("%d", res.StatusCode)),
		}
	}
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (s *Session) GetAccessToken() (string, error) {
	payload := &AccessTokenStruct{
		GrantType:    "refresh_token",
		RefreshToken: s.Refresh,
		ClientID:     fmt.Sprintf("%s@AMER.OAUTHMAP", s.ConsumerKey),
		RedirectUri:  "http://127.0.0.1",
	}

	url := fmt.Sprintf("%s/oauth2/token", s.RootUrl)

	v, err := query.Values(payload)
	if err != nil {
		return "", &ApiError{
			Reason: "Could not querystring your payload",
			Err:    err,
		}
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return "", &ApiError{
			Reason: "Could not validate request",
			Err:    errors.New("GetAccessToken() http.newRequest cannot be validated"),
		}
	}

	res, err := s.HttpClient.Do(req)

	if err != nil {
		return "", &ApiError{
			Reason: "Could not process request",
			Err:    err,
		}
	}

	if err = getHttpError(res); err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", &ApiError{
			Reason: "Could not read TDA Response",
			Err:    err,
		}
	}

	var tokenResponse AccessTokenResponse
	json.Unmarshal(body, &tokenResponse)

	return tokenResponse.AccessToken, nil
}
