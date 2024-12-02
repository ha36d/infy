package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
)

var cfgFile string
var verbose bool
var C config

type config struct {
	Components []map[string]any  `mapstructure:"components"`
	Metadata   map[string]string `mapstructure:"metadata"`
	Cloud      string            `mapstructure:"cloud"`
	Account    string            `mapstructure:"account"`
	Region     string            `mapstructure:"region"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "infy",
	Short: "abstraction on top of pulumi",
	Long:  `abstraction on top of pulumi`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.infy.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbosity")
	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		log.Printf("Failed to bind flag: %v\n", err)
	}
	viper.SetDefault("verbose", false)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		// Search config in working directory, and then home directory with name ".infy" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".infy")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No config file found! Counting on flags!")
	}

	// Load JSON schema from file
	schemaFile, err := os.ReadFile("config_schema.json")
	if err != nil {
		log.Fatalf("Error reading schema file: %v", err)
	}
	schemaLoader := gojsonschema.NewStringLoader(string(schemaFile))
	documentLoader := gojsonschema.NewGoLoader(viper.AllSettings())
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		log.Fatalf("Error validating config: %v", err)
	}

	if !result.Valid() {
		log.Println("The configuration is not valid:")
		for _, desc := range result.Errors() {
			log.Printf("- %s\n", desc)
		}
		os.Exit(1)
	}

	// Unmarshal the validated configuration
	if err := viper.Unmarshal(&C); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	verbose = viper.GetBool("verbose")
	if verbose {
		// If a config file is found, read it in.
		fmt.Fprintln(os.Stderr, "Using config:", viper.ConfigFileUsed())
		log.Println("--- Configuration ---")
		for s, i := range viper.AllSettings() {
			log.Printf("\t%s = %s\n", s, i)
		}
		log.Println("---")
	}
}
