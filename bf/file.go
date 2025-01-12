package bf

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type File struct {
	Metadata    FileMetadata
	DataOffset  int
	ResourceKey int
}

func (file *File) PrintFile() {
	fmt.Printf("%x,%d,%s\n", file.ResourceKey, file.Metadata.DirIndex, file.Metadata.Filename)
}

func (file *File) WriteToDisk(target string, bfPath string) {
	err := os.MkdirAll(target, os.ModePerm)
	if err != nil {
		panic(err)
	}

	bf, err := os.Open(bfPath)
	if err != nil {
		panic(err)
	}

	bf.Seek(int64(file.DataOffset), io.SeekStart)
	fileSizeBytes := make([]byte, 4)
	bf.Read(fileSizeBytes)
	fileSize := int(binary.LittleEndian.Uint32(fileSizeBytes))
	if fileSize != file.Metadata.Size {
		panic("Filesize mismatch")
	}

	fileBytes := make([]byte, fileSize)
	bf.Read(fileBytes)
	bf.Close()

	err = os.WriteFile(target+"/"+ConvertToUnicode(file.Metadata.Filename), fileBytes, 0644)
	if err != nil {
		panic(err)
	}
}
