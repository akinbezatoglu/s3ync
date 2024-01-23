package destroy

import (
	"fmt"

	"github.com/akinbezatoglu/s3ync/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdDestroy(cfg config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "destroy",
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("s3ync destroy")
		},
	}

	return cmd
}
