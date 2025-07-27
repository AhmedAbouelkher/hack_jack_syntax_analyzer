package main

import (
	"fmt"
	"runtime"
	"strings"
)

type AnalyzerError struct {
	Err     error
	Line    string
	LineNum int
	Stack   string
}

func (e *AnalyzerError) Error() string {
	return e.Err.Error()
}

func NewTokenErr(token Token, msg string, args ...any) *AnalyzerError {
	return &AnalyzerError{
		Err:     fmt.Errorf(msg, args...),
		Line:    token.tokenValue,
		LineNum: token.lineNum,
		Stack:   string(getStack()),
	}
}

// remove the first 3 lines of the stack trace
func getStack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			// remove the first 3 lines
			lines := strings.Split(string(buf[:n]), "\n")
			linesN := 3*2 + 1
			lines = lines[linesN:]
			return []byte(strings.Join(lines, "\n"))
		}
		buf = make([]byte, 2*len(buf))
	}
}
