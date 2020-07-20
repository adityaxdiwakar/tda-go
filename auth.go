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

// Session is the client structure that is used for all relevant parts of this
// library. The session struct takes in 3 required parameters which is the
// refresh key, consumer key, and the root url (which is usually the same for
// all endusers).
type Session struct {
	// RefreshKey of your account given by the TDAmeritrade OAuth2 Workflow
	Refresh string

	// ConsumerKey from the TDAmeritrade Developer Application Portal
	ConsumerKey string

	// RootUrl for the API, this will usually be https://api.tdameritrade.com/v1
	RootUrl string

	// HttpClient is the library specific client, defaulted to initialize with
	// a 10 second limit for requests, in case of errors on TDA servers to
	// prevent stalling
	HttpClient http.Client
}

// AccessTokenStruct is the internal structure used for establishing the
// request between TDAmeritrade Servers and the library, this is used to create
// a URLEncoded Form object which is sent to retrieve a user's access token
type AccessTokenStruct struct {
	// GrantType is the type of grant for the access token, this will always be
	// 'refresh_token' for this library
	GrantType string `url:"grant_type"`

	// RefreshToken provided by the Session intialization struct
	RefreshToken string `url:"refresh_token"`

	// ClientID is the CONSUMER_KEY with '@AMER.OAUTHAPP' added
	ClientID string `url:"client_id"`

	// RedirectUri of the app from TDAmeritrade's Development Portal
	RedirectUri string `url:"redirect_uri"`
}

// InitSession is to be called after initializing the Session struct in order
// to create the root library HTTP client with a default timeout of 10 seconds
func (s *Session) InitSession() {
	s.HttpClient = http.Client{Timeout: time.Second * 10}
}

// getHttpError is an internal error handler function to allow for easier error
// handling with respondes from the TDAmeritrade API
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

// AccessTokenResponse is a struct to load the TDAmeritrade Response into, this
// takes many fields and the user is returned the AccessToken using
// *Session.GetAccessToken()
type AccessTokenResponse struct {
	// AccessToken provided from the endpoint to be used to make subsequent
	// requests
	AccessToken string `json:"access_token"`

	// Scope for what endpoints are allowed, for properly configured apps, this
	// will be all scopes
	Scope string `json:"scope"`

	// Expiry for the access token, all access tokens expire after 1800s (30m)
	// by default
	ExpiresIn int `json:"expires_in"`

	// TokenType is 'Bearer' for the TDA API
	TokenType string `json:"token_type"`
}

// GetAccessToken is the session function to retrieve the access token, this
// returns a multitude of errors and is only intended for internal library
// usage but is exposed in case the AccessToken is important to be accessed
// externally
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
