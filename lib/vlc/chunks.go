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
	HexChunk      string
	HexChunks     []HexChunk
)

func (h HexChunks) ToString() string {
	switch len(h) {
	case 0:
		return ""
	case 1:
		return string(h[0])
	}

	var buf strings.Builder
	buf.WriteString(string(h[0]))

	for _, hc := range h[1:] {
		buf.WriteString(hexChunksSeparator)
		buf.WriteString(string(hc))
	}
	return buf.String()
}

func (b BinaryChunks) ToHex() HexChunks {
	res := make(HexChunks, 0, len(b))

	for _, chunk := range b {
		// chunk -> HexChunk
		hexChunk := chunk.ToHex()
		res = append(res, hexChunk)
	}
	return res
}

func (b BinaryChunk) ToHex() HexChunk {
	num, err := strconv.ParseUint(string(b), 2, chunksSize)
	if err != nil {
		log.Fatalf("failed to parse binary chunk due error: %v", err)
	}

	res := strings.ToUpper(fmt.Sprintf("%x", num))

	if len(res) == 1 {
		res = "0" + res
	}

	return HexChunk(res)
}

func NewHexChunks(str string) HexChunks {
	parts := strings.Split(str, hexChunksSeparator)
	res := make(HexChunks, 0, len(parts))

	for _, part := range parts {
		res = append(res, HexChunk(part))
	}
	return res
}

func (h HexChunks) ToBinary() BinaryChunks {
	res := make(BinaryChunks, 0, len(h))

	for _, chunk := range h {
		bChunk := chunk.ToBinary()
		res = append(res, bChunk)
	}
	return res
}

func (h HexChunk) ToBinary() BinaryChunk {
	num, err := strconv.ParseUint(string(h), 16, chunksSize)
	if err != nil {
		log.Fatalf("failed to parse HexChunk ToBinary() %v due error: %v", h, err)
	}
	res := fmt.Sprintf("%08b", num)
	return BinaryChunk(res)
}

// Join joins chunks into one line and returns as string
func (b BinaryChunks) Join() string {
	var buf strings.Builder
	for _, bc := range b {
		buf.WriteString(string(bc))
	}
	return buf.String()
}
