package services

import (
	"io"
	"net/http"
	"strings"

	"NotificationManagement/domain"
	"NotificationManagement/types"
)

type CurlServiceImpl struct{}

func NewCurlServiceImpl() domain.CurlService {
	return &CurlServiceImpl{}
}

func (s *CurlServiceImpl) ExecuteCurl(req types.CurlRequest) (types.CurlResponse, error) {
	client := &http.Client{}
	request, err := http.NewRequest(req.Method, req.URL, io.NopCloser(strings.NewReader(req.Body)))
	if err != nil {
		return types.CurlResponse{ErrMessage: err.Error()}, err
	}
	for k, v := range req.Headers {
		request.Header.Set(k, v)
	}

	resp, err := client.Do(request)
	if err != nil {
		return types.CurlResponse{ErrMessage: err.Error()}, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	headers := map[string]string{}
	for k, v := range resp.Header {
		headers[k] = v[0]
	}

	return types.CurlResponse{
		Status:  resp.StatusCode,
		Headers: headers,
		Body:    string(respBody),
	}, nil
}
