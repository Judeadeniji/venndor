package cli_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_EndToEnd(t *testing.T) {
	// Build the CLI binary first
	binName := "venndor-cli"
	if os.Getenv("GOOS") == "windows" {
		binName += ".exe"
	}
	binPath := filepath.Join(t.TempDir(), binName)
	
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/vendor")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build cli: %v", err)
	}

	// Create test environment
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Create dummy package.json
	os.WriteFile("package.json", []byte(`{"name":"test"}`), 0644)

	runCLI := func(args ...string) (string, error) {
		cmd := exec.Command(binPath, args...)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}

	// 1. vendor init
	out, err := runCLI("init")
	if err != nil {
		t.Fatalf("init failed: %v\nOutput: %s", err, out)
	}
	if _, err := os.Stat("vendor.yaml"); os.IsNotExist(err) {
		t.Fatal("init should create vendor.yaml")
	}

	// 2. vendor add is-even@1.0.0
	out, err = runCLI("add", "is-even@1.0.0")
	if err != nil {
		t.Fatalf("add failed: %v\nOutput: %s", err, out)
	}
	if _, err := os.Stat("vendor/is-even/index.js"); os.IsNotExist(err) {
		t.Fatal("expected is-even to be vendored")
	}

	// 3. vendor diff (no changes)
	out, err = runCLI("diff", "is-even")
	if err != nil {
		t.Fatalf("diff failed: %v\nOutput: %s", err, out)
	}
	if !strings.Contains(out, "clean") {
		t.Fatalf("expected diff to say clean, got: %s", out)
	}

	// Make a change to trigger a diff
	os.WriteFile("vendor/is-even/index.js", []byte("module.exports = function() { return true; };\n"), 0644)
	
	// 4. vendor diff (with changes)
	out, err = runCLI("diff", "is-even")
	if err != nil {
		t.Fatalf("diff after changes failed: %v\nOutput: %s", err, out)
	}
	if _, err := os.Stat("patches/is-even.patch"); os.IsNotExist(err) {
		t.Fatal("expected patch to be generated")
	}

	// 5. vendor status
	out, err = runCLI("status")
	if err != nil {
		t.Fatalf("status failed: %v\nOutput: %s", err, out)
	}
	if !strings.Contains(out, "[PATCHED]") {
		t.Fatalf("status should show PATCHED: %s", out)
	}

	// 6. vendor update is-even (will update to latest and reapply patch)
	out, err = runCLI("update", "is-even")
	if err != nil {
		t.Fatalf("update failed: %v\nOutput: %s", err, out)
	}

	// Verify patch was applied after update
	b, _ := os.ReadFile("vendor/is-even/index.js")
	if !strings.Contains(string(b), "module.exports = function() { return true; };") {
		t.Fatalf("patch was not preserved after update")
	}

	// 7. vendor sync
	out, err = runCLI("sync")
	if err != nil {
		t.Fatalf("sync failed: %v\nOutput: %s", err, out)
	}

	// 8. vendor remove
	out, err = runCLI("remove", "is-even")
	if err != nil {
		t.Fatalf("remove failed: %v\nOutput: %s", err, out)
	}
	if _, err := os.Stat("vendor/is-even"); !os.IsNotExist(err) {
		t.Fatal("remove should delete vendor/is-even")
	}
	if _, err := os.Stat("patches/is-even.patch"); !os.IsNotExist(err) {
		t.Fatal("remove should delete patches")
	}
}
