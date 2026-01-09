package cmd

import (
	"fmt"
	"os"

	"github.com/antonpodkur/remember/internal/storage"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Long:  `Print all note names alphabetically (without .md extension)`,
	Args:  cobra.NoArgs,
	Run:   runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	notes, err := storage.ListNotes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	for _, note := range notes {
		fmt.Println(note)
	}
}
