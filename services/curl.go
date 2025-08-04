package services

import (
	"NotificationManagement/domain"
	"NotificationManagement/logger"
	"NotificationManagement/models"
	"NotificationManagement/repositories"
	"NotificationManagement/types"
	"NotificationManagement/utils"
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

func (s *CurlServiceImpl) GetModelByID(id uint) (*models.CurlRequest, error) {
	return s.CurlRepo.GetByID(s.GetInstance().GetContext(), id, &[]string{"AdditionalFields"})
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

	// Trim whitespace and collapse multiple spaces
	raw = strings.TrimSpace(raw)
	raw = regexp.MustCompile(`\s+`).ReplaceAllString(raw, " ")

	// Match single, double, or escaped single quotes for the URL
	re := regexp.MustCompile(`curl\s+(?:'([^']+)'|"([^"]+)"|\\'([^\\']+)\\')`)
	matches := re.FindStringSubmatch(raw)
	if len(matches) < 2 {
		err = errors.New("could not parse URL from curl command")
		return
	}
	// Find the first non-empty group
	for i := 1; i < len(matches); i++ {
		if matches[i] != "" {
			url = matches[i]
			break
		}
	}

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

func (s *CurlServiceImpl) ExecuteCurl(req *models.CurlRequest) (*types.CurlResponse, error) {

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

	return &types.CurlResponse{
		Status:  resp.StatusCode,
		Headers: respHeaders,
		Body:    respBodyVal,
	}, nil
}

func (s *CurlServiceImpl) UpdateModel(id uint, model *models.CurlRequest) (*models.CurlRequest, error) {
	ctx := context.Background()
	existing, err := s.CurlRepo.GetByID(ctx, id, nil)
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

		if len(idsToCheck) > 0 {
			background := context.Background()
			f := []repositories.Filter{
				{Field: "request_id", Op: "=", Value: id},
			}
			background = context.WithValue(background, repositories.ContextStruct{}, &repositories.ContextStruct{Filter: &f})
			existingAdditionalFields, err := s.AdditionalFieldRepo.GetByIDs(background, idsToCheck, nil)
			if err != nil {
				return nil, err
			}

			existingIDsMap := make(map[uint]bool)
			for _, existingAf := range existingAdditionalFields {
				existingIDsMap[existingAf.ID] = true
			}

			for i, af := range *model.AdditionalFields {
				if af.ID != 0 && !existingIDsMap[af.ID] {
					(*model.AdditionalFields)[i].ID = 0
				}

			}
		}
	}
	updatedAssoc, err := utils.SyncHasManyAssociation(s.CurlRepo.GetDB(ctx), &existing, "AdditionalFields", model.AdditionalFields)
	if err != nil {
		return nil, err
	}

	model, err = s.CommonService.UpdateModel(id, model)
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
