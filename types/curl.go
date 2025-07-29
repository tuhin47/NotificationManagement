package types

import (
	"NotificationManagement/models"
	"encoding/json"
)

type CurlRequest struct {
	URL                    string                        `json:"url"`
	Method                 string                        `json:"method"`
	Headers                map[string]string             `json:"headers,omitempty"`
	Body                   string                        `json:"body,omitempty"`
	RawCurl                string                        `json:"rawCurl,omitempty"`
	OllamaFormatProperties []OllamaFormatPropertyRequest `json:"ollama_format_properties"`
}

type OllamaFormatPropertyRequest struct {
	PropertyName string `json:"property_name"`
	Type         string `json:"type"`
	Description  string `json:"description,omitempty"`
	RequestID    uint   `json:"request_id,omitempty"`
	ID           uint   `json:"id,omitempty"`
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
	var props []models.OllamaFormatProperty
	for _, p := range cr.OllamaFormatProperties {
		props = append(props, models.OllamaFormatProperty{
			PropertyName: p.PropertyName,
			Type:         p.Type,
			Description:  p.Description,
			RequestID:    p.RequestID,
			ID:           p.ID,
		})
	}
	return &models.CurlRequest{
		URL:                    cr.URL,
		Method:                 cr.Method,
		Headers:                string(headersJSON),
		Body:                   cr.Body,
		RawCurl:                cr.RawCurl,
		OllamaFormatProperties: &props,
	}, nil
}
