package cache

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// StashBaseline copies the vendored package directory into the .vendor-cache.
// srcDir: e.g. "vendor/is-odd"
// pkgName: e.g. "is-odd"
// version: e.g. "3.0.1"
func StashBaseline(srcDir, pkgName, version string) error {
	cacheDir := filepath.Join(".vendor-cache", fmt.Sprintf("%s-%s", pkgName, version))
	
	// Remove cache dir if it exists to ensure a clean stash
	if err := os.RemoveAll(cacheDir); err != nil {
		return err
	}

	return copyDir(srcDir, cacheDir)
}

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}
