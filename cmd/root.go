package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     "aliyun-nuke",
	Short:   "aliyun-nuke removes all resources in your Alibaba Cloud account",
	Version: "0.1.3",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
