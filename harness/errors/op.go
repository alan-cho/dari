package errors

import (
	"runtime"
	"strings"
)

const unknownOp = "unknown"

func callerOp(skip int) string {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return unknownOp
	}
	function := runtime.FuncForPC(pc)
	if function == nil {
		return unknownOp
	}
	return shortName(function.Name())
}

func shortName(qualifiedName string) string {
	return qualifiedName[strings.LastIndex(qualifiedName, "/")+1:]
}
