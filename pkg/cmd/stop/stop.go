package stop

import (
	"fmt"

	"github.com/akinbezatoglu/s3ync/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdStop(cfg config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "stop",
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("s3ync stop")
		},
	}
	return cmd
}
