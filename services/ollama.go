package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/models"
	"NotificationManagement/types"
	"NotificationManagement/utils/errutil"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OllamaServiceImpl struct {
	Repo        domain.DeepseekModelRepository
	CurlService domain.CurlService
}

func NewOllamaService(repo domain.DeepseekModelRepository, curlService domain.CurlService) domain.OllamaService {
	return &OllamaServiceImpl{Repo: repo, CurlService: curlService}
}

func (s *OllamaServiceImpl) MakeAIRequest(mod *models.DeepseekModel, requestId uint) (*types.OllamaResponse, error) {

	// Get the CurlRequest by ID
	curlRequest, err := s.CurlService.GetCurlRequestByID(requestId)
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrRecordNotFound, err)
	}

	curlResponse, err := s.CurlService.ExecuteCurl(curlRequest)

	if err != nil {
		return nil, err
	}
	// Make HTTP request to Ollama
	respBody, err := deepseekCall(mod, curlResponse)
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}

	// Parse the response
	var ollamaResp types.OllamaResponse
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return nil, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}

	return &ollamaResp, nil
}

func (s *OllamaServiceImpl) PullModel(model models.DeepseekModel) error {
	// Implementation for PullModel will be added later
	return fmt.Errorf("PullModel not implemented yet")
}

func deepseekCall(model *models.DeepseekModel, response *types.CurlResponse) ([]byte, error) {
	// Build assistant content from CurlResponse.Body
	var assistantContent string
	if response.ErrMessage == "" && response.Body != nil {
		bodyBytes, err := json.Marshal(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response body: %w", err)
		}
		assistantContent = "Here is a json string  `" + string(bodyBytes) + "`"
	} else {
		assistantContent = "No content available"
	}

	properties := &map[string]types.OllamaFormatProperty{
		"IsCorrect": {
			Type:        "boolean",
			Description: "This holds the true or false value for the Statement",
		},
		"Rate": {
			Type:        "number",
			Description: "Rate from Json",
		},
		"TargetRate": {
			Type:        "number",
			Description: "Targeted Rate",
		},
	}
	ollamaReq := types.OllamaRequest{
		Model: model.ModelName,
		Messages: []*types.OllamaMessage{
			{
				Role:    "assistant",
				Content: assistantContent,
			},
			{
				Role:    "user",
				Content: "Please check the current rate from the json.Is it greater than 125 ? Return Json Response ",
			},
		},
		Stream: false,
		Format: &types.OllamaFormat{
			Type:       "object",
			Properties: *properties,
			Required:   []string{"IsCorrect", "Rate", "TargetRate"},
		},
		Options: &types.OllamaOptions{
			Temperature: 0.5,
		},
		Think: true,
	}

	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/chat", model.BaseURL)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}
