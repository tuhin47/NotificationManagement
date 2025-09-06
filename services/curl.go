package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"NotificationManagement/repositories"
	"NotificationManagement/services/helper"
	"NotificationManagement/types"
	"NotificationManagement/utils/errutil"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type CurlServiceImpl struct {
	domain.CommonService[models.CurlRequest]
	CurlRepo            domain.CurlRequestRepository
	AdditionalFieldRepo domain.AdditionalFieldsRepository
}

func (s *CurlServiceImpl) GetModelById(c context.Context, id uint, preloads *[]string) (*models.CurlRequest, error) {
	if preloads == nil {
		preloads = &[]string{"AdditionalFields"}
	}
	return s.CurlRepo.GetByID(s.GetInstance().ProcessContext(c), id, preloads)
}

func NewCurlService(repo domain.CurlRequestRepository, fieldsRepository domain.AdditionalFieldsRepository) domain.CurlService {
	service := &CurlServiceImpl{
		CurlRepo:            repo,
		AdditionalFieldRepo: fieldsRepository,
	}
	service.CommonService = NewCommonService(repo, service)
	return service
}

func parseBasicCurl(raw string) (method, url string, headers map[string]string, body string, err error) {
	headers = make(map[string]string)
	method = "GET"
	body = ""

	// Trim whitespace
	raw = strings.TrimSpace(raw)

	// Regex to find the URL. It can be unquoted, single-quoted, or double-quoted.
	// Also handles the optional -s flag.
	reURL := regexp.MustCompile(`curl(?:\s+-s)?\s+(?:'([^']+)'|"([^"]+)"|\\'([^\\']+)\\'|([^\s'"]+))`)
	matchesURL := reURL.FindStringSubmatch(raw)

	if len(matchesURL) < 2 {
		err = errutil.NewAppError(errutil.ErrCurlParseError, errors.New("could not parse URL from curl command"))
		return
	}

	// Find the first non-empty group for the URL
	for i := 1; i < len(matchesURL); i++ {
		if matchesURL[i] != "" {
			url = matchesURL[i]
			break
		}
	}

	// Find -X or --request for method
	reMethod := regexp.MustCompile(`(?:-X|--request)\s+([A-Za-z]+)`)
	methodMatch := reMethod.FindStringSubmatch(raw)
	if len(methodMatch) > 1 {
		method = strings.ToUpper(methodMatch[1])
	}

	// Find all -H 'Header: value' or -H "Header: value"
	reHeader := regexp.MustCompile(`-H\s+(?:'([^:]+):\s?([^']+)'|"([^:]+):\s?([^"]+)")`)
	headersFound := reHeader.FindAllStringSubmatch(raw, -1)
	for _, h := range headersFound {
		if h[1] != "" && h[2] != "" { // Single quoted header
			headers[h[1]] = h[2]
		} else if h[3] != "" && h[4] != "" { // Double quoted header
			headers[h[3]] = h[4]
		}
	}

	// Find --data 'body' or -d 'body' or --data "body" or -d "body"
	reBody := regexp.MustCompile(`(--data|-d)\s+(?:'([^']+)'|"([^"]+)")`)
	bodyMatch := reBody.FindStringSubmatch(raw)
	if len(bodyMatch) > 1 {
		if bodyMatch[2] != "" { // Single quoted body
			body = bodyMatch[2]
		} else if bodyMatch[3] != "" { // Double quoted body
			body = bodyMatch[3]
		}
		if method == "GET" { // Default to POST if body is present and method is not explicitly set
			method = "POST"
		}
	}

	return
}

func (s *CurlServiceImpl) ProcessCurlRequest(c context.Context, req *models.CurlRequest) (*types.CurlResponse, error) {

	var method, urlStr, body string
	headers := map[string]string{}

	if req.RawCurl != "" {
		m, u, h, b, err := parseBasicCurl(req.RawCurl)
		if err != nil {
			return &types.CurlResponse{}, errutil.NewAppErrorWithMessage(
				errutil.ErrInvalidRequestBody,
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
		body = req.Body
		// TODO
		//headers = req.Headers
	}

	logger.Info("Executing HTTP request", "method", method, "url", urlStr, "headers", headers, "body", body)

	transport := &http.Transport{}
	if strings.Contains(req.RawCurl, "--insecure") {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	client := &http.Client{
		Transport: transport,
	}
	request, err := http.NewRequest(method, urlStr, io.NopCloser(strings.NewReader(body)))
	if err != nil {
		return &types.CurlResponse{}, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	resp, err := client.Do(request)
	if err != nil {
		return &types.CurlResponse{}, errutil.NewAppError(errutil.ErrExternalServiceError, err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	respHeaders := map[string]string{}
	for k, v := range resp.Header {
		respHeaders[k] = v[0]
	}

	var respBodyVal interface{}
	if req.ResponseType == types.ResponseTypeHTML {
		respBodyVal = string(respBody)
	} else if json.Valid(respBody) {
		var jsonBody map[string]interface{}
		if err := json.Unmarshal(respBody, &jsonBody); err == nil {
			respBodyVal = jsonBody
		} else {
			respBodyVal = string(respBody)
		}
	} else {
		respBodyVal = string(respBody)
	}

	return &types.CurlResponse{
		Status:  resp.StatusCode,
		Headers: respHeaders,
		Body:    respBodyVal,
		UserID:  req.UserID,
	}, nil
}

func (s *CurlServiceImpl) UpdateModel(c context.Context, id uint, model *models.CurlRequest) (*models.CurlRequest, error) {
	existing, err := s.CurlRepo.GetByID(c, id, nil)
	if err != nil {
		return nil, err
	}

	if model.AdditionalFields != nil && len(*model.AdditionalFields) > 0 {
		var idsToCheck []uint
		for i, af := range *model.AdditionalFields {
			if af.ID != 0 {
				idsToCheck = append(idsToCheck, af.ID)
			}
			(*model.AdditionalFields)[i].RequestID = id
		}
		existingIDsMap := make(map[uint]bool)

		if len(idsToCheck) > 0 {
			if txContext, ok := repositories.GetTxContext(c); ok {
				filters := append(txContext.Filter, repositories.NewFilter("request_id", "=", id))
				txContext.Filter = filters
			}
			existingAdditionalFields, err := s.AdditionalFieldRepo.GetByIDs(c, idsToCheck, nil)
			if err != nil {
				return nil, err
			}

			for _, existingAf := range existingAdditionalFields {
				existingIDsMap[existingAf.ID] = true
			}
		}
		for i, af := range *model.AdditionalFields {
			if af.ID != 0 && !existingIDsMap[af.ID] {
				(*model.AdditionalFields)[i].ID = 0
			}
		}
	}

	updatedAssoc, err := helper.SyncHasManyAssociation(s.CurlRepo.GetDB(c), &existing, "AdditionalFields", model.AdditionalFields)
	if err != nil {
		return nil, err
	}

	model, err = s.CommonService.UpdateModel(c, id, model)
	if err != nil {
		return nil, err
	}
	if updatedAssoc != nil {
		if props, ok := updatedAssoc.(*[]models.AdditionalFields); ok {
			existing.AdditionalFields = props
		}
	}

	return model, nil
}
