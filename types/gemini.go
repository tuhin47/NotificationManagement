package types

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

type GeminiRequest struct {
	Contents []*GeminiMessage `json:"contents"`
	// Add other fields like generationConfig, safetySettings if needed
	Model string `json:"model"` // This might be part of the URL, but often included in request body too
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiMessage struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}
