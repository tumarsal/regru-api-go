package regru

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the main API client for REG.RU.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	auth       AuthParams
	input      InputParams
	sigAuth    *SignatureAuth
}

// ClientOption configures the Client.
type ClientOption func(*Client)

// NewClient creates a new REG.RU API client.
func NewClient(username, password string, opts ...ClientOption) *Client {
	c := &Client{
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		auth: AuthParams{
			Username: username,
			Password: password,
		},
		input: InputParams{
			OutputContentType: OutputJSON,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = hc
	}
}

// WithBaseURL sets a custom base URL (e.g., for testing).
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithSignatureAuth switches to signature-based authentication.
// The password field is cleared and sig is used instead.
func WithSignatureAuth(sig string) ClientOption {
	return func(c *Client) {
		c.auth.Password = ""
		c.auth.Sig = sig
	}
}

// WithOutputFormat sets the response format (json, yaml, xml, plain).
func WithOutputFormat(format string) ClientOption {
	return func(c *Client) {
		c.input.OutputContentType = format
	}
}

// NewTLSClientAuthOption configures mutual TLS (client certificate) for API requests.
// Required when an SSL certificate is uploaded in REG.RU API settings.
func NewTLSClientAuthOption(certFile, keyFile string) (ClientOption, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("load TLS client cert: %w", err)
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		},
	}
	return func(c *Client) {
		if c.HTTPClient == nil {
			c.HTTPClient = &http.Client{Timeout: 30 * time.Second}
		}
		c.HTTPClient.Transport = transport
	}, nil
}

// request builds and executes an API call.
func (c *Client) request(ctx context.Context, category, function string, params map[string]interface{}) (*Response, error) {
	apiURL := fmt.Sprintf("%s%s/%s/%s", c.BaseURL, DefaultBaseURLPath, category, function)

	formData := url.Values{}

	// Add output format
	if c.input.OutputContentType != "" {
		formData.Set("output_content_type", c.input.OutputContentType)
	}

	// Check if we need JSON input format for complex structures
	hasComplex := false
	for _, v := range params {
		switch v.(type) {
		case []interface{}, map[string]interface{}, []string, []ServiceID:
			hasComplex = true
		}
	}

	if hasComplex {
		// Use JSON input_data
		jsonData, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("marshal input_data: %w", err)
		}
		formData.Set("input_format", InputJSON)
		formData.Set("input_data", string(jsonData))
	} else {
		// Use plain form parameters
		for k, v := range params {
			switch val := v.(type) {
			case string:
				formData.Set(k, val)
			case int:
				formData.Set(k, fmt.Sprintf("%d", val))
			case float64:
				formData.Set(k, fmt.Sprintf("%g", val))
			case bool:
				if val {
					formData.Set(k, "1")
				} else {
					formData.Set(k, "0")
				}
			default:
				formData.Set(k, fmt.Sprintf("%v", v))
			}
		}
	}

	// Auth: password or per-request RSA signature
	if c.auth.Username != "" {
		formData.Set("username", c.auth.Username)
	}
	if c.sigAuth != nil {
		sig, err := c.sigAuth.Sign(buildSignParams(c.auth.Username, c.input.OutputContentType, hasComplex, params))
		if err != nil {
			return nil, fmt.Errorf("sign request: %w", err)
		}
		formData.Set("sig", sig)
	} else if c.auth.Sig != "" {
		formData.Set("sig", c.auth.Sig)
	} else if c.auth.Password != "" {
		formData.Set("password", c.auth.Password)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d: %s", resp.StatusCode, string(body))
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w (body: %s)", err, string(body))
	}

	return &result, nil
}

// call is a helper for simple API calls without custom parameters.
func (c *Client) call(ctx context.Context, category, function string, params map[string]interface{}, answer interface{}) error {
	resp, err := c.request(ctx, category, function, params)
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return fmt.Errorf("api error: %s", resp.Error())
	}
	if answer != nil && len(resp.Answer) > 0 {
		if err := json.Unmarshal(resp.Answer, answer); err != nil {
			return fmt.Errorf("unmarshal answer: %w", err)
		}
	}
	return nil
}

// nop calls are test/heartbeat functions.
func (c *Client) nop(ctx context.Context, path string) error {
	return c.call(ctx, path, "nop", nil, nil)
}

// Nop checks API availability with test credentials (no auth needed).
func (c *Client) Nop(ctx context.Context) error {
	_, err := c.request(ctx, "nop", "nop", nil)
	return err
}

// ResellerNop checks API availability for resellers.
func (c *Client) ResellerNop(ctx context.Context) error {
	return c.nop(ctx, "reseller_nop")
}

// GetUserID returns the authenticated user's ID.
func (c *Client) GetUserID(ctx context.Context) (string, error) {
	var answer struct {
		UserID string `json:"user_id"`
	}
	if err := c.call(ctx, "", "get_user_id", nil, &answer); err != nil {
		return "", err
	}
	return answer.UserID, nil
}

// GetServiceID returns the service ID for a domain or service.
func (c *Client) GetServiceID(ctx context.Context, service ServiceID) (string, error) {
	params := serviceToMap(service)
	var answer struct {
		ServiceID string `json:"service_id"`
	}
	if err := c.call(ctx, "", "get_service_id", params, &answer); err != nil {
		return "", err
	}
	return answer.ServiceID, nil
}

// serviceToMap converts ServiceID to map.
func serviceToMap(s ServiceID) map[string]interface{} {
	m := make(map[string]interface{})
	if s.DomainName != "" {
		m["domain_name"] = s.DomainName
	}
	if s.ServiceID != "" {
		m["service_id"] = s.ServiceID
	}
	if s.ServiceType != "" {
		m["servtype"] = s.ServiceType
	}
	if s.SubType != "" {
		m["subtype"] = s.SubType
	}
	if s.UpLinkServiceID != "" {
		m["uplink_service_id"] = s.UpLinkServiceID
	}
	return m
}

// toJSON helper.
func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// fromXML helper.
func fromXML(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func buildSignParams(username, outputFormat string, hasComplex bool, params map[string]interface{}) map[string]interface{} {
	m := map[string]interface{}{
		"username": username,
	}
	if outputFormat != "" {
		m["output_content_type"] = outputFormat
	}
	if hasComplex {
		m["input_format"] = InputJSON
	}
	for k, v := range params {
		m[k] = v
	}
	return m
}
