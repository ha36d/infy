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

var (
	timeout     time.Duration
	forceUpdate bool
)

// buildCmd represents the infrastructure build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the infrastructure",
	Long: `Build the infrastructure from a yaml file.
This command will create or update infrastructure resources based on the configuration.`,
	Example: `  infy build
  infy build --timeout 30m
  infy build --force`,
	RunE: runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Add build-specific flags
	buildCmd.Flags().DurationVar(&timeout, "timeout", 2*time.Hour, "timeout for the build operation")
	buildCmd.Flags().BoolVar(&forceUpdate, "force", false, "force update even if no changes detected")
}

func runBuild(cmd *cobra.Command, args []string) error {
	// Create cancellable context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Handle OS interrupts
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		log.Warn("Received interrupt signal, initiating graceful shutdown...")
		cancel()
	}()

	// Validate configuration before proceeding
	if err := validateBuildConfig(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Log build parameters
	logBuildParameters()

	// Execute the build
	log.Info("Starting infrastructure build...")
	startTime := time.Now()

	err := corepulumi.Up(ctx, config.Metadata, config.Cloud, config.Account, config.Region, config.Components)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("build operation timed out after %v", timeout)
		}
		return fmt.Errorf("build failed: %w", err)
	}

	// Log success and duration
	duration := time.Since(startTime)
	log.WithFields(log.Fields{
		"duration": duration.Round(time.Second),
		"cloud":    config.Cloud,
		"region":   config.Region,
	}).Info("Infrastructure build completed successfully")

	return nil
}

func validateBuildConfig() error {
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

func logBuildParameters() {
	log.WithFields(log.Fields{
		"cloud":      config.Cloud,
		"region":     config.Region,
		"components": len(config.Components),
		"timeout":    timeout,
		"force":      forceUpdate,
	}).Info("Build parameters")
}
