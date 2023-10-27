package hooks

import (
	"fmt"
	"github.com/captainhook-go/captainhook/cli"
	"github.com/spf13/cobra"
)

func SetupHookPrepareCommitMsgCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pre-push",
		Short: "Execute pre-push actions",
		Long:  "Execute all actions configured for pre-push",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("PREPARE COMMIT MSG HOOK")
		},
	}

	cli.ConfigurationAware(cmd)
	cli.RepositoryAware(cmd)

	return cmd
}
