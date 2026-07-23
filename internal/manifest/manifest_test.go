package manifest

import (
	"os"
	"testing"
)

func TestRecordAndRemoveManifest(t *testing.T) {
	// Isolate test in a temp directory
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	err := RecordManifest("is-even", "1.0.0", "http://example.com/tar.tgz", "vendor/is-even")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, err := GetPackageConfig("is-even")
	if err != nil || cfg.Version != "1.0.0" {
		t.Fatalf("expected 1.0.0, got %v", cfg)
	}

	err = RemoveManifest("is-even")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = GetPackageConfig("is-even")
	if err == nil {
		t.Fatalf("expected error after removing manifest")
	}
}

func TestMarkPatched(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	RecordManifest("testpkg", "2.0.0", "url", "dir")
	if err := MarkPatched("testpkg", "patches/testpkg.patch"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, _ := GetPackageConfig("testpkg")
	if !cfg.Patched {
		t.Fatalf("expected patched to be true")
	}

	lock, _ := GetPackageLock("testpkg")
	if lock.PatchFile != "patches/testpkg.patch" {
		t.Fatalf("expected patchFile to be set, got %s", lock.PatchFile)
	}
}
