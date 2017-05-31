package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL     = "https://api.hello.is"
	authURL     = baseURL + "/v1/oauth2/token"
	devicesURL  = baseURL + "/v2/devices"
	roomInfoURL = baseURL + "/v1/room/current?temp_unit=c"

	clientID     = "8d3c1664-05ae-47e4-bcdb-477489590aa4"
	clientSecret = "4f771f6f-5c10-4104-bbc6-3333f5b11bf9"
)

var (
	// ErrWrongCredentials is used when the credentials are not correct.
	ErrWrongCredentials = errors.New("wrong credentials")
)

// HelloClient provides a client interface to the Hello API.
type HelloClient struct {
	username string
	password string
	token    authToken
	client   *http.Client
}

type authToken struct {
	Token   string
	Expires time.Time
}

// NewClient creates a new client instance.
func NewClient(username, password string) *HelloClient {
	return &HelloClient{
		username: username,
		password: password,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Devices returns a list of devices connected to the account.
func (c *HelloClient) Devices() (result Devices, err error) {
	err = c.checkAuthentication()
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest(http.MethodGet, devicesURL, nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token.Token)

	res, err := c.client.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// RoomInfo gets the current room information from the API.
func (c *HelloClient) RoomInfo() (result RoomInfo, err error) {
	err = c.checkAuthentication()
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest(http.MethodGet, roomInfoURL, nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token.Token)

	res, err := c.client.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (c *HelloClient) checkAuthentication() error {
	if len(c.token.Token) > 0 && c.token.Expires.After(time.Now()) {
		return nil
	}

	data := url.Values{
		"grant_type":    []string{"password"},
		"client_id":     []string{clientID},
		"client_secret": []string{clientSecret},
		"username":      []string{c.username},
		"password":      []string{c.password},
	}
	res, err := c.client.PostForm(authURL, data)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return ErrWrongCredentials
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var token TokenInfo
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return err
	}

	c.token.Expires = time.Now().Add(time.Duration(token.ExpiresIn) * time.Millisecond)
	c.token.Token = token.AccessToken

	return nil
}
