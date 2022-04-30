package parser

type COMMAND_TYPES uint8

const (
	UNDEFINED COMMAND_TYPES = iota
	C_ARITHMETIC
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

var ARITHMETIC_COMMANDS = map[string]string{
	"add": "ADD",
	"sub": "SUB",
	"neg": "NEG",
	"eq":  "EQ",
	"gt":  "GT",
	"lt":  "LT",
	"and": "AND",
	"or":  "OR",
	"not": "NOT",
}

func (s COMMAND_TYPES) String() string {
	switch s {
	case C_ARITHMETIC:
		return "arithmetic"
	case C_PUSH:
		return "push"
	case C_POP:
		return "pop"
	case C_LABEL:
		return "label"
	case C_GOTO:
		return "goto"
	case C_IF:
		return "if"
	case C_FUNCTION:
		return "function"
	case C_RETURN:
		return "return"
	case C_CALL:
		return "call"
	}
	return "unknown"
}

// func (s ARITHMETIC_COMMANDS) String() string {
// 	switch s {
// 	case ADD:
// 		return "Arithmetic command: add"
// 	case SUB:
// 		return "Arithmetic command: sub"
// 	case NEG:
// 		return "Arithmetic command: neg"
// 	case EQ:
// 		return "Arithmetic command: eq"
// 	case GT:
// 		return "Arithmetic command: gt"
// 	case LT:
// 		return "Arithmetic command: lt"
// 	case AND:
// 		return "Arithmetic command: and"
// 	case OR:
// 		return "Arithmetic command: or"
// 	case NOT:
// 		return "Arithmetic command: not"
// 	}
// 	return "Arithmetic command: unknown"
// }

// func (s MEMORY_ACCESS_COMMANDS) String() string {
// 	switch s {
// 	case POP:
// 		return "Memory access command: pop"
// 	case PUSH:
// 		return "Memory access command: push"
// 	}
// 	return "Memory access command: unkown"
// }
