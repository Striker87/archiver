package cmd

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Striker87/archiver/lib/vlc"
	"github.com/spf13/cobra"
)

const packedExtension = "vlc"

var vlcPackCmd = &cobra.Command{
	Use:   "vlc",
	Short: "Pack file using variable-length code",
	Run:   pack,
}

var ErrEmptyPath = errors.New("path to file is not specified")

func pack(_ *cobra.Command, args []string) {
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

	packed := vlc.Encode(string(data))
	err = os.WriteFile(packedFileName(filePath), packed, 0644)
	if err != nil {
		log.Fatalf("failed to WriteFile, %s due error: %v", filePath, err)
	}
}

func packedFileName(path string) string {
	// /path/fileName.txt ->fileName.vlc
	fileName := filepath.Base(path)

	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(fileName)) + "." + packedExtension
}

func init() {
	packCmd.AddCommand(vlcPackCmd)
}
