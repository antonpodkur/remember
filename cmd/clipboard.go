package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/antonpodkur/remember/internal/storage"
	"github.com/antonpodkur/remember/internal/ui"
	"github.com/spf13/cobra"
)

var clipboardCmd = &cobra.Command{
	Use:   "clipboard <name|latest>",
	Short: "Append clipboard content to a note",
	Long:  `Read content from system clipboard and append it to the specified note`,
	Args:  cobra.ExactArgs(1),
	Run:   runClipboard,
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
	rootCmd.AddCommand(clipboardCmd)
}

func runClipboard(cmd *cobra.Command, args []string) {
	name := args[0]

	if name == "latest" {
		latestName, err := storage.GetLatestNote()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
		name = latestName
	} else {
		if storage.IsReservedName(name) {
			fmt.Fprintf(os.Stderr, "Error: %s is a reserved command\n", name)
			os.Exit(1)
		}

		if err := storage.ValidateName(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
	}

	content, err := readClipboard()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(content) == "" {
		fmt.Fprintln(os.Stderr, "Error: clipboard is empty")
		os.Exit(1)
	}

	if err := storage.AppendToNote(name, content); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s %s\n", ui.Success.Render("✓ Added clipboard to"), ui.NoteName.Render(name))
}

func readClipboard() (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbpaste")
	case "linux":
		// Try xclip first, then xsel
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--output")
		} else {
			return "", fmt.Errorf("clipboard tool not found (install xclip or xsel)")
		}
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to read clipboard: %w", err)
	}

	return string(output), nil
}
