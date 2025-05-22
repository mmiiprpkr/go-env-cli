package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `To load completions:

Bash:

  $ source <(go-env-cli completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ go-env-cli completion bash > /etc/bash_completion.d/go-env-cli
  # macOS:
  $ go-env-cli completion bash > /usr/local/etc/bash_completion.d/go-env-cli

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ go-env-cli completion zsh > "${fpath[1]}/_go-env-cli"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ go-env-cli completion fish | source

  # To load completions for each session, execute once:
  $ go-env-cli completion fish > ~/.config/fish/completions/go-env-cli.fish

PowerShell:

  PS> go-env-cli completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> go-env-cli completion powershell > go-env-cli.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
