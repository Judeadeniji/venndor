package patch

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateAndApply(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// Create pristine (two levels deep to simulate /tmp/.../pkg)
	pristineDir := filepath.Join("tmp", "pristine")
	os.MkdirAll(pristineDir, 0755)
	os.WriteFile(filepath.Join(pristineDir, "test.txt"), []byte("hello\n"), 0644)

	// Create vendored (two levels deep to simulate vendor/pkg)
	vendoredDir := filepath.Join("vendor", "vendored")
	os.MkdirAll(vendoredDir, 0755)
	os.WriteFile(filepath.Join(vendoredDir, "test.txt"), []byte("hello world\n"), 0644)

	found, err := Generate("testpkg", pristineDir, vendoredDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatalf("expected differences to be found")
	}

	// Check if patch file exists
	patchFile := filepath.Join("patches", "testpkg.patch")
	if _, err := os.Stat(patchFile); os.IsNotExist(err) {
		t.Fatalf("expected patch file to be generated")
	}

	// Create a new target dir matching pristine
	targetDir := "target"
	os.MkdirAll(targetDir, 0755)
	os.WriteFile(filepath.Join(targetDir, "test.txt"), []byte("hello\n"), 0644)

	// Apply patch
	if err := Apply(patchFile, targetDir); err != nil {
		t.Fatalf("failed to apply patch: %v", err)
	}

	// Verify applied patch
	b, _ := os.ReadFile(filepath.Join(targetDir, "test.txt"))
	if string(b) != "hello world\n" {
		t.Fatalf("expected 'hello world\\n', got %q", string(b))
	}
}
