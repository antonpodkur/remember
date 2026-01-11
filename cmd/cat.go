package cmd

import (
	"fmt"
	"os"

	"github.com/antonpodkur/remember/internal/storage"
	"github.com/spf13/cobra"
)

var catCmd = &cobra.Command{
	Use:   "cat <name|latest>",
	Short: "Print note contents to terminal",
	Long:  `Print the raw contents of a note to stdout`,
	Args:  cobra.ExactArgs(1),
	Run:   runCat,
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
	rootCmd.AddCommand(catCmd)
}

func runCat(cmd *cobra.Command, args []string) {
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

	content, err := storage.ReadNoteContent(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Print(content)
}
