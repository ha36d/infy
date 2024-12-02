package cmd

import (
	"context"

	corepulumi "github.com/ha36d/infy/pkg/pulumi/core"
	"github.com/spf13/cobra"
)

// buildCmd represents the copy command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the infrastructure",
	Long:  `Build the infrastructure from a yaml file`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		corepulumi.Up(ctx, C.Metadata, C.Cloud, C.Account, C.Region, C.Components)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
