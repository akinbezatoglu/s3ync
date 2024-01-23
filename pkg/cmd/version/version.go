package version

import (
	"fmt"

	"github.com/akinbezatoglu/s3ync/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdVersion(cfg config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Root().Annotations["versionInfo"])
		},
	}
	return cmd
}
