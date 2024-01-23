package list

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdList() *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "list",
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("s3ync sync list")
		},
	}
	return cmd
}
