package init

import (
	"github.com/akinbezatoglu/s3ync/internal/config"
	"github.com/spf13/cobra"
)

func NewInitCmd(cfg config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "init",
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			//	profiles := s3.GetLocalAwsProfilesWithDefaultRegion()
			//	// add profiles to config file with the region values
			//	for i := 0; i < len(profiles)/2; i++ {
			//		cfg.Set([]string{"s3", "profiles", profiles[i*2]}, "")
			//		cfg.Set([]string{"s3", "profiles", profiles[i*2], "region"}, profiles[i*2+1])
			//		cfg.Set([]string{"s3", "profiles", profiles[i*2], "syncs"}, "")
			//	}
			//	// write config file
			//	if err := cfg.Write(); err != nil {
			//		fmt.Println(err)
			//		os.Exit(1)
			//	}
			//	fmt.Println("Succesfully initialize the config file:", cfg.GetConfigFilePath())
		},
	}
	return cmd
}
