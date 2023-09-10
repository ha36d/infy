package cmd

import (
	"context"

	corepulumi "github.com/ha36d/infy/pkg/pulumi/core"
	"github.com/spf13/cobra"
)

// previewCmd represents the copy command
var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview the infrastructure",
	Long:  `Preview the infrastructure from a yaml file`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		corepulumi.Preview(ctx, C.Name, C.Team, C.Env, C.Cloud, C.Account, C.Region, C.Components)
	},
}

func init() {
	rootCmd.AddCommand(previewCmd)
}
