package npm

import (
	"os"
	"testing"
)

func TestFetchMetadata(t *testing.T) {
	version, url, err := FetchMetadata("is-odd", "3.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "3.0.0" {
		t.Fatalf("expected 3.0.0, got %s", version)
	}
	if url == "" {
		t.Fatalf("expected tarball URL, got empty")
	}
}

func TestFetchAndExtract(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	version, url, err := FetchAndExtract("is-odd", "3.0.0", "vendor/is-odd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "3.0.0" {
		t.Fatalf("expected 3.0.0, got %s", version)
	}
	if url == "" {
		t.Fatalf("expected tarball URL, got empty")
	}

	// Verify extracted file
	if _, err := os.Stat("vendor/is-odd/index.js"); os.IsNotExist(err) {
		t.Fatalf("expected index.js to be extracted")
	}

	// Verify cache file
	cachePath := GetCachePath("is-odd", "3.0.0")
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Fatalf("expected cache file at %s", cachePath)
	}
}
