package types

import (
	"NotificationManagement/models"
	"NotificationManagement/utils/errutil"
	"encoding/json"
	"fmt"
	"os"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CurlRequest struct {
	URL              string                   `json:"url"`
	Method           string                   `json:"method"`
	Headers          map[string]string        `json:"headers,omitempty"`
	Body             string                   `json:"body,omitempty"`
	RawCurl          string                   `json:"rawCurl,omitempty"`
	ResponseType     string                   `json:"responseType,omitempty"`
	UserID           uint                     `json:"user_id"`
	AdditionalFields []AdditionalFieldRequest `json:"additional_fields"`
}

func (cr *CurlRequest) Validate() error {
	return validation.ValidateStruct(cr,
		validation.Field(&cr.ResponseType, validation.In(ResponseTypeJSON, ResponseTypeXML, ResponseTypeHTML, ResponseTypeText)),
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
	UserID     uint              `json:"user_id,omitempty"`
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
		ResponseType:     cr.ResponseType,
		UserID:           cr.UserID,
		AdditionalFields: &props,
	}, nil
}

func (response *CurlResponse) GetAssistantContent(respType string) (*string, error) {
	if response.ErrMessage != "" {
		return nil, errutil.NewAppError(errutil.ErrExternalServiceError, fmt.Errorf("error : %s", response.ErrMessage))
	}
	if response.Body == nil {
		return nil, errutil.NewAppError(errutil.ErrEmptyResponse, fmt.Errorf("Body not valid "))
	}

	switch respType {
	case ResponseTypeJSON:
		bodyBytes, err := json.Marshal(response.Body)
		if err != nil {
			return nil, errutil.NewAppError(errutil.ErrCurlMarshalResponseBodyFailed, err)
		}
		s := "Here is a json string  `" + string(bodyBytes) + "`"
		return &s, nil
	case ResponseTypeHTML:
		htmlContent, ok := response.Body.(string)
		if !ok {
			return nil, errutil.NewAppError(errutil.ErrCurlInvalidResponseBodyType, fmt.Errorf("response body is not a string for HTML type"))
		}
		tmpfile, err := os.CreateTemp("", "response-*.html")
		if err != nil {
			return nil, errutil.NewAppError(errutil.ErrCurlCreateTempFileFailed, err)
		}
		defer tmpfile.Close()

		_, err = tmpfile.WriteString(htmlContent)
		if err != nil {
			return nil, errutil.NewAppError(errutil.ErrCurlWriteTempFileFailed, err)
		}
		filePath := tmpfile.Name()
		return &filePath, nil
	case ResponseTypeXML:
		xmlContent, ok := response.Body.(string)
		if !ok {
			return nil, errutil.NewAppError(errutil.ErrCurlInvalidResponseBodyType, fmt.Errorf("response body is not a string for XML type"))
		}
		s := "Here is an XML string  `" + xmlContent + "`"
		return &s, nil
	case ResponseTypeText:
		textContent, ok := response.Body.(string)
		if !ok {
			return nil, errutil.NewAppError(errutil.ErrCurlInvalidResponseBodyType, fmt.Errorf("response body is not a string for Text type"))
		}
		s := "Here is a text string  `" + textContent + "`"
		return &s, nil
	default:
		return nil, errutil.NewAppError(errutil.ErrCurlUnsupportedResponseType, fmt.Errorf("unsupported response type: %s", respType))
	}
}
