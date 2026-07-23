package cli

import (
	"fmt"
	"os"
)

// Execute runs the CLI application
func Execute() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add":
		handleAdd(os.Args[2:])
	case "remove":
		handleRemove(os.Args[2:])
	case "update":
		handleUpdate(os.Args[2:])
	case "diff":
		handleDiff(os.Args[2:])
	case "status":
		handleStatus(os.Args[2:])
	case "sync", "install":
		handleSync(os.Args[2:])
	case "init":
		handleInit(os.Args[2:])
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`vendor-cli — Vendor npm packages directly into your repo

Usage:
  vendor <command> [arguments]

Commands:
  add <pkg>[@version]  Vendor a package from npm
  remove <pkg>         Remove a vendored package
  update [pkg]         Update a vendored package (or all if omitted)
  diff <pkg>           Show local modifications vs. baseline
  status               List vendored packages
  sync                 Re-apply workspace registration and install
  init                 First-time setup`)
}

func handleAdd(args []string) {
	fmt.Println("TODO: Implement add")
}

func handleRemove(args []string) {
	fmt.Println("TODO: Implement remove")
}

func handleUpdate(args []string) {
	fmt.Println("TODO: Implement update")
}

func handleDiff(args []string) {
	fmt.Println("TODO: Implement diff")
}

func handleStatus(args []string) {
	fmt.Println("TODO: Implement status")
}

func handleSync(args []string) {
	fmt.Println("TODO: Implement sync")
}

func handleInit(args []string) {
	fmt.Println("TODO: Implement init")
}
