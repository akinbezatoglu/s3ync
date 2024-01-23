package list

import (
	"fmt"
	"os"

	"github.com/akinbezatoglu/s3ync/internal/config"
	"github.com/spf13/cobra"
)

func NewListCmd(cfg config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:  "list",
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			profiles, err := cfg.GetProfileNames()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, profile := range profiles {
				syncs, err := cfg.GetSyncListFromProfile(profile)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if syncs == nil {
					fmt.Printf("Profile: (%v), there is no sync data.\n", profile)
				} else {
					fmt.Printf("Profile: (%v)\n", profile)
					for i := 0; i < len(syncs)/2; i++ {
						fmt.Printf("- local: %v\n", syncs[i*2])
						fmt.Printf("  bucket: %v\n", syncs[i*2+1])
					}
				}
			}
		},
	}
	return cmd
}
