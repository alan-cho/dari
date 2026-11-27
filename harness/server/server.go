package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"dari/harness/agent"
	"dari/harness/api"
	herrors "dari/harness/errors"
)

func NewHandler(client agent.Client) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+api.MessagesPath, sendHandler(client))
	return mux
}

func sendHandler(client agent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := send(r.Context(), client, r.Body)
		if err != nil {
			log.Println(err)
			status, body := errorResponse(err)
			writeJSON(w, status, body)
			return
		}
		writeJSON(w, http.StatusOK, response)
	}
}

func send(ctx context.Context, client agent.Client, body io.Reader) (api.SendResponse, error) {
	request, err := decode(body)
	if err != nil {
		return api.SendResponse{}, err
	}
	messages, reply, usage, err := agent.Send(ctx, client, request.Messages, request.Input)
	if err != nil {
		return api.SendResponse{}, err
	}
	return api.SendResponse{Messages: messages, Reply: reply, Usage: usage}, nil
}

func decode(body io.Reader) (api.SendRequest, error) {
	var request api.SendRequest
	if err := json.NewDecoder(body).Decode(&request); err != nil {
		return api.SendRequest{}, herrors.New(herrors.BadRequest, err)
	}
	return request, validate(request)
}

func validate(request api.SendRequest) error {
	if strings.TrimSpace(request.Input) == "" {
		return herrors.New(herrors.BadRequest, errors.New("input is empty"))
	}
	for index, message := range request.Messages {
		if message.Role != api.User && message.Role != api.Assistant {
			return herrors.New(herrors.BadRequest, fmt.Errorf("message %d has unknown role %q", index, message.Role))
		}
	}
	return nil
}

func errorResponse(err error) (int, api.ErrorResponse) {
	var harnessErr *herrors.Error
	if !errors.As(err, &harnessErr) {
		harnessErr = herrors.New(herrors.Internal, err).(*herrors.Error)
	}
	return statusFor(harnessErr.Code), api.ErrorResponse{Code: harnessErr.Code, Description: harnessErr.Description()}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Println(err)
	}
}
