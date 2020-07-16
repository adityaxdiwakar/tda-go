package tda

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ApiError struct {
	Reason string
	Err    error
}

func (r *ApiError) Error() string {
	return fmt.Sprintf("reason %d: err %v", r.Reason, r.Err)
}

type Session struct {
	Refresh     string
	ConsumerKey string
	RootUrl     string
	HttpClient  *http.Client
}

type AccessTokenStruct struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	RedirectUri  string `json:"redirect_uri"`
}

func (s *Session) InitSession() {
	s.HttpClient = &http.Client{Timeout: time.Second * 10}
}

func (s *Session) GetAccessToken() (string, error) {
	payload := &AccessTokenStruct{
		GrantType:    "refresh_token",
		RefreshToken: s.Refresh,
		ClientID:     fmt.Sprintf("%s@AMER.OAUTHMAP", s.ConsumerKey),
		RedirectUri:  "http://127.0.0.1",
	}

	url := fmt.Sprintf("%s/oauth2/token", s.RootUrl)

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(payload)
	req, err := http.NewRequest("POST", url, buf)

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

}
