package codewriter

import (
	"bufio"
	"fmt"
	"vm-translator/lib/parser"
)

var writer *bufio.Writer

func Init(w *bufio.Writer) {
	writer = w
}

func WritePushPop(command string, segment string, index int) {
	fmt.Printf("%v %v %v\n", command, segment, index)

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

			fmt.Fprintf(writer, "// %v %v %v\n", command, segment, index)
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
			fmt.Fprintf(writer, "// %v %v %v\n", command, segment, index)

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

	} else if command == parser.C_POP.String() {
		if segment == "local" || segment == "argument" || segment == "this" || segment == "that" || segment == "temp" {
			fmt.Fprintf(writer, "// %v %v %v\n", command, segment, index)

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
	}

	fmt.Println("skipped")
}

func WriteArithmetic(command string) {
	switch command {
	case "add":
		fmt.Fprintf(writer, "// %v\n", command)

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n") // *sp
		fmt.Fprintf(writer, "D=M\n") // D = *sp

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M-1\n") // sp--

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n")   // *sp
		fmt.Fprintf(writer, "D=D+M\n") // D = D + *sp

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n") // *sp
		fmt.Fprintf(writer, "M=D\n") // *sp = D

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M+1\n") // sp++

	case "sub":
		fmt.Fprintf(writer, "// %v\n", command)

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

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "A=M\n") // *sp
		fmt.Fprintf(writer, "M=D\n") // *sp = D

		fmt.Fprintf(writer, "@SP\n")
		fmt.Fprintf(writer, "M=M+1\n") // sp++

	}
}
