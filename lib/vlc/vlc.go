package vlc

import (
	"log"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	chunksSize         = 8
	hexChunksSeparator = " "
)

func Encode(str string) string {
	// prepare text: M ->!m
	str = prepareText(str)

	// encode to binary: some text-> 10101000
	// split binary by chunks(8): bits to bytes -> '10101000 10101000 10101000'
	chunks := splitByChunks(encodeBin(str), chunksSize)

	// bytes to HEX -> '20 30 3C'
	return chunks.ToHex().ToString()
}

// splitByChunks splits binary string by chunks with given size
// i.g. '101010001010100010101000' -> '10101000 10101000 10101000'
func splitByChunks(bStr string, chunkSize int) BinaryChunks {
	strLen := utf8.RuneCountInString(bStr)
	chunksCount := strLen / chunkSize

	if strLen/chunkSize != 0 {
		chunksCount++
	}
	res := make(BinaryChunks, 0, chunksCount)

	var buf strings.Builder
	for i, ch := range bStr {
		buf.WriteString(string(ch))

		if (i+1)%chunkSize == 0 {
			res = append(res, BinaryChunk(buf.String()))
			buf.Reset()
		}
	}

	if buf.Len() != 0 {
		lastChunk := buf.String()
		lastChunk += strings.Repeat("0", chunkSize-len(lastChunk))
		res = append(res, BinaryChunk(lastChunk))
	}
	return res
}

// prepareText prepares text to be fit for encode:
// changes upper case letters to: ! + lower case letter
// i.g. My name is Ted ->! my name is !ted
func prepareText(str string) string {
	var buf strings.Builder

	for _, ch := range str {
		if unicode.IsUpper(ch) {
			buf.WriteRune('!')
			buf.WriteRune(unicode.ToLower(ch))
		} else {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}

// encodeBin encodes str into binary codes string without spaces
func encodeBin(str string) string {
	var buf strings.Builder

	for _, ch := range str {
		buf.WriteString(bin(ch))
	}

	return buf.String()
}

func bin(ch rune) string {
	table := getEncodingTable()

	res, ok := table[ch]
	if !ok {
		log.Fatalf("unknown character: %s", string(ch))
	}
	return res
}

func getEncodingTable() encodingTable {
	return encodingTable{
		' ': "11",
		't': "1001",
		'n': "10000",
		's': "0101",
		'r': "01000",
		'd': "00101",
		'!': "001000",
		'c': "000101",
		'm': "000011",
		'g': "0000100",
		'b': "0000010",
		'v': "00000001",
		'k': "0000000001",
		'q': "000000000001",
		'e': "101",
		'o': "10001",
		'a': "011",
		'i': "01001",
		'h': "0011",
		'l': "001001",
		'u': "00011",
		'f': "000100",
		'p': "0000101",
		'w': "0000011",
		'y': "0000001",
		'j': "000000001",
		'x': "00000000001",
		'z': "000000000000",
	}
}
func Decode(encodedText string) string {
	// hex chunks -> binary chunks
	// bChunks -> binary string
	bString := NewHexChunks(encodedText).ToBinary().Join()

	// build decoding tree
	dTree := getEncodingTable().DecodingTree()

	// bString (dTree) -> text
	return exportText(dTree.Decode(bString)) // My name id Ted -> !my name is !ted
}

// exportText is opposite to prepareText, it prepares decoded text to export
// changes: ! + <lower case letter> -> to upper case letter.
// i.g.: !my name is !ted -> My name is Ted
func exportText(str string) string {
	var (
		buf       strings.Builder
		isCapital bool
	)

	for _, ch := range str {
		if isCapital {
			buf.WriteRune(unicode.ToUpper(ch))
			isCapital = false
			continue
		}

		if ch == '!' {
			isCapital = true
			continue
		} else {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}
