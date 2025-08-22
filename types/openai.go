package types

import (
	"NotificationManagement/models"
	"encoding/json"
)

type OpenAIModelResponse struct {
	ID         uint                   `json:"id"`
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	ModelName  string                 `json:"model"`
	APISecret  models.EncryptedString `json:"api_secret"`
	ModifiedAt string                 `json:"modified_at"`
	CreatedAt  string                 `json:"created_at"`
	UpdatedAt  string                 `json:"updated_at"`
}

func FromOpenAIModel(model *models.OpenAIModel) *OpenAIModelResponse {
	return &OpenAIModelResponse{
		ID:        model.ID,
		Type:      model.Type,
		Name:      model.Name,
		ModelName: model.ModelName,
		APISecret: model.APISecret,
		CreatedAt: model.CreatedAt.Format(ResponseDateFormat),
		UpdatedAt: model.UpdatedAt.Format(ResponseDateFormat),
	}
}

// JSONSchemaProperty represents a property in a JSON schema
type JSONSchemaProperty struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// JSONSchema represents a JSON schema structure
type JSONSchema struct {
	Type                 string                        `json:"type"`
	Properties           map[string]JSONSchemaProperty `json:"properties"`
	Required             []string                      `json:"required"`
	AdditionalProperties bool                          `json:"additionalProperties"`
}

// MarshalJSON implements json.Marshaler interface
func (j JSONSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type                 string                        `json:"type"`
		Properties           map[string]JSONSchemaProperty `json:"properties"`
		Required             []string                      `json:"required"`
		AdditionalProperties bool                          `json:"additionalProperties"`
	}{
		Type:                 j.Type,
		Properties:           j.Properties,
		Required:             j.Required,
		AdditionalProperties: j.AdditionalProperties,
	})
}
