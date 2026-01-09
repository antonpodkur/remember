package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/antonpodkur/remember/internal/storage"
	"github.com/antonpodkur/remember/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "remember <name> [content]",
	Short: "Zero-friction note-taking for developers",
	Long:  `remember is a minimal CLI for capturing notes. Append timestamped text to markdown files stored in ~/.remember/`,
	Args:  cobra.MinimumNArgs(0),
	Run:   runRoot,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		notes, err := storage.ListNotes()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return notes, cobra.ShellCompDirectiveNoFileComp
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func runRoot(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	name := args[0]

	if storage.IsReservedName(name) {
		fmt.Fprintf(os.Stderr, "Error: %s is a reserved command\n", name)
		os.Exit(1)
	}

	if err := storage.ValidateName(name); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	var content string
	if len(args) > 1 {
		content = strings.Join(args[1:], " ")
		if strings.TrimSpace(content) == "" {
			fmt.Fprintln(os.Stderr, "Error: content cannot be empty")
			os.Exit(1)
		}
	} else {
		// Check if stdin is a terminal (interactive mode) or a pipe
		if fileInfo, _ := os.Stdin.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
			// Interactive terminal - use styled input
			var ok bool
			content, ok = ui.RunInteractiveInput(name)
			if !ok {
				fmt.Fprintln(os.Stderr, "Cancelled")
				os.Exit(0)
			}
		} else {
			// Piped input - read from stdin
			content = readFromPipe()
		}
		if strings.TrimSpace(content) == "" {
			fmt.Fprintln(os.Stderr, "Error: nothing to save")
			os.Exit(1)
		}
	}

	if err := storage.AppendToNote(name, content); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s %s\n", ui.Success.Render("✓ Added to"), ui.NoteName.Render(name))
}

func readFromPipe() string {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return strings.Join(lines, "\n")
}
