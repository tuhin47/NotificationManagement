package types

type OllamaRequest struct {
	Model    string           `json:"model" validate:"required"`
	Messages []*OllamaMessage `json:"messages" validate:"required,min=1"`
	Stream   bool             `json:"stream"`
	Format   *OllamaFormat    `json:"format,omitempty"`
	Options  *OllamaOptions   `json:"options,omitempty"`
	Think    bool             `json:"think,omitempty"`
}
type OllamaPullRequest struct {
	Name string `json:"name" validate:"required"`
}

type OllamaMessage struct {
	Role    string `json:"role" validate:"required,oneof=system user assistant"`
	Content string `json:"content" validate:"required"`
}

type OllamaFormat struct {
	Type       string                          `json:"type" validate:"required"`
	Properties map[string]OllamaFormatProperty `json:"properties" validate:"required"`
	Required   []string                        `json:"required" validate:"required"`
}

type OllamaOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
}

type OllamaFormatProperty struct {
	Type        string `json:"type" validate:"required"`
	Description string `json:"description,omitempty"`
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
