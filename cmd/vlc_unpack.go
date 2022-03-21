package cmd

import (
	"arhiver/lib/vlc"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// TODO: take extension from file
const unpackedExtension = "txt"

var vlcUnpackCmd = &cobra.Command{
	Use:   "vlc",
	Short: "Unpack file using variable-length code",
	Run:   unpack,
}

func unpack(_ *cobra.Command, args []string) {
	filePath := args[0]
	if len(args) == 0 || filePath == "" {
		log.Fatal(ErrEmptyPath)
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

	packed := vlc.Decode(string(data))
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

func init() {
	unpackCmd.AddCommand(vlcUnpackCmd)
}
