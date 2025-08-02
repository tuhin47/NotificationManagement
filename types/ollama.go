package types

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type OllamaRequest struct {
	Model    string           `json:"model"`
	Messages []*OllamaMessage `json:"messages"`
	Stream   bool             `json:"stream"`
	Format   *OllamaFormat    `json:"format,omitempty"`
	Options  *OllamaOptions   `json:"options,omitempty"`
	Think    bool             `json:"think,omitempty"`
}

func (r *OllamaRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Model, validation.Required),
		validation.Field(&r.Messages, validation.Required, validation.Length(1, 0), validation.Each(validation.By(func(value interface{}) error {
			if v, ok := value.(*OllamaMessage); ok {
				return v.Validate()
			}
			return nil
		}))),
	)
}

type OllamaPullRequest struct {
	Name string `json:"name"`
}

func (r *OllamaPullRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required),
	)
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (m *OllamaMessage) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Role, validation.Required, validation.In("system", "user", "assistant")),
		validation.Field(&m.Content, validation.Required),
	)
}

type OllamaFormat struct {
	Type       string                          `json:"type"`
	Properties map[string]OllamaFormatProperty `json:"properties"`
	Required   []string                        `json:"required"`
}

func (f *OllamaFormat) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Type, validation.Required),
		validation.Field(&f.Properties, validation.Required),
		validation.Field(&f.Required, validation.Required),
	)
}

type OllamaOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
}

type OllamaFormatProperty struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

func (p *OllamaFormatProperty) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Type, validation.Required),
	)
}

type OllamaResponse struct {
	Model              string         `json:"model"`
	CreatedAt          string         `json:"created_at"`
	Message            *OllamaMessage `json:"message"`
	DoneReason         string         `json:"done_reason"`
	Done               bool           `json:"done"`
	TotalDuration      int64          `json:"total_duration"`
	LoadDuration       int64          `json:"load_duration"`
	PromptEvalCount    int            `json:"prompt_eval_count"`
	PromptEvalDuration int64          `json:"prompt_eval_duration"`
	EvalCount          int            `json:"eval_count"`
	EvalDuration       int64          `json:"eval_duration"`
}
