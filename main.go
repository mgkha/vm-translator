package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	codewriter "vm-translator/lib/code-writer"
	parser "vm-translator/lib/parser"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	inputFileName := os.Args[1]
	outputFilename := strings.TrimSuffix(inputFileName, filepath.Ext(inputFileName)) + ".asm"

	className := strings.TrimSuffix(filepath.Base(inputFileName), filepath.Ext(inputFileName))

	inFile, err := os.Open(inputFileName)
	check(err)
	defer inFile.Close()
	reader := bufio.NewReader(inFile)

	outFile, err := os.Create(outputFilename)
	check(err)
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)

	parser.Init(reader)
	codewriter.Init(writer, className)

	for parser.HasMoreLines() {
		parser.Advance()

		index, _ := strconv.Atoi(parser.Arg2())

		if parser.CommandType() == parser.C_PUSH || parser.CommandType() == parser.C_POP {
			codewriter.WritePushPop(parser.CommandType().String(), parser.Arg1(), index)
		} else if parser.CommandType() == parser.C_ARITHMETIC {
			codewriter.WriteArithmetic(parser.Arg1())
		}
	}
	writer.Flush()
}
