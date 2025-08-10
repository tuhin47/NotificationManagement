package ollama

import "github.com/go-ozzo/ozzo-validation/v4"

type Request struct {
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
	Stream   bool       `json:"stream"`
	Format   *Format    `json:"format,omitempty"`
	Options  *Options   `json:"options,omitempty"`
	Think    bool       `json:"think,omitempty"`
}

func (r *Request) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Model, validation.Required),
		validation.Field(&r.Messages, validation.Required, validation.Length(1, 0), validation.Each(validation.By(func(value interface{}) error {
			if v, ok := value.(*Message); ok {
				return v.Validate()
			}
			return nil
		}))),
	)
}

type PullRequest struct {
	Name string `json:"name"`
}

func (r *PullRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required),
	)
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (m *Message) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Role, validation.Required, validation.In("system", "user", "assistant")),
		validation.Field(&m.Content, validation.Required),
	)
}

type Format struct {
	Type       string                    `json:"type"`
	Properties map[string]FormatProperty `json:"properties"`
	Required   []string                  `json:"required"`
}

func (f *Format) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Type, validation.Required),
		validation.Field(&f.Properties, validation.Required),
		validation.Field(&f.Required, validation.Required),
	)
}

type Options struct {
	Temperature float64 `json:"temperature,omitempty"`
}

type FormatProperty struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

func (p *FormatProperty) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Type, validation.Required),
	)
}

type Response struct {
	Model              string   `json:"model"`
	CreatedAt          string   `json:"created_at"`
	Message            *Message `json:"message"`
	DoneReason         string   `json:"done_reason"`
	Done               bool     `json:"done"`
	TotalDuration      int64    `json:"total_duration"`
	LoadDuration       int64    `json:"load_duration"`
	PromptEvalCount    int      `json:"prompt_eval_count"`
	PromptEvalDuration int64    `json:"prompt_eval_duration"`
	EvalCount          int      `json:"eval_count"`
	EvalDuration       int64    `json:"eval_duration"`
}
