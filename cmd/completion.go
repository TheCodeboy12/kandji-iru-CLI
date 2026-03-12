package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate shell completion script for kandji-iru-cli.

To load completions:

Bash:
  source <(kandji-iru-cli completion bash)
  # or add to ~/.bashrc:
  echo 'source <(kandji-iru-cli completion bash)' >> ~/.bashrc

Zsh:
  source <(kandji-iru-cli completion zsh)
  # or add to ~/.zshrc:
  echo 'source <(kandji-iru-cli completion zsh)' >> ~/.zshrc

Fish:
  kandji-iru-cli completion fish | source
  # or: kandji-iru-cli completion fish > ~/.config/fish/completions/kandji-iru-cli.fish

PowerShell:
  kandji-iru-cli completion powershell | Out-String | Invoke-Expression
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
