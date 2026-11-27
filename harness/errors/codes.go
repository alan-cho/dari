package errors

type Code string

const (
	EndOfInput         Code = "END_OF_INPUT"
	InputRead          Code = "INPUT_READ"
	APIRequest         Code = "API_REQUEST"
	MissingModel       Code = "MISSING_MODEL"
	Truncated          Code = "TRUNCATED"
	Refused            Code = "REFUSED"
	UnsupportedToolUse Code = "UNSUPPORTED_TOOL_USE"
	UnexpectedStop     Code = "UNEXPECTED_STOP"
)

const unrecognizedDescription = "unrecognized error code"

var descriptions = map[Code]string{
	EndOfInput:         "input was closed",
	InputRead:          "failed to read from the terminal",
	APIRequest:         "request to the model failed",
	MissingModel:       "MODEL environment variable is not set",
	Truncated:          "model reply was cut off at the token limit",
	Refused:            "model declined to respond",
	UnsupportedToolUse: "model requested a tool, which is not supported yet",
	UnexpectedStop:     "model stopped for an unhandled reason",
}

func describe(code Code) string {
	if description, ok := descriptions[code]; ok {
		return description
	}
	return unrecognizedDescription
}
