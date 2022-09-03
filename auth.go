package tda

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

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

	// httpClient is the library specific client, defaulted to initialize with
	// a 10 second limit for requests, in case of errors on TDA servers to
	// prevent stalling
	httpClient http.Client

	// accessStatPath is a path that can be used to determine when the last
	// access token was created and also to receive the value. Since this
	// access token can be used for sensitive use cases, it is recommended
	// that this is set safely; can only be set, cannot be changed.
	accessStatPath string
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

// SessionOption is an option that can be provided to NewSession in order to
// modify the internal state of the created Session.
type SessionOption func(*Session)

// WithHttpClient is an option that returns a SessionOption that can be used
// to change the internal HTTP Client used by the TDA Session.
func WithHttpClient(client http.Client) SessionOption {
	return func(s *Session) {
		s.httpClient = client
	}
}

// WithStatPath is an option that sets the stat path to prevent creating
// a new access token on every request. This is recommended to be used if
// server is frequently making requests.
func WithStatPath(path string) SessionOption {
	return func(s *Session) {
		s.accessStatPath = path
	}
}

// NewSession constructs a new TDA Go Session taking in options.
func NewSession(refresh, consumerKey, rootUrl string, opts ...SessionOption) *Session {
	s := &Session{
		Refresh:     refresh,
		ConsumerKey: consumerKey,
		RootUrl:     rootUrl,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// getHttpError is an internal error handler function to allow for easier error
// handling with respondes from the TDAmeritrade API
func getHttpError(res *http.Response) error {
	switch res.StatusCode {
	case 200:
		return nil
	case 400:
		return errors.New("Bad request was made, are you sure your keys are properly typed?")
	case 401:
		return errors.New("Invalid refresh token or consumer key, try again")
	default:
		return fmt.Errorf("Server could not handle request, returned status: %d", res.StatusCode)
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

	if s.accessStatPath != "" {
		// First, check if there is a stat path set to see if token
		// generation is required.
		info, err := os.Stat(s.accessStatPath)
		if errors.Is(err, os.ErrNotExist) {
			file, err := os.Create(s.accessStatPath)
			file.Close()
			if err != nil {
				return "", fmt.Errorf("could not create storage in stat path")
			}

			info, err = os.Stat(s.accessStatPath)
			if err != nil {
				return "", fmt.Errorf("could not stat recently created stat file")
			}
		} else if err != nil {
			return "", fmt.Errorf("could not stat access path")
		}

		fileTime := info.ModTime()
		dat, err := os.ReadFile(s.accessStatPath)
		if err != nil {
			return "", fmt.Errorf("could not open access token file: %w", err)
		}

		if !(fileTime.Before(time.Now().Add(-25*time.Minute)) || string(dat) == "") {
			// token exists and is valid, no need to generate new token
			return string(dat), nil
		}
	}

	// generate new token
	payload := &AccessTokenStruct{
		GrantType:    "refresh_token",
		RefreshToken: s.Refresh,
		ClientID:     fmt.Sprintf("%s@AMER.OAUTHMAP", s.ConsumerKey),
		RedirectUri:  "http://127.0.0.1",
	}

	url := fmt.Sprintf("%s/oauth2/token", s.RootUrl)

	v, err := query.Values(payload)
	if err != nil {
		return "", fmt.Errorf("could not querystring keys: %w", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return "", fmt.Errorf("could not validate request")
	}

	res, err := s.httpClient.Do(req)

	if err != nil {
		return "", fmt.Errorf("could not process request: %w", err)
	}

	if err = getHttpError(res); err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("could not read TDA Response: %w", err)
	}

	var tokenResponse AccessTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", fmt.Errorf("could not response into token response struct: %w", err)
	}

	return tokenResponse.AccessToken, nil
}
