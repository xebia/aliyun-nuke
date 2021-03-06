package cmd

import (
	"fmt"
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/nuker"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var regions []string
var excludedIds []string
var force bool

func init() {
	defaultRegions := make([]string, 0)
	for _, region := range account.AllRegions {
		defaultRegions = append(defaultRegions, string(region))
	}

	destroyCmd.Flags().StringSliceVarP(&regions, "regions", "r", defaultRegions, "Specify list of regions to destroy resources in")
	destroyCmd.Flags().StringSliceVarP(&excludedIds, "exclude", "e", []string{}, "List of resource IDs to exclude from deletion")
	destroyCmd.Flags().BoolVarP(&force, "force", "f", false, "Force deletion of all resources, including system resources")
	rootCmd.AddCommand(destroyCmd)
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys all resources in your account",
	Run: func(cmd *cobra.Command, args []string) {
		accessKeyId, ok := os.LookupEnv("ALIYUN_NUKE_ACCESS_KEY_ID")
		accessKeySecret, ok := os.LookupEnv("ALIYUN_NUKE_ACCESS_KEY_SECRET")

		if !ok {
			log.Fatal("credential error: ALIYUN_NUKE_ACCESS_KEY_ID and ALIYUN_NUKE_ACCESS_KEY_SECRET undefined")
		}

		currentAccount := account.Account{
			Credentials: account.Credentials{
				AccessKeyID:     accessKeyId,
				AccessKeySecret: accessKeySecret,
			},
		}

		var regionsToDestroy []account.Region
		for _, region := range regions {
			regionsToDestroy = append(regionsToDestroy, account.Region(region))
		}

		log.Println(fmt.Sprintf("Starting destruction in regions: %s", regionsToDestroy))
		results := nuker.NukeItAll(currentAccount, regionsToDestroy, excludedIds, force)
		for result := range results {
			if result.Success {
				log.Println(fmt.Sprintf("Nuked: %s - %s", result.Resource.Type(), result.Resource.Id()))
			} else if result.Skipped {
				log.Println(fmt.Sprintf("Skipped: %s - %s", result.Resource.Type(), result.Resource.Id()))
			} else {
				log.Println(result.Error)
			}
		}
		log.Println("Account has converged or max. retries exhausted")
	},
}
