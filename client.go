package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// constants
const (
	SuccessStatusCode = 0
)

// Client api client
type Client struct {
	client *http.Client
	config *Config
	logger *log.Logger
}

// Config api client config
type Config struct {
	BaseEndpoint string
	APIKey       string
	APISecret    string
	Debug        bool
}

// NewClient creates api client
func NewClient(cfg *Config, c *http.Client, logger *log.Logger) *Client {
	if logger == nil {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}
	return &Client{
		client: c,
		config: cfg,
		logger: logger,
	}
}

func (c *Client) createSig(nonce int64, url, payload string) string {
	sig := hmac.New(sha256.New, []byte(c.config.APISecret))
	message := fmt.Sprintf("%d%s%s", nonce, url, payload)
	sig.Write([]byte(message))
	return hex.EncodeToString(sig.Sum(nil))
}

func (c *Client) call(
	ctx context.Context, method, pathStr string, request interface{}, response interface{}) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "failed to marshal request")
	}
	if c.config.Debug {
		c.logger.Printf("request: %s", payload)
	}

	endpoint := fmt.Sprintf("%s%s", c.config.BaseEndpoint, pathStr)
	req, err := http.NewRequest(method, endpoint, strings.NewReader(string(payload)))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.WithContext(ctx)

	nonce := time.Now().UnixNano()
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("ACCESS-KEY", c.config.APIKey)
	req.Header.Add("ACCESS-NONCE", fmt.Sprintf("%d", nonce))
	req.Header.Add("ACCESS-SIGNATURE", c.createSig(nonce, endpoint, string(payload)))

	res, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("status code: %d, body: %s", res.StatusCode, res.Body)
	}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(response); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	if c.config.Debug {
		c.logger.Printf("request: %s", response)
	}
	return nil
}

// HelloRequest request
type HelloRequest struct {
	Name string `json:"name"`
}

// HelloResponse response
type HelloResponse struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

// Hello service
func (c *Client) Hello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	pathStr := "/v1/api/hello"
	method := "GET"
	var res HelloResponse
	if err := c.call(ctx, method, pathStr, req, &res); err != nil {
		return nil, errors.Wrapf(err, "%s %s failed", method, pathStr)
	}
	if res.StatusCode != SuccessStatusCode {
		return nil, errors.Errorf("StatusCode: %d", res.StatusCode)
	}
	return &res, nil
}
