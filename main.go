package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
)

const appVersion = "0.4"

const appVersionString = "wavdump " + appVersion

const usage = `
The wavdump tool shows the header structure of WAV or WAVE64 files.

Usage:
  wavdump <input_files>...

Options:
  -h --help             Show this screen.
  --version             Show version of the program.
`

type cmdArgs struct {
	Input_files []string
}

func dataToStr(data []byte) string {
	if data == nil {
		return ""
	}

	res := ""
	for i, b := range data {
		if i > 0 && i%4 == 0 {
			res += " "
		}

		res += fmt.Sprintf("%02x ", b)
	}

	return res
}

func printLines(lines []Line) error {

	for _, line := range lines {

		if line.data == nil {
			fmt.Println()
			continue
		}

		dataStr := dataToStr(line.data)
		value := line.value
		if value == nil {
			value = ""
		}

		fmt.Printf("%08x", line.address)
		fmt.Printf("%57s", dataStr)
		fmt.Printf("%15v", value)
		fmt.Printf("    %s\n", line.description)
	}
	return nil
}

func main() {
	args := cmdArgs{}

	parser := docopt.Parser{}
	opts, _ := parser.ParseArgs(usage, os.Args[1:], appVersionString)

	if err := opts.Bind(&args); err != nil {
		fmt.Println(err)
		parser.HelpHandler(err, usage)
		os.Exit(1)
	}

	res := 0
	for _, file := range args.Input_files {
		fmt.Println("\n", file)
		parser := Parser{}
		err := parser.parse(file)

		if err == nil {
			err = printLines(parser.lines)
		}

		if err != nil {
			fmt.Println(err)
			res = 2
		}


		fmt.Printf("\naudio data hash:    %x\n", parser.hash)
	}

	os.Exit(res)
}
