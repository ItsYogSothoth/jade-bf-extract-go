package bf

import (
	"bytes"
	"encoding/binary"
)

type FileMetadata struct {
	Size      int
	PrevIndex int
	NextIndex int
	DirIndex  int
	Timestamp int
	Filename  string
}

func MakeFileMetadata(inBytes []byte) *FileMetadata {
	return &FileMetadata{
		Size:      int(binary.LittleEndian.Uint32(inBytes[:4])),
		PrevIndex: int(binary.LittleEndian.Uint32(inBytes[4:8])),
		NextIndex: int(binary.LittleEndian.Uint32(inBytes[8:12])),
		DirIndex:  int(binary.LittleEndian.Uint32(inBytes[12:16])),
		Timestamp: int(binary.LittleEndian.Uint32(inBytes[16:20])),
		Filename:  string(bytes.Trim(inBytes[20:], "\x00")),
	}
}
