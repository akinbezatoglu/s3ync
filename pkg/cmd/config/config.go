package config

import (
	"fmt"
	"os"

	"github.com/akinbezatoglu/s3ync/internal/config"
	"github.com/akinbezatoglu/s3ync/internal/sync/s3"
	addCmd "github.com/akinbezatoglu/s3ync/pkg/cmd/config/add"
	initCmd "github.com/akinbezatoglu/s3ync/pkg/cmd/config/init"
	listCmd "github.com/akinbezatoglu/s3ync/pkg/cmd/config/list"
	"github.com/spf13/cobra"
)

func NewCmdConfig(cfg config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "config",
		Short: "Manage configuration for s3ync",
		Args:  cobra.ExactArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// ~/.aws/config file does not exist or never configured in aws cli.
			if !s3.IsAwsCliConfigured() {
				fmt.Println("To get started with s3ync, you need to install and configure aws-cli (https://github.com/aws/aws-cli)")
				os.Exit(1)
			}
		},
	}

	cmd.AddCommand(initCmd.NewInitCmd(cfg))
	cmd.AddCommand(listCmd.NewListCmd(cfg))
	cmd.AddCommand(addCmd.NewCmdAdd(cfg))

	return cmd
}
