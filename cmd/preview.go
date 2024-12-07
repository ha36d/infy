package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	corepulumi "github.com/ha36d/infy/pkg/pulumi/core"
)

var previewTimeout time.Duration

// previewCmd represents the infrastructure preview command
var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview the infrastructure",
	Long: `Preview the infrastructure changes from a yaml file.
This command shows what changes would be made without actually applying them.`,
	Example: `  infy preview
  infy preview --timeout 10m`,
	RunE: runPreview,
}

func init() {
	rootCmd.AddCommand(previewCmd)

	// Add preview-specific flags
	previewCmd.Flags().DurationVar(&previewTimeout, "timeout", 30*time.Minute, "timeout for the preview operation")
}

func runPreview(cmd *cobra.Command, args []string) error {
	// Create cancellable context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), previewTimeout)
	defer cancel()

	// Handle OS interrupts
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		log.Warn("Received interrupt signal, canceling preview...")
		cancel()
	}()

	// Validate configuration before proceeding
	if err := validatePreviewConfig(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Log preview parameters
	logPreviewParameters()

	// Execute the preview
	log.Info("Starting infrastructure preview...")
	startTime := time.Now()

	err := corepulumi.Preview(ctx, config.Metadata, config.Cloud, config.Account, config.Region, config.Components)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("preview operation timed out after %v", previewTimeout)
		}
		return fmt.Errorf("preview failed: %w", err)
	}

	// Log success and duration
	duration := time.Since(startTime)
	log.WithFields(log.Fields{
		"duration": duration.Round(time.Second),
		"cloud":    config.Cloud,
		"region":   config.Region,
	}).Info("Infrastructure preview completed successfully")

	return nil
}

func validatePreviewConfig() error {
	if config.Cloud == "" {
		return fmt.Errorf("cloud provider not specified")
	}
	if config.Account == "" {
		return fmt.Errorf("account not specified")
	}
	if config.Region == "" {
		return fmt.Errorf("region not specified")
	}
	if len(config.Components) == 0 {
		return fmt.Errorf("no components specified in configuration")
	}
	return nil
}

func logPreviewParameters() {
	log.WithFields(log.Fields{
		"cloud":      config.Cloud,
		"region":     config.Region,
		"components": len(config.Components),
		"timeout":    previewTimeout,
	}).Info("Preview parameters")
}
