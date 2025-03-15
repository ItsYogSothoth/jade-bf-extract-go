package bf

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type BigFile struct {
	Path                    string
	Version                 []byte
	FileCount               int
	DirCount                int
	OffsetTableLength       int
	InitialKey              []byte
	OffsetTableOffset       int
	FileMetadataTableOffset int
	DirMetadataTableOffset  int

	OffsetArray []Offset
	Files       []File
	Dirs        []Directory
}

func (bf *BigFile) ProcessHeader() {
	f, err := os.Open(bf.Path)
	if err != nil {
		panic(err)
	}

	f.Seek(4, io.SeekStart)
	f.Read(bf.Version)

	fileCountBytes := make([]byte, 4)
	f.Read(fileCountBytes)
	bf.FileCount = int(binary.LittleEndian.Uint32(fileCountBytes))

	dirCountBytes := make([]byte, 4)
	f.Read(dirCountBytes)
	bf.DirCount = int(binary.LittleEndian.Uint32(dirCountBytes))

	f.Seek(16, io.SeekCurrent)

	offsetTableLenthBytes := make([]byte, 4)
	f.Read(offsetTableLenthBytes)
	bf.OffsetTableLength = int(binary.LittleEndian.Uint32(offsetTableLenthBytes))

	f.Seek(4, io.SeekCurrent)

	f.Read(bf.InitialKey)

	f.Seek(8, io.SeekCurrent)

	offsetTableOffsetBytes := make([]byte, 4)
	f.Read(offsetTableOffsetBytes)
	bf.OffsetTableOffset = int(binary.LittleEndian.Uint32(offsetTableOffsetBytes))

	f.Close()

	bf.FileMetadataTableOffset = bf.OffsetTableOffset + bf.OffsetTableLength*8
	bf.DirMetadataTableOffset = bf.FileMetadataTableOffset + bf.OffsetTableLength*0x54
}

func (bf *BigFile) PopulateOffsetArray() {
	f, err := os.Open(bf.Path)
	if err != nil {
		panic(err)
	}

	bf.OffsetArray = make([]Offset, bf.OffsetTableLength)
	f.Seek(int64(bf.OffsetTableOffset), io.SeekStart)

	for i := range bf.OffsetTableLength {
		fileDataOffsetBytes := make([]byte, 4)
		resourceKeyBytes := make([]byte, 4)

		f.Read(fileDataOffsetBytes)
		f.Read(resourceKeyBytes)

		bf.OffsetArray[i] = Offset{
			FileDataOffset: int(binary.LittleEndian.Uint32(fileDataOffsetBytes)),
			ResourceKey:    int(binary.LittleEndian.Uint32(resourceKeyBytes)),
		}
	}

	f.Close()
}

func (bf *BigFile) PopulateDirArray() {
	f, err := os.Open(bf.Path)
	if err != nil {
		panic(err)
	}

	bf.Dirs = make([]Directory, bf.DirCount)
	f.Seek(int64(bf.DirMetadataTableOffset), io.SeekStart)

	for i := range bf.DirCount {
		dirBytes := make([]byte, 84)
		f.Read(dirBytes)
		bf.Dirs[i] = *MakeDirectory(dirBytes)

		if bf.Dirs[i].ParentIndex != 0xFFFFFFFF {
			bf.Dirs[bf.Dirs[i].ParentIndex].Children = append(bf.Dirs[bf.Dirs[i].ParentIndex].Children, &bf.Dirs[i])
		}
	}

	f.Close()
}

func (bf *BigFile) PopulateFileArray() {
	f, err := os.Open(bf.Path)
	if err != nil {
		panic(err)
	}

	bf.Files = make([]File, bf.FileCount)
	f.Seek(int64(bf.FileMetadataTableOffset), io.SeekStart)

	for i := range bf.FileCount {
		metadataBytes := make([]byte, 84)
		f.Read(metadataBytes)
		bf.Files[i] = File{
			Metadata:    *MakeFileMetadata(metadataBytes),
			DataOffset:  bf.OffsetArray[i].FileDataOffset,
			ResourceKey: bf.OffsetArray[i].ResourceKey,
		}

		bf.Dirs[bf.Files[i].Metadata.DirIndex].Files = append(bf.Dirs[bf.Files[i].Metadata.DirIndex].Files, &bf.Files[i])
	}

	f.Close()
}

func (bf *BigFile) ExtractFile(idx int, target string) {
	bf.Files[idx].WriteToDisk(target, bf.Path)
}

func (bf *BigFile) ExtractDir(idx int, target string, incEmpty bool) {
	bf.Dirs[idx].ExtractDir(target, bf.Path, incEmpty)
}

func (bf *BigFile) PrintInfo() {
	fmt.Printf(
		`BF path: %s
BF version: %x
File count: %d
Dir count: %d
Offset table length: %d
Initial key: %x
Offset table offset: %x
File metadata table offset: %x
Dir metadata table offset: %x
`,
		bf.Path,
		bf.Version,
		bf.FileCount,
		bf.DirCount,
		bf.OffsetTableLength,
		bf.InitialKey,
		bf.OffsetTableOffset,
		bf.FileMetadataTableOffset,
		bf.DirMetadataTableOffset,
	)
}

func (bf *BigFile) PrintOffsetArray() {
	fmt.Printf("index,file_data_offset,resouce_key,filename\n")
	for i := range len(bf.OffsetArray) {
		fmt.Printf("%d,%x,%x,%s\n", i, bf.OffsetArray[i].FileDataOffset, bf.OffsetArray[i].ResourceKey, bf.Files[i].Metadata.Filename)
	}
}

func MakeBfFile(path string) *BigFile {
	bf := BigFile{
		Path:       path,
		Version:    make([]byte, 4),
		InitialKey: make([]byte, 4),
	}
	bf.ProcessHeader()
	bf.PopulateOffsetArray()
	bf.PopulateDirArray()
	bf.PopulateFileArray()

	return &bf
}

func ConvertToUnicode(name string) string {
	decoder := charmap.ISO8859_15.NewDecoder()
	result, _, err := transform.String(decoder, name)
	if err != nil {
		return name
	}
	return result
}
