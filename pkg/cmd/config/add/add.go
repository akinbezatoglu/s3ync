package add

import (
	"fmt"

	"github.com/akinbezatoglu/s3ync/internal/config"
	"github.com/akinbezatoglu/s3ync/internal/sync/s3"
	"github.com/spf13/cobra"
)

func NewCmdAdd(cfg config.Config) *cobra.Command {
	var newProfile string
	var cmd = &cobra.Command{
		Use:  "add",
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			isAwsCfg := s3.IsProfileExistInAwsConfig(newProfile)
			isInCfgFile := cfg.IsProfileExistInConfigFile(newProfile)
			if isAwsCfg && !isInCfgFile {
				p_with_region := s3.GetLocalAwsProfilesWithDefaultRegion()
				for i := 0; i < len(p_with_region)/2; i++ {
					if p_with_region[i*2] == newProfile {
						cfg.Set([]string{"s3", "profiles", newProfile}, "")
						cfg.Set([]string{"s3", "profiles", newProfile, "region"}, p_with_region[i*2+1])
						cfg.Set([]string{"s3", "profiles", newProfile, "syncs"}, "")
						cfg.Write()
					}
				}
				fmt.Printf("Succesfully, %v added to config file", newProfile)
			} else if isAwsCfg && isInCfgFile {
				fmt.Printf("%v is already configured", newProfile)
			} else {
				fmt.Printf("%v is not configured in aws cli. Please configure the %v in aws cli first.", newProfile, newProfile)
			}
		},
	}

	cmd.PersistentFlags().StringVar(&newProfile, "profile", "", "Add new profile to config file")

	return cmd
}
