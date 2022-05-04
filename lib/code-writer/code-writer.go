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
	fmt.Printf("%v %v %v\n", command, segment, index)
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

			fmt.Fprintf(writer, "@%v\n", index)
			fmt.Fprintf(writer, "D=A\n") // D=10

			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "A=M\n") // *sp
			fmt.Fprintf(writer, "M=D\n")

			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "M=M+1\n") // sp++
			return
		}

		if segment == "local" || segment == "argument" || segment == "this" || segment == "that" || segment == "temp" {

			fmt.Fprintf(writer, "@%v\n", index)
			fmt.Fprintf(writer, "D=A\n")

			if segment == "temp" {
				fmt.Fprintf(writer, "@5\n")    // temp base address is 5
				fmt.Fprintf(writer, "D=D+A\n") // D = 5 + index
			} else {
				fmt.Fprintf(writer, "@%v\n", symbol)
				fmt.Fprintf(writer, "D=D+M\n") // D = segmentPointer + index
			}

			fmt.Fprintf(writer, "@addr\n")
			fmt.Fprintf(writer, "M=D\n") // addr = D

			fmt.Fprintf(writer, "@addr\n")
			fmt.Fprintf(writer, "A=M\n") // *addr
			fmt.Fprintf(writer, "D=M\n")

			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "A=M\n") // *sp
			fmt.Fprintf(writer, "M=D\n")

			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "M=M+1\n") // sp++

			return
		}

		if segment == "static" {

			fmt.Fprintf(writer, "@%v.%v\n", fileName, index)
			fmt.Fprintf(writer, "D=M\n") // D = fileName.index

			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "A=M\n") // *sp
			fmt.Fprintf(writer, "M=D\n")

			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "M=M+1\n") // sp++

			return
		}

		if segment == "pointer" {
			switch index {
			case 0:
				fmt.Fprintf(writer, "@THIS\n")
			case 1:
				fmt.Fprintf(writer, "@THAT\n")
			default:
				panic("invalid index of pointer segment")
			}

			fmt.Fprintf(writer, "D=M\n") // D = this/that

			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "A=M\n") // *sp
			fmt.Fprintf(writer, "M=D\n") // *sp = this/that

			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "M=M+1\n") // sp++

			return
		}

	} else if command == parser.C_POP.String() {
		if segment == "local" || segment == "argument" || segment == "this" || segment == "that" || segment == "temp" {

			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "M=M-1\n") // sp--

			fmt.Fprintf(writer, "@%v\n", index)
			fmt.Fprintf(writer, "D=A\n")

			if segment == "temp" {
				fmt.Fprintf(writer, "@5\n")    // temp base address is 5
				fmt.Fprintf(writer, "D=D+A\n") // D = 5 + index
			} else {
				fmt.Fprintf(writer, "@%v\n", symbol)
				fmt.Fprintf(writer, "D=D+M\n") // D = segmentPointer + index
			}

			fmt.Fprintf(writer, "@addr\n")
			fmt.Fprintf(writer, "M=D\n") // addr = D

			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "A=M\n") // *sp
			fmt.Fprintf(writer, "D=M\n")

			fmt.Fprintf(writer, "@addr\n")
			fmt.Fprintf(writer, "A=M\n") // *addr
			fmt.Fprintf(writer, "M=D\n")

			return
		}

		if segment == "static" {

			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "M=M-1\n") // sp--

			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "A=M\n") // *sp
			fmt.Fprintf(writer, "D=M\n")

			fmt.Fprintf(writer, "@%v.%v\n", fileName, index)
			fmt.Fprintf(writer, "M=D\n") // fileName.index = D

			return
		}

		if segment == "pointer" {
			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "M=M-1\n") // sp--

			fmt.Fprintf(writer, "@SP\n")
			fmt.Fprintf(writer, "A=M\n") // *sp
			fmt.Fprintf(writer, "D=M\n") // D = *sp

			switch index {
			case 0:
				fmt.Fprintf(writer, "@THIS\n")
			case 1:
				fmt.Fprintf(writer, "@THAT\n")
			default:
				panic("invalid index of pointer segment")
			}

			fmt.Fprintf(writer, "M=D\n") // this/that = D

			return
		}
	}

	panic("invalid command")
}

func WriteArithmetic(command string) {
	fmt.Fprintf(writer, "// %v\n", command)

	switch command {
	case "add":
		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n") // *sp
		fmt.Fprintf(writer, "D=M\n") // D = *sp

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n")   // *sp
		fmt.Fprintf(writer, "M=D+M\n") // *sp = D + *sp

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M+1\n") // sp++

	case "sub":
		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n") // *sp
		fmt.Fprintf(writer, "D=M\n") // D = *sp

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n")   // *sp
		fmt.Fprintf(writer, "M=M-D\n") // *sp = *sp - D

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M+1\n") // sp++

	case "eq":
		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n") // *sp
		fmt.Fprintf(writer, "D=M\n") // D = *sp

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n")   // *sp
		fmt.Fprintf(writer, "D=M-D\n") // D = *sp - D

		compareLogic("D;JEQ")

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M+1\n") // sp++

	case "lt":
		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n") // *sp
		fmt.Fprintf(writer, "D=M\n") // D = *sp

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n")   // *sp
		fmt.Fprintf(writer, "D=M-D\n") // D = *sp - D

		compareLogic("D;JLT")

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M+1\n") // sp++

	case "gt":
		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n") // *sp
		fmt.Fprintf(writer, "D=M\n") // D = *sp

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n")   // *sp
		fmt.Fprintf(writer, "D=M-D\n") // D = *sp - D

		compareLogic("D;JGT")

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M+1\n") // sp++

	case "and":
		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n") // *sp
		fmt.Fprintf(writer, "D=M\n") // D = *sp

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n")   // *sp
		fmt.Fprintf(writer, "M=D&M\n") // *sp = *sp & D

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M+1\n") // sp++

	case "or":
		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n") // *sp
		fmt.Fprintf(writer, "D=M\n") // D = *sp

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n")   // *sp
		fmt.Fprintf(writer, "M=D|M\n") // *sp = *sp | D

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M+1\n") // sp++

	case "not":
		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n")  // *sp
		fmt.Fprintf(writer, "M=!M\n") // *sp = !*sp

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M+1\n") // sp++

	case "neg":
		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n")  // *sp
		fmt.Fprintf(writer, "M=-M\n") // *sp = -*sp

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M+1\n") // sp++

	}
}

func compareLogic(jump string) {
	fmt.Fprintf(writer, "@LABEL_%v\n", suffix)
	fmt.Fprintf(writer, "%v\n", jump)

	fmt.Fprintf(writer, "@SP\n")
	fmt.Fprintf(writer, "A=M\n") // *sp
	fmt.Fprintf(writer, "M=0\n") // *sp = false
	fmt.Fprintf(writer, "@END_%v\n", suffix)
	fmt.Fprintf(writer, "0;JMP\n")

	fmt.Fprintf(writer, "(LABEL_%v)\n", suffix)
	fmt.Fprintf(writer, "@SP\n")
	fmt.Fprintf(writer, "A=M\n")  // *sp
	fmt.Fprintf(writer, "M=-1\n") // *sp = true

	fmt.Fprintf(writer, "(END_%v)\n", suffix)

	suffix++
}
