package types

import (
	"NotificationManagement/models"
	"encoding/json"
)

type CurlRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
	RawCurl string            `json:"rawCurl,omitempty"`
}

type CurlResponse struct {
	Status     int               `json:"status"`
	Headers    map[string]string `json:"headers"`
	Body       interface{}       `json:"body"`
	ErrMessage string            `json:"error,omitempty"`
}

// ToModel converts a types.CurlRequest to a models.CurlRequest
func (cr *CurlRequest) ToModel() (*models.CurlRequest, error) {
	headersJSON, err := json.Marshal(cr.Headers)
	if err != nil {
		return nil, err
	}
	return &models.CurlRequest{
		URL:     cr.URL,
		Method:  cr.Method,
		Headers: string(headersJSON),
		Body:    cr.Body,
		RawCurl: cr.RawCurl,
	}, nil
}
