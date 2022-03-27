package cmd

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Striker87/archiver/lib/compression"
	"github.com/Striker87/archiver/lib/compression/vlc"
	"github.com/spf13/cobra"
)

// TODO: take extension from file
const unpackedExtension = "txt"

var unpackCmd = &cobra.Command{
	Use:   "unpack",
	Short: "Unpack file",
	Run:   unpack,
}

func init() {
	rootCmd.AddCommand(unpackCmd)

	unpackCmd.Flags().StringP("method", "m", "", "decompression method: vlc")
	if err := unpackCmd.MarkFlagRequired("method"); err != nil {
		log.Fatal("packCmd.MarkFlagRequired error:", err)
	}
}

func unpack(cmd *cobra.Command, args []string) {
	filePath := args[0]
	if len(args) == 0 || filePath == "" {
		log.Fatal(ErrEmptyPath)
	}

	method := cmd.Flag("method").Value.String()
	var decoder compression.Decoder

	switch method {
	case "vlc":
		decoder = vlc.New()
	default:
		cmd.PrintErrf("unknown decompress method")
	}

	r, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file %s due error: %v", filePath, err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		log.Fatalf("failed to io.ReadAll due error: %v", err)
	}

	packed := decoder.Decode(data)
	err = os.WriteFile(unpackedFileName(filePath), []byte(packed), 0644)
	if err != nil {
		log.Fatalf("failed to WriteFile, %s due error: %v", filePath, err)
	}
}

func unpackedFileName(path string) string {
	// /path/fileName.txt ->fileName.vlc
	fileName := filepath.Base(path)

	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(fileName)) + "." + unpackedExtension
}
