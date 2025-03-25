package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ItsYogSothoth/jade-bf-extract-go/bf"
)

var mode = "extract"
var bfPath string
var extractPath string

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		PrintUsage()
		os.Exit(0)
	}

	ProcessCommandline(args)

	bfFile := bf.MakeBfFile(bfPath)

	switch mode {
	case "extract":
		if extractPath != "" {
			bfFile.ExtractDir(0, extractPath, false)
		} else {
			fmt.Printf("Missing extract path\n")
		}
	case "print-offset":
		bfFile.PrintOffsetArray()
	case "print-info":
		bfFile.PrintInfo()
	}
}

func ProcessCommandline(args []string) {
	for i := range args {
		switch args[i] {
		case "--print-offset-array":
			mode = "print-offset"
		case "--print-info":
			mode = "print-info"
		default:
			if(args[i][0:2] != "--") {
				bfPath = args[i]
				if (i + 1) != len(args) {
					i++
					extractPath = args[i]
				}
			}
		}
	}
}

func PrintUsage() {
	fmt.Printf("Usage: %s [options] <bf-file-path> <output-dir>\n", filepath.Base(os.Args[0]))
	fmt.Printf("\nOptions:\n --print-offset-array - prints offset array with resource keys in CSV format\n --print-info - prints BF file info\n")
}
