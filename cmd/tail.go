package cmd

import (
	"fmt"
	"os"

	"github.com/antonpodkur/remember/internal/storage"
	"github.com/antonpodkur/remember/internal/ui"
	"github.com/spf13/cobra"
)

var tailCount int

var tailCmd = &cobra.Command{
	Use:   "tail <name|latest>",
	Short: "Show last N entries of a note",
	Long:  `Show the last N entries of a note (default: 3)`,
	Args:  cobra.ExactArgs(1),
	Run:   runTail,
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
	tailCmd.Flags().IntVarP(&tailCount, "count", "n", 3, "Number of entries to show")
	rootCmd.AddCommand(tailCmd)
}

func runTail(cmd *cobra.Command, args []string) {
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

	if len(entries) == 0 {
		fmt.Fprintf(os.Stderr, "Note '%s' has no entries\n", name)
		os.Exit(0)
	}

	// Get last N entries
	start := len(entries) - tailCount
	if start < 0 {
		start = 0
	}

	fmt.Printf("%s (last %d entries)\n\n", ui.NoteName.Render(name), len(entries[start:]))

	for i, entry := range entries[start:] {
		fmt.Printf("## %s\n\n%s\n", entry.Timestamp, entry.Content)
		if i < len(entries[start:])-1 {
			fmt.Println()
		}
	}
}
