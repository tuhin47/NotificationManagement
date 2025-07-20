package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

	"NotificationManagement/domain"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"NotificationManagement/types"
	"NotificationManagement/utils/errutil"
)

type CurlServiceImpl struct {
	Repo domain.CurlRequestRepository
}

func NewCurlServiceImpl(repo domain.CurlRequestRepository) domain.CurlService {
	return &CurlServiceImpl{Repo: repo}
}

func parseBasicCurl(raw string) (method, url string, headers map[string]string, body string, err error) {
	headers = make(map[string]string)
	method = "GET"
	body = ""

	// Trim whitespace and collapse multiple spaces
	raw = strings.TrimSpace(raw)
	raw = regexp.MustCompile(`\s+`).ReplaceAllString(raw, " ")

	// Match both single and double quotes for the URL
	re := regexp.MustCompile(`curl\s+['"]([^'"]+)['"]`)
	matches := re.FindStringSubmatch(raw)
	if len(matches) < 2 {
		err = errors.New("could not parse URL from curl command")
		return
	}
	url = matches[1]

	// Find -X or --request for method
	if strings.Contains(raw, " -X ") {
		parts := strings.Split(raw, " -X ")
		if len(parts) > 1 {
			methodPart := strings.Fields(parts[1])
			if len(methodPart) > 0 {
				method = strings.ToUpper(methodPart[0])
			}
		}
	}

	// Find all -H 'Header: value'
	reHeader := regexp.MustCompile(`-H\s+'([^:]+):\s?([^']+)'`)
	headersFound := reHeader.FindAllStringSubmatch(raw, -1)
	for _, h := range headersFound {
		if len(h) == 3 {
			headers[h[1]] = h[2]
		}
	}

	// Find --data or -d for body
	if strings.Contains(raw, "--data '") || strings.Contains(raw, "-d '") {
		reBody := regexp.MustCompile(`(--data|-d)\s+'([^']+)'`)
		bodyMatch := reBody.FindStringSubmatch(raw)
		if len(bodyMatch) == 3 {
			body = bodyMatch[2]
			method = "POST"
		}
	}

	return
}

func (s *CurlServiceImpl) ExecuteCurl(req types.CurlRequest) (types.CurlResponse, error) {

	var method, urlStr, body string
	headers := map[string]string{}

	if req.RawCurl != "" {
		m, u, h, b, err := parseBasicCurl(req.RawCurl)
		if err != nil {
			return types.CurlResponse{}, errutil.NewAppErrorWithMessage(
				errutil.ErrInvalidInput,
				err,
				"Failed to parse raw curl command",
			)
		}
		method = m
		urlStr = u
		headers = h
		body = b
	} else {
		method = req.Method
		urlStr = req.URL
		headers = req.Headers
		body = req.Body
	}

	logger.Info("Executing HTTP request", "method", method, "url", urlStr, "headers", headers, "body", body)

	// Store the request in the database
	modelReq, err := req.ToModel()
	modelReq.Method = method
	modelReq.URL = urlStr
	if err == nil {
		_ = s.Repo.Create(context.Background(), modelReq)
	}
	client := &http.Client{}
	request, err := http.NewRequest(method, urlStr, io.NopCloser(strings.NewReader(body)))
	if err != nil {
		return types.CurlResponse{}, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	resp, err := client.Do(request)
	if err != nil {
		return types.CurlResponse{}, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	respHeaders := map[string]string{}
	for k, v := range resp.Header {
		respHeaders[k] = v[0]
	}

	var respBodyVal interface{}
	if json.Valid(respBody) {
		var jsonBody map[string]interface{}
		if err := json.Unmarshal(respBody, &jsonBody); err == nil {
			respBodyVal = jsonBody
		} else {
			respBodyVal = string(respBody)
		}
	} else {
		respBodyVal = string(respBody)
	}

	return types.CurlResponse{
		Status:  resp.StatusCode,
		Headers: respHeaders,
		Body:    respBodyVal,
	}, nil
}

func (s *CurlServiceImpl) GetCurlRequestByID(id uint) (*models.CurlRequest, error) {
	curlRequest, err := s.Repo.GetByID(context.Background(), id)
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrRecordNotFound, err)
	}
	return curlRequest, nil
}

func (s *CurlServiceImpl) UpdateCurlRequest(id uint, req types.CurlRequest) (*models.CurlRequest, error) {
	// First check if the record exists
	existing, err := s.Repo.GetByID(context.Background(), id)
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrRecordNotFound, err)
	}

	// Convert the request to model
	modelReq, err := req.ToModel()
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrInvalidInput, err)
	}

	// Update the existing record with new data
	existing.URL = modelReq.URL
	existing.Method = modelReq.Method
	existing.Headers = modelReq.Headers
	existing.Body = modelReq.Body
	existing.RawCurl = modelReq.RawCurl

	// Save the updated record
	err = s.Repo.Update(context.Background(), existing)
	if err != nil {
		return nil, errutil.NewAppError(errutil.ErrDatabaseQuery, err)
	}

	return existing, nil
}

func (s *CurlServiceImpl) DeleteCurlRequest(id uint) error {
	// First check if the record exists
	_, err := s.Repo.GetByID(context.Background(), id)
	if err != nil {
		return errutil.NewAppError(errutil.ErrRecordNotFound, err)
	}

	err = s.Repo.Delete(context.Background(), id)
	if err != nil {
		return errutil.NewAppError(errutil.ErrDatabaseQuery, err)
	}

	return nil
}
