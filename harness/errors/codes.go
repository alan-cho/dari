package errors

type Code string

const (
	EndOfInput Code = "END_OF_INPUT"
	InputRead  Code = "INPUT_READ"
	APIRequest Code = "API_REQUEST"
)

const unrecognizedDescription = "unrecognized error code"

var descriptions = map[Code]string{
	EndOfInput: "input was closed",
	InputRead:  "failed to read from the terminal",
	APIRequest: "request to the model failed",
}

func describe(code Code) string {
	if description, ok := descriptions[code]; ok {
		return description
	}
	return unrecognizedDescription
}
