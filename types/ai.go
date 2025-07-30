package types

type MakeAIRequestPayload struct {
	CurlRequestID uint `json:"curl_request_id" validate:"required"`
	ModelID       uint `json:"model_id" validate:"required"`
}
