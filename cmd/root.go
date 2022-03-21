package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "Simple archiver",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("failed to execute due error: %v", err)
	}
}
