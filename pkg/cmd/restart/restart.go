package restart

import (
	"fmt"

	"github.com/akinbezatoglu/s3ync/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdRestart(cfg config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "restart",
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("s3ync restart")
		},
	}
	return cmd
}
