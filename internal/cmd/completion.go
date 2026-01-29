package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewCompletionCmd(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for Mandor.

To load completions:

Bash:

  $ source <(mandor completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ mandor completion bash > /etc/bash_completion.d/mandor
  # macOS:
  $ mandor completion bash > /usr/local/etc/bash_completion.d/mandor

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ mandor completion zsh > "${fpath[1]}/_mandor"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ mandor completion fish | source

  # To load completions for each session, execute once:
  $ mandor completion fish > ~/.config/fish/completions/mandor.fish
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish"},
		Args:                  cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				return rootCmd.GenZshCompletion(os.Stdout)
			case "fish":
				return rootCmd.GenFishCompletion(os.Stdout, true)
			}
			return nil
		},
	}

	return cmd
}
