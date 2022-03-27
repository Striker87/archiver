package vlc

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type (
	BinaryChunk   string
	BinaryChunks  []BinaryChunk
	encodingTable map[rune]string
)

func NewBinChunks(data []byte) BinaryChunks {
	res := make(BinaryChunks, 0, len(data))
	for _, code := range data {
		res = append(res, NewBinChunk(code))
	}

	return res
}

func NewBinChunk(code byte) BinaryChunk {
	return BinaryChunk(fmt.Sprintf("%08b", code))
}

func (b BinaryChunks) Bytes() []byte {
	res := make([]byte, 0, len(b))
	for _, bc := range b {
		res = append(res, bc.Byte())
	}

	return res
}

func (b BinaryChunk) Byte() byte {
	num, err := strconv.ParseUint(string(b), 2, chunksSize)
	if err != nil {
		log.Fatalf("failed to parse binary chunk: %v", err)
	}

	return byte(num)
}

// Join joins chunks into one line and returns as string
func (b BinaryChunks) Join() string {
	var buf strings.Builder
	for _, bc := range b {
		buf.WriteString(string(bc))
	}

	return buf.String()
}
