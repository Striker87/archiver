package cmd

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Striker87/archiver/lib/compression"
	"github.com/Striker87/archiver/lib/compression/vlc"
	"github.com/spf13/cobra"
)

const packedExtension = "vlc"

var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "Pack file",
	Run:   pack,
}

func init() {
	rootCmd.AddCommand(packCmd)
	packCmd.Flags().StringP("method", "m", "", "compression method: vlc")
	if err := packCmd.MarkFlagRequired("method"); err != nil {
		log.Fatal("packCmd.MarkFlagRequired error:", err)
	}
}

var ErrEmptyPath = errors.New("path to file is not specified")

func pack(cmd *cobra.Command, args []string) {
	filePath := args[0]
	if len(args) == 0 || filePath == "" {
		log.Fatal(ErrEmptyPath)
	}

	method := cmd.Flag("method").Value.String()
	var encoder compression.Encoder

	switch method {
	case "vlc":
		encoder = vlc.New()
	default:
		cmd.PrintErrf("unknown compress method")
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

	packed := encoder.Encode(string(data))
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
