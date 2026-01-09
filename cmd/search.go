package cmd

import (
	"fmt"
	"os"

	"github.com/antonpodkur/remember/internal/storage"
	"github.com/antonpodkur/remember/internal/ui"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Fuzzy search note names",
	Long:  `Search for notes by name using fuzzy matching`,
	Args:  cobra.ExactArgs(1),
	Run:   runSearch,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		notes, _ := storage.ListNotes()
		return notes, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) {
	query := args[0]

	notes, err := storage.ListNotes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if len(notes) == 0 {
		return
	}

	matches := fuzzy.Find(query, notes)
	for _, match := range matches {
		fmt.Println(ui.NoteName.Render(match.Str))
	}
}
