package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ItsYogSothoth/jade-bf-extract-go/bf"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		PrintUsage()
		os.Exit(0)
	}

	bfFile := bf.MakeBfFile(args[0])

	bfFile.ExtractDir(0, args[1], false)
}

func PrintUsage() {
	fmt.Printf("Usage: %s <bf-file-path> <output-dir>\n", filepath.Base(os.Args[0]))
}
