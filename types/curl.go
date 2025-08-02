package types

import (
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"
	"encoding/json"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CurlRequest struct {
	URL              string                   `json:"url"`
	Method           string                   `json:"method"`
	Headers          map[string]string        `json:"headers,omitempty"`
	Body             string                   `json:"body,omitempty"`
	RawCurl          string                   `json:"rawCurl,omitempty"`
	AdditionalFields []AdditionalFieldRequest `json:"additional_fields"`
}

func (cr *CurlRequest) Validate() error {
	return validation.ValidateStruct(cr,
		validation.Field(&cr.URL, validation.Required, is.URL, validation.Length(0, 2048)),
		validation.Field(&cr.Method, validation.Required, validation.Length(1, 10)),
		validation.Field(&cr.AdditionalFields, validation.Each(validation.By(func(value interface{}) error {
			if v, ok := value.(AdditionalFieldRequest); ok {
				return v.Validate()
			}
			if v, ok := value.(*AdditionalFieldRequest); ok {
				return v.Validate()
			}
			return nil // Or return an error if the type is unexpected
		}))),
	)
}

type AdditionalFieldRequest struct {
	PropertyName string `json:"property_name"`
	Type         string `json:"type"`
	Description  string `json:"description,omitempty"`
	RequestID    uint   `json:"request_id,omitempty"`
	ID           uint   `json:"id,omitempty"`
}

func (ar *AdditionalFieldRequest) Validate() error {
	return validation.ValidateStruct(ar,
		validation.Field(&ar.PropertyName, validation.Required, validation.Length(1, 100)),
		validation.Field(&ar.Type, validation.Required, validation.In("number", "boolean", "text"), validation.Length(1, 10)),
	)
}

type CurlResponse struct {
	Status     int               `json:"status"`
	Headers    map[string]string `json:"headers"`
	Body       interface{}       `json:"body"`
	ErrMessage string            `json:"error,omitempty"`
}

// ToModel converts a types.CurlRequest to a models.CurlRequest
func (cr *CurlRequest) ToModel() (*models.CurlRequest, error) {
	err := cr.Validate()
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrInvalidRequestBody, err)
	}
	headersJSON, err := json.Marshal(cr.Headers)
	if err != nil {
		return nil, err
	}
	var props []models.AdditionalFields
	for _, p := range cr.AdditionalFields {
		props = append(props, models.AdditionalFields{
			PropertyName: p.PropertyName,
			Type:         p.Type,
			Description:  p.Description,
			RequestID:    p.RequestID,
			ID:           p.ID,
		})
	}
	return &models.CurlRequest{
		URL:              cr.URL,
		Method:           cr.Method,
		Headers:          string(headersJSON),
		Body:             cr.Body,
		RawCurl:          cr.RawCurl,
		AdditionalFields: &props,
	}, nil
}

func (response *CurlResponse) GetAssistantContent() (string, error) {
	if response.ErrMessage == "" && response.Body != nil {
		bodyBytes, err := json.Marshal(response.Body)
		if err != nil {
			return "", fmt.Errorf("failed to marshal response body: %w", err)
		}
		return "Here is a json string  `" + string(bodyBytes) + "`", nil
	}
	return "No content available", nil
}
