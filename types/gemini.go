package types

import (
	"NotificationManagement/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []GeminiPart `json:"parts"`
			Role  string       `json:"role"`
		} `json:"content"`
		FinishReason  string `json:"finishReason"`
		Index         int    `json:"index"`
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"candidates"`
	PromptFeedback struct {
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"promptFeedback"`
}

type GeminiModelResponse struct {
	ID         uint   `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	ModelName  string `json:"model"`
	APISecret  string `json:"api_secret"`
	ModifiedAt string `json:"modified_at"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type GeminiRequest struct {
	Contents []*GeminiMessage `json:"contents"`
	// Add other fields like generationConfig, safetySettings if needed
	Model string `json:"model"` // This might be part of the URL, but often included in request body too
}

func (r *GeminiRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Model, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Contents, validation.Required, validation.Each(validation.By(func(value interface{}) error {
			if v, ok := value.(*GeminiMessage); ok {
				return v.Validate()
			}
			return nil
		}))),
	)
}

type GeminiPart struct {
	Text string `json:"text"`
}

func (p *GeminiPart) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Text, validation.Required),
	)
}

type GeminiMessage struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

func (m *GeminiMessage) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Role, validation.Required, validation.In("system", "user", "assistant")),
		validation.Field(&m.Parts, validation.Required, validation.Each(validation.By(func(value interface{}) error {
			if v, ok := value.(GeminiPart); ok {
				return v.Validate()
			}
			return nil
		}))),
	)
}

func FromGeminiModel(model *models.GeminiModel) *GeminiModelResponse {
	return &GeminiModelResponse{
		ID:        model.ID,
		Type:      model.Type,
		Name:      model.Name,
		ModelName: model.ModelName,
		APISecret: model.APISecret,
		CreatedAt: model.CreatedAt.Format(ResponseDateFormat),
		UpdatedAt: model.UpdatedAt.Format(ResponseDateFormat),
	}
}
