package unsync

import (
	"fmt"

	"github.com/akinbezatoglu/s3ync/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdUnsync(cfg config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "unsync",
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("s3ync unsync")
		},
	}
	return cmd
}
