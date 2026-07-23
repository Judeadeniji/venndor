package pm

import (
	"os"
	"testing"
)

func TestDetect(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	// 1. No package.json -> defaults to npm, corepack=false
	manager, useCorepack, err := Detect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if manager != "npm" || useCorepack != false {
		t.Fatalf("expected npm/false, got %s/%v", manager, useCorepack)
	}

	// 2. packageManager field -> parsed correctly, corepack=true
	os.WriteFile("package.json", []byte(`{"packageManager": "yarn@3.2.1"}`), 0644)
	manager, useCorepack, err = Detect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if manager != "yarn" || useCorepack != true {
		t.Fatalf("expected yarn/true, got %s/%v", manager, useCorepack)
	}

	// 3. invalid or unversioned packageManager like bun
	os.WriteFile("package.json", []byte(`{"packageManager": "bun"}`), 0644)
	manager, useCorepack, err = Detect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if manager != "bun" || useCorepack != true {
		t.Fatalf("expected bun/true, got %s/%v", manager, useCorepack)
	}
}
