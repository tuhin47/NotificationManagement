package types

type CurlRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
	RawCurl string            `json:"rawCurl,omitempty"`
}

type CurlResponse struct {
	Status     int               `json:"status"`
	Headers    map[string]string `json:"headers"`
	Body       interface{}       `json:"body"`
	ErrMessage string            `json:"error,omitempty"`
}
