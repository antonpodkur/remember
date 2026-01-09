package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/antonpodkur/remember/internal/storage"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "remember <name> [content]",
	Short: "Zero-friction note-taking for developers",
	Long:  `remember is a minimal CLI for capturing notes. Append timestamped text to markdown files stored in ~/.remember/`,
	Args:  cobra.MinimumNArgs(0),
	Run:   runRoot,
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
		content = readInteractive()
		if strings.TrimSpace(content) == "" {
			fmt.Fprintln(os.Stderr, "Error: nothing to save")
			os.Exit(1)
		}
	}

	if err := storage.AppendToNote(name, content); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Added to %s\n", name)
}

func readInteractive() string {
	if fileInfo, _ := os.Stdin.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		fmt.Println("(Ctrl+D to save)")
	}

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return strings.Join(lines, "\n")
}
