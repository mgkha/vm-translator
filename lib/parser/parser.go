package parser

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

var reader *bufio.Reader
var currentCommand []string

func Init(r *bufio.Reader) {
	reader = r
}

func HasMoreLines() bool {
	p, err := reader.Peek(1)
	return err != io.EOF || len(p) > 0
}

func Advance() {
	line, _, _ := reader.ReadLine()
	str := string(line)

	parsed := whiteSpaceParser(str)
	if len(parsed) > 0 {
		currentCommand = strings.Split(str, " ")
		return
	}
	Advance()
}

func CommandType() (commandType COMMAND_TYPES) {
	_, arithmetic := ARITHMETIC_COMMANDS[currentCommand[0]]

	switch {
	case currentCommand[0] == C_PUSH.String():
		commandType = C_PUSH
	case currentCommand[0] == C_POP.String():
		commandType = C_POP
	case arithmetic:
		commandType = C_ARITHMETIC
	default:
		panic("error parsing command")
	}
	return
}

func Arg1() (arg1 string) {
	if CommandType() == C_POP || CommandType() == C_PUSH {
		arg1 = currentCommand[1]
	} else if CommandType() == C_ARITHMETIC {
		arg1 = currentCommand[0]
	}
	return
}

func Arg2() (arg2 string) {
	if CommandType() == C_POP || CommandType() == C_PUSH || CommandType() == C_FUNCTION || CommandType() == C_CALL {
		arg2 = currentCommand[2]
	}
	return
}

// Parse the given instruction string to remove
// - Empty lines / indentation
// - Line comments
// - In-line comments
func whiteSpaceParser(s string) string {
	r := regexp.MustCompile(`//.*|\s`)
	return r.ReplaceAllString(s, "")
}
