package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for remember.

To load completions:

Bash:
  $ source <(remember completion bash)
  # To load completions for each session, execute once:
  $ remember completion bash >> ~/.bashrc

Zsh:
  $ source <(remember completion zsh)
  # To load completions for each session, execute once:
  $ remember completion zsh >> ~/.zshrc

Fish:
  $ remember completion fish | source
  # To load completions for each session, execute once:
  $ remember completion fish > ~/.config/fish/completions/remember.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			rootCmd.GenFishCompletion(os.Stdout, true)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
