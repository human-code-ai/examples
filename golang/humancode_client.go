package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
	"time"
)

type ClientResponse[T any] struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result T      `json:"result"`
}

type GetSessionIdResult struct {
	SessionId string `json:"session_id"`
}

type VerifyResult struct {
	HumanId string `json:"human_id"`
}

// ClientConfig holds the configuration for the human code client.
type ClientConfig struct {
	// BaseUrl is the base URL for the API.
	BaseUrl string
	// Debug indicates whether debug mode is enabled.
	Debug bool
	// AppId is the application ID.
	AppId string
	// AppKey is the application key used for signing requests.
	AppKey string
}

// HumanCodeClient is a client for interacting with the human code API.
type HumanCodeClient struct {
	// apiClient is the HTTP client used to make requests.
	apiClient *resty.Client
	// config holds the configuration for the client.
	config *ClientConfig
}

// NewHumanCodeClient creates a new instance of HumanCodeClient.
// It initializes a resty client with the provided HTTP client and configuration.
//
// Parameters:
//   - httpClient: The underlying HTTP client to use.
//   - config: The configuration for the human code client.
//
// Returns:
//   - A pointer to a new HumanCodeClient instance.
func NewHumanCodeClient(httpClient http.Client, config *ClientConfig) *HumanCodeClient {
	// Create a new resty client with the provided HTTP client.
	c := resty.NewWithClient(&httpClient)
	// Set the base URL for the client.
	c.BaseURL = config.BaseUrl
	// Enable debug mode if configured.
	c.Debug = config.Debug
	// Set the default content type for requests.
	c.SetHeader("Content-Type", "application/json")

	return &HumanCodeClient{
		apiClient: c,
		config:    config,
	}
}

// genSign generates a HMAC-SHA256 signature for the given data using the application key.
//
// Parameters:
//   - data: The data to be signed.
//
// Returns:
//   - A hex-encoded string representing the signature.
//   - An error if there is an issue writing the data to the hash.
func (h *HumanCodeClient) genSign(data string) (string, error) {
	// Create a new HMAC-SHA256 hash using the application key.
	sha256hash := hmac.New(sha256.New, []byte(h.config.AppKey))
	// Write the data to the hash.
	if _, err := io.WriteString(sha256hash, data); err != nil {
		return "", err
	}
	// Return the hex-encoded hash sum.
	return hex.EncodeToString(sha256hash.Sum(nil)), nil
}

// GetConfig returns the configuration of the human code client.
//
// Returns:
//   - A pointer to the ClientConfig struct.
func (h *HumanCodeClient) GetConfig() *ClientConfig {
	return h.config
}

// GetSessionId retrieves a session ID from the API.
// It sends a POST request with a timestamp and nonce string, signed with the application key.
//
// Parameters:
//   - nonceStr: A nonce string for the request.
//
// Returns:
//   - A pointer to a GetSessionIdResult struct containing the session ID.
//   - An error if there is an issue with JSON marshaling, signing, the HTTP request, or the API response.
func (h *HumanCodeClient) GetSessionId(nonceStr string) (*GetSessionIdResult, error) {
	// Get the current timestamp in milliseconds.
	timeStamp := time.Now().UnixMilli()

	// Prepare the POST request body as a JSON object.
	postBody, err := json.Marshal(map[string]string{
		"timestamp": fmt.Sprintf("%d", timeStamp),
		"nonce_str": nonceStr,
	})
	// Handle JSON marshaling errors.
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Generate the signature for the request body.
	sign, err := h.genSign(string(postBody))
	fmt.Println(sign)
	// Handle signature generation errors.
	if err != nil {
		return nil, err
	}

	// Create a new instance to hold the API response.
	result := &ClientResponse[GetSessionIdResult]{}

	// Send the POST request to the API.
	resp, err := h.apiClient.R().
		SetBody(postBody).
		SetResult(result).
		Post(fmt.Sprintf("/api/session/v2/get_id?app_id=%s&sign=%s", h.config.AppId, sign))
	// Handle HTTP request errors.
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	// Check if the API response indicates an error.
	if resp.IsError() || result.Code != 0 {
		return nil, fmt.Errorf("API error: code: %d, msg: %s", result.Code, result.Msg)
	}

	// Return the session ID result.
	return &result.Result, nil
}

// GenRegistrationUrl generates a registration URL for the given session ID and callback URL.
//
// Parameters:
//   - sessionId: The session ID for the registration.
//   - callBackUrl: The URL to redirect to after registration.
//
// Returns:
//   - A string representing the registration URL.
//   - An error if there is an issue formatting the URL.
func (h *HumanCodeClient) GenRegistrationUrl(sessionId string, callBackUrl string) (string, error) {
	// Get the current timestamp in milliseconds.
	timeStamp := time.Now().UnixMilli()
	// Generate the registration URL.
	return fmt.Sprintf("%s/authentication/index.html?session_id=%s&callback_url=%s&ts=%d#/", h.config.BaseUrl, sessionId, callBackUrl, timeStamp), nil
}

// GenVerificationUrl generates a verification URL for the given session ID, human ID, and callback URL.
//
// Parameters:
//   - sessionId: The session ID for the verification.
//   - humanId: The human ID for the verification.
//   - callBackUrl: The URL to redirect to after verification.
//
// Returns:
//   - A string representing the verification URL.
//   - An error if there is an issue formatting the URL.
func (h *HumanCodeClient) GenVerificationUrl(sessionId string, humanId string, callBackUrl string) (string, error) {
	// Get the current timestamp in milliseconds.
	timeStamp := time.Now().UnixMilli()
	// Generate the verification URL.
	return fmt.Sprintf("%s/authentication/index.html?session_id=%s&human_id=%s&callback_url=%s&ts=%d#/", h.config.BaseUrl, sessionId, humanId, callBackUrl, timeStamp), nil
}

// Verify verifies a verification code for a given session ID.
// It sends a POST request with the session ID, verification code, timestamp, and nonce string, signed with the application key.
//
// Parameters:
//   - sessionId: The session ID for the verification.
//   - vCode: The verification code to be verified.
//   - nonceStr: A nonce string for the request.
//
// Returns:
//   - A pointer to a VerifyResult struct containing the verification result.
//   - An error if there is an issue with JSON marshaling, signing, the HTTP request, or the API response.
func (h *HumanCodeClient) Verify(sessionId string, vCode string, nonceStr string) (*VerifyResult, error) {
	// Get the current timestamp in milliseconds.
	timeStamp := time.Now().UnixMilli()

	// Prepare the POST request body as a JSON object.
	postBody, err := json.Marshal(map[string]string{
		"session_id": sessionId,
		"vcode":      vCode,
		"timestamp":  fmt.Sprintf("%d", timeStamp),
		"nonce_str":  nonceStr,
	})
	// Handle JSON marshaling errors.
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Generate the signature for the request body.
	sign, err := h.genSign(string(postBody))
	// Handle signature generation errors.
	if err != nil {
		return nil, err
	}

	// Create a new instance to hold the API response.
	result := &ClientResponse[VerifyResult]{}

	// Send the POST request to the API.
	resp, err := h.apiClient.R().
		SetBody(postBody).
		SetResult(result).
		Post(fmt.Sprintf("/api/vcode/v2/verify?app_id=%s&sign=%s", h.config.AppId, sign))
	// Handle HTTP request errors.
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	// Check if the API response indicates an error.
	if resp.IsError() || result.Code != 0 {
		return nil, fmt.Errorf("API error: code: %d, msg: %s", result.Code, result.Msg)
	}

	// Return the verification result.
	return &result.Result, nil
}
