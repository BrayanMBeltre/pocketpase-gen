package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// AppName is the name of the application.
const AppName = "pb-gen"

// Version is the version of the application (can be set at build time).
var Version = "0.0.1-dev" // Default version

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   AppName,
	Short: "A generator tool for PocketBase Go applications",
	Long: `pb-gen is a CLI tool designed to simplify Go development with PocketBase
by generating type-safe code based on your PocketBase collection schemas.

Currently, it supports generating a schema definition file containing
constants for collection and field names.`,
	Version: Version, // Enables the --version flag automatically
	// If the root command itself shouldn't do anything except show help
	// when called without subcommands, you can leave RunE empty or
	// have it print help explicitly. Cobra often handles this by default.
	// RunE: func(cmd *cobra.Command, args []string) error {
	//  return cmd.Help()
	// },

	// SilenceUsage will prevent the usage message from being printed on error,
	// useful if you manually handle error printing. Defaults to false.
	// SilenceUsage: true,

	// SilenceErrors will prevent Cobra from printing errors; useful if you want
	// full control over error output formatting in main.go. Defaults to false.
	// SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	// Add subcommands here.
	rootCmd.AddCommand(GenerateCmd)

	// Execute the root command. Cobra handles parsing args, flags, and
	// calling the appropriate subcommand RunE function.
	err := rootCmd.Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
	}

	return err // Return the error for main.go to handle exit codes
}

func init() {
	// cobra.OnInitialize(initConfig) // Example: If you had config file loading

	// Define global persistent flags for the root command here.
	// These flags will be available to all subcommands.
	// Example: Add the verbose flag globally if needed by multiple commands.
	// Note: We defined 'verbose' in generate.go and attached it there.
	// If you want it truly global, define the var here and attach it here:
	//
	// var globalVerbose bool // Needs to be accessible by subcommands if they use it
	// rootCmd.PersistentFlags().BoolVarP(&globalVerbose, "verbose", "v", false, "Enable verbose logging output")

	// You can bind flags to Viper config here if using Viper
	// viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Set version template if you want to customize the --version output
	// rootCmd.SetVersionTemplate(...)
}
