package bf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type Directory struct {
	FirstFileIndex   int
	FirstSubdirIndex int
	PrevIndex        int
	NextIndex        int
	ParentIndex      int
	Children         []*Directory
	Files            []*File
	Dirname          string
}

func (dir *Directory) ExtractDir(target string, bfPath string, incEmpty bool) {
	extFlag := true
	if len(dir.Children) == 0 && len(dir.Files) == 0 {
		extFlag = incEmpty
	}
	if extFlag {
		err := os.MkdirAll(target, os.ModePerm)
		if err != nil {
			panic(err)
		}

		for i := range len(dir.Children) {
			dir.Children[i].ExtractDir(target+"/"+ConvertToUnicode(dir.Children[i].Dirname), bfPath, incEmpty)
		}

		for i := range len(dir.Files) {
			dir.Files[i].WriteToDisk(target, bfPath)
		}
	} else {
		fmt.Printf("Directory %s is empty - skipping\n", dir.Dirname)
	}
}

func (dir *Directory) PrintDir() {
	fmt.Printf("%d,%s\n", len(dir.Children), dir.Dirname)
}

func (dir *Directory) PrintDirTree(depth int) {
	var spaces string = ""
	for range depth {
		spaces += " "
	}

	fmt.Printf("%s+ %s\n", spaces, dir.Dirname)

	for i := range len(dir.Files) {
		fmt.Printf("%s|- %s\n", spaces, dir.Files[i].Metadata.Filename)
	}

	for i := range len(dir.Children) {
		dir.Children[i].PrintDirTree(depth + 1)
	}
}

func MakeDirectory(inBytes []byte) *Directory {
	return &Directory{
		FirstFileIndex:   int(binary.LittleEndian.Uint32(inBytes[:4])),
		FirstSubdirIndex: int(binary.LittleEndian.Uint32(inBytes[4:8])),
		PrevIndex:        int(binary.LittleEndian.Uint32(inBytes[8:12])),
		NextIndex:        int(binary.LittleEndian.Uint32(inBytes[12:16])),
		ParentIndex:      int(binary.LittleEndian.Uint32(inBytes[16:20])),
		Children:         make([]*Directory, 0),
		Files:            make([]*File, 0),
		Dirname:          string(bytes.Trim(inBytes[20:], "\x00")),
	}
}
