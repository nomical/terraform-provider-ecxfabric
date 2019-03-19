package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseURL                    = "https://api.equinix.com"
	authEndpoint               = "/oauth2/v1/token"
	ecxV3L2ConnectionsEndpoint = "/ecx/v3/l2/connections"
)

type Option func(*Client) error

func BaseURL(baseURL string) Option {
	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}

func (c *Client) parseOptions(opts ...Option) error {
	for _, option := range opts {
		err := option(c)
		if err != nil {
			return err
		}
	}

	return nil
}

type Client struct {
	baseURL      string
	httpClient   *http.Client
	clientID     string
	clientSecret string
	username     string
	password     string
	accessToken  string
}

func New(opts ...Option) (*Client, error) {
	client := &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}

	if err := client.parseOptions(opts...); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Authenticate(clientID, clientSecret, username, password string) error {
	d, err := json.Marshal(oAuthRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Username:     username,
		Password:     password,
		GrantType:    "password",
	})
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(d)

	url := fmt.Sprintf("%s%s", c.baseURL, authEndpoint)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var oAuthErrorResponse oAuthErrorResponse
	var oAuthResponse oAuthResponse
	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(&oAuthResponse); err != nil {
			return err
		}
		c.accessToken = oAuthResponse.AccessToken

		return nil
	case http.StatusUnauthorized, http.StatusInternalServerError:
		if err := json.NewDecoder(resp.Body).Decode(&oAuthErrorResponse); err != nil {
			return err
		}
	}

	return &AuthError{
		resp.StatusCode,
		oAuthErrorResponse,
	}
}

func (c *Client) CreateL2Connection(input PostConnectionRequest) (*PostConnectionResponse, error) {
	d, err := json.MarshalIndent(input, "", "    ")
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(d)

	url := fmt.Sprintf("%s%s", c.baseURL, ecxV3L2ConnectionsEndpoint)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.accessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var postConnectionResponse PostConnectionResponse
	var errorResponse errorResponse
	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(&postConnectionResponse); err != nil {
			return nil, err
		}

		return &postConnectionResponse, nil
	case http.StatusUnauthorized:
		var oAuthErrorResponse oAuthErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&oAuthErrorResponse); err != nil {
			return nil, err
		}

		return nil, &AuthError{
			resp.StatusCode,
			oAuthErrorResponse,
		}
	case http.StatusBadRequest:
		var errorResponseArray errorResponseArray
		if err := json.NewDecoder(resp.Body).Decode(&errorResponseArray); err != nil {
			return nil, err
		}

		return nil, &ErrorArray{
			resp.StatusCode,
			errorResponseArray,
		}
	default:
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}
	}

	return nil, &Error{
		resp.StatusCode,
		errorResponse,
	}
}

func (c *Client) ReadL2Connection(uuid string) (*GetConnectionByUUIDResponse, error) {
	url := fmt.Sprintf("%s%s/%s", c.baseURL, ecxV3L2ConnectionsEndpoint, uuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.accessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var getConnectionByUUIDResponse GetConnectionByUUIDResponse
	var errorResponse errorResponse
	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(&getConnectionByUUIDResponse); err != nil {
			return nil, err
		}

		return &getConnectionByUUIDResponse, nil
	case http.StatusUnauthorized:
		var oAuthErrorResponse oAuthErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&oAuthErrorResponse); err != nil {
			return nil, err
		}

		return nil, &AuthError{
			resp.StatusCode,
			oAuthErrorResponse,
		}
	case http.StatusBadRequest:
		var errorResponseArray errorResponseArray
		if err := json.NewDecoder(resp.Body).Decode(&errorResponseArray); err != nil {
			return nil, err
		}

		return nil, &ErrorArray{
			resp.StatusCode,
			errorResponseArray,
		}
	default:
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}
	}

	return nil, &Error{
		resp.StatusCode,
		errorResponse,
	}
}

func (c *Client) DeleteL2Connection(uuid string) error {
	url := fmt.Sprintf("%s%s/%s", c.baseURL, ecxV3L2ConnectionsEndpoint, uuid)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.accessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var errorResponse errorResponse
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		var oAuthErrorResponse oAuthErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&oAuthErrorResponse); err != nil {
			return err
		}

		return &AuthError{
			resp.StatusCode,
			oAuthErrorResponse,
		}
	case http.StatusBadRequest:
		var errorResponseArray errorResponseArray
		if err := json.NewDecoder(resp.Body).Decode(&errorResponseArray); err != nil {
			return err
		}

		return &ErrorArray{
			resp.StatusCode,
			errorResponseArray,
		}
	default:
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return err
		}
	}

	return &Error{
		resp.StatusCode,
		errorResponse,
	}
}
