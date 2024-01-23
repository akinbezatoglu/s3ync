package sync

import (
	"fmt"

	"github.com/akinbezatoglu/s3ync/internal/config"
	syncListCmd "github.com/akinbezatoglu/s3ync/pkg/cmd/sync/list"
	"github.com/spf13/cobra"
)

func NewCmdSync(cfg config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "sync",
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("s3ync sync")
		},
	}

	cmd.AddCommand(syncListCmd.NewCmdList())
	return cmd
}
