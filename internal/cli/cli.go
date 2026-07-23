package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	yesFlag bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vendor",
	Short: "vendor-cli — Vendor npm packages directly into your repo",
	Long:  "A CLI tool to vendor npm packages directly into a repo as owned source, while still tracking upstream versions.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Persistent flags
	rootCmd.PersistentFlags().BoolVarP(&yesFlag, "yes", "y", false, "Skip interactive confirmation prompts")

	// Add subcommands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(initCmd)

	// Subcommand specific flags
	statusCmd.Flags().Bool("check-updates", false, "Check registry for newer versions available")
}

var addCmd = &cobra.Command{
	Use:   "add <pkg>[@version]",
	Short: "Vendor a package from npm",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("TODO: Implement add for %s (yes=%v)\n", args[0], yesFlag)
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove <pkg>",
	Short: "Remove a vendored package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("TODO: Implement remove for %s\n", args[0])
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [pkg]",
	Short: "Update a vendored package (or all if omitted)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg := "all"
		if len(args) > 0 {
			pkg = args[0]
		}
		fmt.Printf("TODO: Implement update for %s\n", pkg)
	},
}

var diffCmd = &cobra.Command{
	Use:   "diff <pkg>",
	Short: "Show local modifications vs. baseline",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("TODO: Implement diff for %s\n", args[0])
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List vendored packages",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		checkUpdates, _ := cmd.Flags().GetBool("check-updates")
		fmt.Printf("TODO: Implement status (check-updates=%v)\n", checkUpdates)
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Re-apply workspace registration and install",
	Aliases: []string{"install"},
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: Implement sync")
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "First-time setup",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: Implement init")
	},
}
