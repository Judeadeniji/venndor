package pm

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
)

// Detect attempts to detect the package manager by looking at packageManager field in package.json
// or falling back to lockfiles. It returns the package manager name and a boolean indicating
// whether the "packageManager" field was explicitly set (which means we should use corepack).
func Detect() (string, bool, error) {
	// 1. Check package.json for packageManager (corepack style)
	b, err := os.ReadFile("package.json")
	if err == nil {
		var pkg struct {
			PackageManager string `json:"packageManager"`
		}
		if err := json.Unmarshal(b, &pkg); err == nil && pkg.PackageManager != "" {
			if strings.HasPrefix(pkg.PackageManager, "pnpm") {
				return "pnpm", true, nil
			}
			if strings.HasPrefix(pkg.PackageManager, "yarn") {
				return "yarn", true, nil
			}
			if strings.HasPrefix(pkg.PackageManager, "npm") {
				return "npm", true, nil
			}
			if strings.HasPrefix(pkg.PackageManager, "bun") {
				return "bun", true, nil
			}
		}
	}

	// 2. Fall back to lockfiles (no corepack)
	if _, err := os.Stat("pnpm-lock.yaml"); err == nil {
		return "pnpm", false, nil
	}
	if _, err := os.Stat("yarn.lock"); err == nil {
		return "yarn", false, nil
	}
	if _, err := os.Stat("package-lock.json"); err == nil {
		return "npm", false, nil
	}
	if _, err := os.Stat("bun.lockb"); err == nil {
		return "bun", false, nil
	}

	// Default to npm if nothing else is found
	return "npm", false, nil
}

// Install runs the install command for the specified package manager.
func Install(manager string, useCorepack bool) error {
	args := []string{"install"}
	bin := manager

	if useCorepack && manager != "bun" { // bun doesn't use corepack
		args = append([]string{manager}, args...)
		bin = "corepack"
	}

	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
