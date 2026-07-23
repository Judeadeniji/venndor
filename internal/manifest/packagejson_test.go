package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureWorkspace(t *testing.T) {
	dir := t.TempDir()
	pkgJSON := filepath.Join(dir, "package.json")
	
	// Test on non-existent file
	err := EnsureWorkspace(pkgJSON)
	if err == nil {
		t.Fatal("expected error when package.json does not exist")
	}

	// Create dummy package.json
	os.WriteFile(pkgJSON, []byte(`{"name":"test"}`), 0644)
	
	// Ensure workspace adds "vendor/*"
	if err := EnsureWorkspace(pkgJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, _ := os.ReadFile(pkgJSON)
	if !strings.Contains(string(b), `"workspaces": [`) || !strings.Contains(string(b), `"vendor/*"`) {
		t.Fatalf("expected vendor/* in workspaces, got %s", string(b))
	}
}

func TestEnsureAndRemoveImport(t *testing.T) {
	dir := t.TempDir()
	pkgJSON := filepath.Join(dir, "package.json")
	os.WriteFile(pkgJSON, []byte(`{"name":"test"}`), 0644)

	// Ensure import
	if err := EnsureImport(pkgJSON, "is-odd"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, _ := os.ReadFile(pkgJSON)
	if !strings.Contains(string(b), `"#vendor/is-odd": "./vendor/is-odd"`) {
		t.Fatalf("expected import map entry, got %s", string(b))
	}

	// Remove import
	if err := RemoveImport(pkgJSON, "is-odd"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b2, _ := os.ReadFile(pkgJSON)
	if strings.Contains(string(b2), `"#vendor/is-odd"`) {
		t.Fatalf("expected import map entry to be removed, got %s", string(b2))
	}
}
