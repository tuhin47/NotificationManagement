package types

type DeepseekModelResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	ModelName  string `json:"model"`
	BaseURL    string `json:"base_url"`
	ModifiedAt string `json:"modified_at"`
	Size       int64  `json:"size"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
