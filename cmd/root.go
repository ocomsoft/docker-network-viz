// Package cmd provides the CLI commands for the docker-network-viz tool.
// This file contains the root Cobra command and Viper configuration integration.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// AppName is the name of the application.
	AppName = "docker-network-viz"

	// AppDescription is a short description of the application.
	AppDescription = "Visualize Docker network topology in a tree-style format"
)

var (
	// cfgFile holds the path to the configuration file if specified.
	cfgFile string

	// noColor disables colored output.
	noColor bool

	// rootCmd is the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:   AppName,
		Short: AppDescription,
		Long: `docker-network-viz is a CLI tool for visualizing Docker network topology.

It provides tree-style output showing:
- Networks and their connected containers with aliases
- Container reachability across networks

This helps you understand which containers can communicate with each other
and through which networks.`,
		PersistentPreRunE: initConfig,
		SilenceUsage:      true,
		// Run visualize command by default when no subcommand is provided
		RunE: func(cmd *cobra.Command, args []string) error {
			// If a subcommand is being executed, don't run the default
			if cmd.Flags().Changed("help") {
				return cmd.Help()
			}
			// Execute the visualize command logic
			return runVisualize(cmd, args)
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

// init initializes the root command flags.
func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.docker-network-viz.yaml)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false,
		"disable colored output")

	// Flags for visualize command (also available on root for default behavior)
	rootCmd.Flags().StringVar(&onlyNetwork, "only-network", "",
		"show only the specified network")
	rootCmd.Flags().StringVar(&containerFilter, "container", "",
		"show only the specified container's connectivity")
	rootCmd.Flags().BoolVar(&noAliases, "no-aliases", false,
		"hide container aliases in the output")

	// Bind flags to viper
	_ = viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
	_ = viper.BindPFlag("only-network", rootCmd.Flags().Lookup("only-network"))
	_ = viper.BindPFlag("container", rootCmd.Flags().Lookup("container"))
	_ = viper.BindPFlag("no-aliases", rootCmd.Flags().Lookup("no-aliases"))
}

// initConfig reads in config file and ENV variables if set.
// This is called before any command runs via PersistentPreRunE.
func initConfig(_ *cobra.Command, _ []string) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to find home directory: %w", err)
		}

		// Search config in home directory with name ".docker-network-viz" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".docker-network-viz")
	}

	// Read in environment variables that match
	viper.SetEnvPrefix("DNV")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()

	return nil
}

// GetRootCmd returns the root command for testing purposes.
// This allows tests to access and configure the root command.
func GetRootCmd() *cobra.Command {
	return rootCmd
}

// ResetRootCmd resets the root command for testing purposes.
// This should be called between tests to ensure a clean state.
func ResetRootCmd() {
	rootCmd.ResetFlags()
	rootCmd.ResetCommands()

	// Re-initialize flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.docker-network-viz.yaml)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false,
		"disable colored output")
	rootCmd.Flags().StringVar(&onlyNetwork, "only-network", "",
		"show only the specified network")
	rootCmd.Flags().StringVar(&containerFilter, "container", "",
		"show only the specified container's connectivity")
	rootCmd.Flags().BoolVar(&noAliases, "no-aliases", false,
		"hide container aliases in the output")

	_ = viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
	_ = viper.BindPFlag("only-network", rootCmd.Flags().Lookup("only-network"))
	_ = viper.BindPFlag("container", rootCmd.Flags().Lookup("container"))
	_ = viper.BindPFlag("no-aliases", rootCmd.Flags().Lookup("no-aliases"))

	// Re-add subcommands
	rootCmd.AddCommand(visualizeCmd)
}
