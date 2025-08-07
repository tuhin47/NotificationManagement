package types

import (
	"NotificationManagement/types/ollama"
	"google.golang.org/genai"
)

const ResponseDateFormat = "2006-01-02T15:04:05Z07:00"

const (
	ResponseTypeJSON = "json"
	ResponseTypeXML  = "xml"
	ResponseTypeHTML = "html"
	ResponseTypeText = "text"
)

type AIResponseStruct interface {
	ollama.OllamaResponse | genai.GenerateContentResponse
}
