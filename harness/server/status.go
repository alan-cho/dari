package server

import (
	"net/http"

	herrors "dari/harness/errors"
)

var statuses = map[herrors.Code]int{
	herrors.BadRequest:         http.StatusBadRequest,
	herrors.Truncated:          http.StatusUnprocessableEntity,
	herrors.Refused:            http.StatusUnprocessableEntity,
	herrors.UnsupportedToolUse: http.StatusUnprocessableEntity,
	herrors.APIRequest:         http.StatusBadGateway,
	herrors.UnexpectedStop:     http.StatusBadGateway,
}

func statusFor(code herrors.Code) int {
	if status, ok := statuses[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}
