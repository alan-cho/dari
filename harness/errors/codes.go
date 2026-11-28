package errors

type Code string

const (
	EndOfInput         Code = "END_OF_INPUT"
	InputRead          Code = "INPUT_READ"
	APIRequest         Code = "API_REQUEST"
	MissingEnv         Code = "MISSING_ENV"
	InvalidEnv         Code = "INVALID_ENV"
	Truncated          Code = "TRUNCATED"
	Refused            Code = "REFUSED"
	UnsupportedToolUse Code = "UNSUPPORTED_TOOL_USE"
	UnexpectedStop     Code = "UNEXPECTED_STOP"
	BadRequest         Code = "BAD_REQUEST"
	Internal           Code = "INTERNAL"
	ServerRequest      Code = "SERVER_REQUEST"
	ServerResponse     Code = "SERVER_RESPONSE"
)

const unrecognizedDescription = "unrecognized error code"

var descriptions = map[Code]string{
	EndOfInput:         "input was closed",
	InputRead:          "failed to read from the terminal",
	APIRequest:         "request to the model failed",
	MissingEnv:         "required environment variable is not set",
	InvalidEnv:         "environment variable has an invalid value",
	Truncated:          "model reply was cut off at the token limit",
	Refused:            "model declined to respond",
	UnsupportedToolUse: "model requested a tool, which is not supported yet",
	UnexpectedStop:     "model stopped for an unhandled reason",
	BadRequest:         "request body is malformed",
	Internal:           "server failed unexpectedly",
	ServerRequest:      "failed to reach the server",
	ServerResponse:     "server returned an unreadable response",
}

func describe(code Code) string {
	if description, ok := descriptions[code]; ok {
		return description
	}
	return unrecognizedDescription
}
