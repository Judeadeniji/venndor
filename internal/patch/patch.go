package patch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Generate creates a patch file comparing the pristine directory with the vendored directory.
// Returns true if differences were found and a patch was created, false otherwise.
func Generate(pkgName, pristineDir, vendoredDir string) (bool, error) {
	if err := os.MkdirAll("patches", 0755); err != nil {
		return false, err
	}

	safeName := strings.ReplaceAll(pkgName, "/", "-")
	patchFile := filepath.Join("patches", safeName+".patch")

	// diff -urN pristine vendored > patchFile
	cmd := exec.Command("diff", "-urN", pristineDir, vendoredDir)
	
	out, err := os.Create(patchFile)
	if err != nil {
		return false, err
	}
	defer out.Close()

	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	
	// diff exits with 1 if differences were found
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() == 1 {
			return true, nil // Differences found, patch created successfully
		}
		return false, fmt.Errorf("diff command failed: %v", err)
	}

	if err != nil {
		return false, fmt.Errorf("failed to run diff: %v", err)
	}

	// No differences found
	os.Remove(patchFile) // Clean up empty patch file
	return false, nil
}

// Apply applies a patch file to a target directory.
func Apply(patchFile, targetDir string) error {
	absPatch, err := filepath.Abs(patchFile)
	if err != nil {
		return err
	}

	cmd := exec.Command("patch", "-p2", "-i", absPatch)
	cmd.Dir = targetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
