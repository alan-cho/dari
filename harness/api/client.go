package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	herrors "dari/harness/errors"
)

const MessagesPath = "/messages"

func Send(ctx context.Context, baseURL string, request SendRequest) (SendResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return SendResponse{}, herrors.New(herrors.ServerRequest, err)
	}
	endpoint, err := url.JoinPath(baseURL, MessagesPath)
	if err != nil {
		return SendResponse{}, herrors.New(herrors.ServerRequest, err)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return SendResponse{}, herrors.New(herrors.ServerRequest, err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return SendResponse{}, herrors.New(herrors.ServerRequest, err)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return SendResponse{}, herrors.New(errorCode(httpResponse.Body), fmt.Errorf("server responded %s", httpResponse.Status))
	}
	var response SendResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return SendResponse{}, herrors.New(herrors.ServerResponse, err)
	}
	return response, nil
}

func errorCode(body io.Reader) herrors.Code {
	var response ErrorResponse
	if err := json.NewDecoder(body).Decode(&response); err != nil || response.Code == "" {
		return herrors.ServerResponse
	}
	return response.Code
}
