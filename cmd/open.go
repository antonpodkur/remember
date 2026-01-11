package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/antonpodkur/remember/internal/storage"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [name]",
	Short: "Open a note or the storage folder in your editor",
	Long:  `Open a specific note in $EDITOR, or open the ~/.remember/ folder if no name is provided`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runOpen,
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
	rootCmd.AddCommand(openCmd)
}

func runOpen(cmd *cobra.Command, args []string) {
	editor := resolveEditor()
	if editor == "" {
		fmt.Fprintln(os.Stderr, "Error: no editor found, set $EDITOR")
		os.Exit(1)
	}

	var path string
	if len(args) == 0 {
		path = storage.GetStorageDir()
		if err := storage.EnsureStorageDir(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
	} else {
		name := args[0]

		// Handle "latest" as a special name
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
		path = storage.GetNotePath(name)
	}

	execCmd := exec.Command(editor, path)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func resolveEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	if _, err := exec.LookPath("vim"); err == nil {
		return "vim"
	}

	if _, err := exec.LookPath("nano"); err == nil {
		return "nano"
	}

	return ""
}
