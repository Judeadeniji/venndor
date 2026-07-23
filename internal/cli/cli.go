package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/judeadeniji/venndor/internal/manifest"
	"github.com/judeadeniji/venndor/internal/npm"
	"github.com/judeadeniji/venndor/internal/pm"
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
		pkgArg := args[0]
		
		// Parse pkg and version
		pkgName := pkgArg
		version := ""
		
		// Handle scoped packages like @org/pkg@1.0.0
		idx := strings.LastIndex(pkgArg, "@")
		if idx > 0 { // > 0 to skip the first character in case it's a scoped package
			pkgName = pkgArg[:idx]
			version = pkgArg[idx+1:]
		}

		destDir := filepath.Join("vendor", pkgName)

		fmt.Printf("Vendoring %s into %s...\n", pkgArg, destDir)
		targetVersion, tarballURL, err := npm.FetchAndExtract(pkgName, version, destDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Configuring workspace and imports...\n")
		pkgJSONPath := "package.json"
		
		if err := manifest.EnsureWorkspace(pkgJSONPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to add to workspaces: %v\n", err)
		}
		
		if err := manifest.EnsureImport(pkgJSONPath, pkgName); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to add import alias: %v\n", err)
		}

		fmt.Printf("Updating vendor.yaml and vendor-lock.json...\n")
		if err := manifest.RecordManifest(pkgName, targetVersion, tarballURL, destDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write manifest: %v\n", err)
		}
		
		fmt.Printf("Running install...\n")
		manager, useCorepack, err := pm.Detect()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to detect package manager: %v\n", err)
		} else {
			if err := pm.Install(manager, useCorepack); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: install failed: %v\n", err)
			}
		}

		fmt.Printf("Successfully vendored %s@%s\n", pkgName, targetVersion)
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
