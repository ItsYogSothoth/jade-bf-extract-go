package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ItsYogSothoth/jade-bf-extract-go/bf"
)

var mode = "extract"
var bfPath string
var extractPath string

func init() {
	flag.Usage = PrintUsage
	offsetArrayPtr := flag.Bool("print-offset-array", false, "prints offset array with resource keys in CSV format")
	printInfoPtr := flag.Bool("print-info", false, "prints BF file info")

	flag.Parse()

	if *offsetArrayPtr {
		mode = "print-offset"
	}

	if *printInfoPtr {
		mode = "print-info"
	}

	args := flag.Args()
	if len(args) < 1 || (len(args) < 2 && mode == "extract") {
		PrintUsage()
		os.Exit(0)
	}

	bfPath = flag.Args()[0]
	if mode == "extract" {
		extractPath = flag.Args()[1]
	}
}

func main() {
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

func PrintUsage() {
	fmt.Printf("Usage: %s [options] <bf-file-path> <output-dir>\n", filepath.Base(os.Args[0]))
	fmt.Printf("\nOptions:\n --print-offset-array - prints offset array with resource keys in CSV format\n --print-info - prints BF file info\n")
}
