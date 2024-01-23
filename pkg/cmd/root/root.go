package root

import (
	"github.com/akinbezatoglu/s3ync/internal/config"
	configCmd "github.com/akinbezatoglu/s3ync/pkg/cmd/config"
	destroyCmd "github.com/akinbezatoglu/s3ync/pkg/cmd/destroy"
	restartCmd "github.com/akinbezatoglu/s3ync/pkg/cmd/restart"
	statusCmd "github.com/akinbezatoglu/s3ync/pkg/cmd/status"
	stopCmd "github.com/akinbezatoglu/s3ync/pkg/cmd/stop"
	syncCmd "github.com/akinbezatoglu/s3ync/pkg/cmd/sync"
	unsyncCmd "github.com/akinbezatoglu/s3ync/pkg/cmd/unsync"
	versionCmd "github.com/akinbezatoglu/s3ync/pkg/cmd/version"
	"github.com/spf13/cobra"
)

func NewCmdRoot(cfg config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:         "s3ync",
		Short:       "Synchronize your local folders with buckets",
		Long:        `Automate the synchronization of your logs, local backups or changes with buckets seamlessly`,
		Annotations: map[string]string{},
	}

	cmd.AddCommand(configCmd.NewCmdConfig(cfg))
	cmd.AddCommand(statusCmd.NewCmdStatus(cfg))
	cmd.AddCommand(syncCmd.NewCmdSync(cfg))
	cmd.AddCommand(unsyncCmd.NewCmdUnsync(cfg))
	cmd.AddCommand(versionCmd.NewCmdVersion(cfg))
	cmd.AddCommand(restartCmd.NewCmdRestart(cfg))
	cmd.AddCommand(stopCmd.NewCmdStop(cfg))
	cmd.AddCommand(destroyCmd.NewCmdDestroy(cfg))

	return cmd
}
