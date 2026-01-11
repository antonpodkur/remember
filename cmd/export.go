package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/antonpodkur/remember/internal/storage"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export <name|latest>",
	Short: "Export note in minified format for LLM context",
	Long:  `Export a note in delimiter-separated format (timestamp|content) with multiline content joined to single lines`,
	Args:  cobra.ExactArgs(1),
	Run:   runExport,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		notes, _ := storage.ListNotes()
		if len(notes) > 0 {
			notes = append([]string{"latest"}, notes...)
		}
		return notes, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) {
	name := args[0]

	if name == "latest" {
		latestName, err := storage.GetLatestNote()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
		name = latestName
	}

	if !storage.NoteExists(name) {
		fmt.Fprintf(os.Stderr, "Error: note '%s' not found\n", name)
		os.Exit(1)
	}

	entries, err := storage.ParseNoteEntries(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		content := strings.ReplaceAll(entry.Content, "\n", " ")
		content = strings.Join(strings.Fields(content), " ")
		fmt.Printf("%s|%s\n", entry.Timestamp, content)
	}
}
