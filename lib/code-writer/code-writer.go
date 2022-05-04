package codewriter

import (
	"bufio"
	"fmt"
	"vm-translator/lib/parser"
)

var writer *bufio.Writer
var fileName string
var suffix uint16

func Init(w *bufio.Writer, name string) {
	writer = w
	fileName = name

}

func WritePushPop(command string, segment string, index int) {
	// fmt.Printf("%v %v %v\n", command, segment, index)
	fmt.Fprintf(writer, "// %v %v %v\n", command, segment, index)

	var symbol string
	switch segment {
	case "local":
		symbol = "LCL"
	case "argument":
		symbol = "ARG"
	case "this":
		symbol = "THIS"
	case "that":
		symbol = "THAT"
	}

	if command == parser.C_PUSH.String() {
		if segment == "constant" {

			loadD_Value(index) // D = index

			pushD_To_Stack()

			incrementStackPointer()

			return
		}

		if segment == "local" || segment == "argument" || segment == "this" || segment == "that" || segment == "temp" {

			loadD_Value(index)

			if segment == "temp" {
				fmt.Fprintf(writer, "@5\n")    // temp base address is 5
				fmt.Fprintf(writer, "D=D+A\n") // D = 5 + index
			} else {
				fmt.Fprintf(writer, "@%v\n", symbol)
				fmt.Fprintf(writer, "D=D+M\n") // D = segmentPointer + D (index)
			}

			loadAddr_From_D("segmentAddr") // segmentAddr = D

			loadA_From_Addr("segmentAddr") // *segmentAddr

			fmt.Fprintf(writer, "D=M\n")

			pushD_To_Stack()

			incrementStackPointer()

			return
		}

		if segment == "static" {

			loadD_Addr(fmt.Sprintf("%v.%v", fileName, index)) // D = @fileName.index

			pushD_To_Stack()

			incrementStackPointer()

			return
		}

		if segment == "pointer" {
			switch index {
			case 0:
				loadD_Addr("THIS")
			case 1:
				loadD_Addr("THAT")
			default:
				panic("invalid index of pointer segment")
			}

			pushD_To_Stack() // *sp = D (THIS/THAT)

			incrementStackPointer()

			return
		}

	} else if command == parser.C_POP.String() {
		if segment == "local" || segment == "argument" || segment == "this" || segment == "that" || segment == "temp" {

			decrementStackPointer()

			loadD_Value(index)

			if segment == "temp" {
				fmt.Fprintf(writer, "@5\n")    // temp base address is 5
				fmt.Fprintf(writer, "D=D+A\n") // D = 5 + index
			} else {
				fmt.Fprintf(writer, "@%v\n", symbol)
				fmt.Fprintf(writer, "D=D+M\n") // D = segmentPointer + D (index)
			}

			loadAddr_From_D("segmentAddr") // segmentAddr = D

			loadD_From_Stack()

			loadA_From_Addr("segmentAddr") // *segmentAddr

			fmt.Fprintf(writer, "M=D\n")

			return
		}

		if segment == "static" {

			decrementStackPointer()

			loadD_From_Stack()

			loadAddr_From_D(fmt.Sprintf("%v.%v", fileName, index))

			return
		}

		if segment == "pointer" {
			decrementStackPointer()

			loadD_From_Stack()

			switch index {
			case 0:
				loadAddr_From_D("THIS")
			case 1:
				loadAddr_From_D("THAT")
			default:
				panic("invalid index of pointer segment")
			}

			return
		}
	}

	panic("invalid command")
}

func WriteArithmetic(command string) {
	fmt.Fprintf(writer, "// %v\n", command)

	switch command {
	case "add":
		decrementStackPointer()

		loadD_From_Stack() // D = *sp

		decrementStackPointer()

		loadA_From_Stack()

		fmt.Fprintf(writer, "M=D+M\n") // *sp = D + *sp

		incrementStackPointer()

	case "sub":
		decrementStackPointer()

		loadD_From_Stack() // D = *sp

		decrementStackPointer()

		loadA_From_Stack()

		fmt.Fprintf(writer, "M=M-D\n") // *sp = *sp - D

		incrementStackPointer()

	case "eq":
		decrementStackPointer()

		loadD_From_Stack() // D = *sp

		decrementStackPointer()

		loadA_From_Stack()

		fmt.Fprintf(writer, "D=M-D\n") // D = *sp - D

		compareLogic("D;JEQ")

		incrementStackPointer()

	case "lt":
		decrementStackPointer()

		loadD_From_Stack() // D = *sp

		decrementStackPointer()

		loadA_From_Stack()

		fmt.Fprintf(writer, "D=M-D\n") // D = *sp - D

		compareLogic("D;JLT")

		incrementStackPointer()

	case "gt":
		decrementStackPointer()

		loadD_From_Stack() // D = *sp

		decrementStackPointer()

		loadA_From_Stack()

		fmt.Fprintf(writer, "D=M-D\n") // D = *sp - D

		compareLogic("D;JGT")

		incrementStackPointer()

	case "and":
		decrementStackPointer()

		loadD_From_Stack() // D = *sp

		decrementStackPointer()

		loadA_From_Stack()

		fmt.Fprintf(writer, "M=D&M\n") // *sp = *sp & D

		incrementStackPointer()

	case "or":
		decrementStackPointer()

		loadD_From_Stack() // D = *sp

		decrementStackPointer()

		loadA_From_Stack()

		fmt.Fprintf(writer, "M=D|M\n") // *sp = *sp | D

		incrementStackPointer()

	case "not":
		decrementStackPointer()

		loadA_From_Stack()

		fmt.Fprintf(writer, "M=!M\n") // *sp = !*sp

		incrementStackPointer()

	case "neg":
		decrementStackPointer()

		loadA_From_Stack()

		fmt.Fprintf(writer, "M=-M\n") // *sp = -*sp

		incrementStackPointer()

	}
}

func compareLogic(jump string) {
	fmt.Fprintf(writer, "@LABEL_%v\n", suffix)
	fmt.Fprintf(writer, "%v\n", jump)

	loadA_From_Stack()

	fmt.Fprintf(writer, "M=0\n") // *sp = false
	fmt.Fprintf(writer, "@END_%v\n", suffix)
	fmt.Fprintf(writer, "0;JMP\n")

	fmt.Fprintf(writer, "(LABEL_%v)\n", suffix)

	loadA_From_Stack()

	fmt.Fprintf(writer, "M=-1\n") // *sp = true

	fmt.Fprintf(writer, "(END_%v)\n", suffix)

	suffix++
}

func incrementStackPointer() {
	fmt.Fprintf(writer, "@SP\n")
	fmt.Fprintf(writer, "M=M+1\n") // sp++
}

func decrementStackPointer() {
	fmt.Fprintf(writer, "@SP\n")
	fmt.Fprintf(writer, "M=M-1\n") // sp--
}

func loadD_Value(value int) {
	fmt.Fprintf(writer, "@%v\n", value)
	fmt.Fprintf(writer, "D=A\n") // D = value
}

func loadD_Addr(addr string) {
	fmt.Fprintf(writer, "@%v\n", addr)
	fmt.Fprintf(writer, "D=M\n") // D = addr
}

// addr = D
func loadAddr_From_D(addr string) {
	fmt.Fprintf(writer, "@%v\n", addr)
	fmt.Fprintf(writer, "M=D\n")
}

// A = addr
func loadA_From_Addr(addr string) {
	fmt.Fprintf(writer, "@%v\n", addr)
	fmt.Fprintf(writer, "A=M\n")
}

// A = sp
func loadA_From_Stack() {
	fmt.Fprintf(writer, "@SP\n")
	fmt.Fprintf(writer, "A=M\n") // A = sp
}

// D = *sp
func loadD_From_Stack() {
	loadA_From_Stack()
	fmt.Fprintf(writer, "D=M\n")
}

// *sp = D
func pushD_To_Stack() {
	loadA_From_Stack()
	fmt.Fprintf(writer, "M=D\n")
}
