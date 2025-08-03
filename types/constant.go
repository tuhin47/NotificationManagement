package types

import "google.golang.org/genai"

const ResponseDateFormat = "2006-01-02T15:04:05Z07:00"

type AIResponseStruct interface {
	OllamaResponse | genai.GenerateContentResponse
}
