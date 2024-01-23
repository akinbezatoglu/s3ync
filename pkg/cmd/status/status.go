package status

import (
	"fmt"

	"github.com/akinbezatoglu/s3ync/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdStatus(cfg config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "status",
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("s3ync status")
		},
	}

	return cmd
}
