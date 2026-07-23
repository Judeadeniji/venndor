package npm

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// PackageMeta represents the metadata returned by the npm registry.
type PackageMeta struct {
	Name     string             `json:"name"`
	Versions map[string]Version `json:"versions"`
	DistTags map[string]string  `json:"dist-tags"`
}

// Version represents a specific version's metadata.
type Version struct {
	Version string `json:"version"`
	Dist    struct {
		Tarball   string `json:"tarball"`
		Integrity string `json:"integrity"`
	} `json:"dist"`
}

// FetchAndExtract resolves the package version, downloads the tarball, and extracts it to the destination directory.
// If version is empty, it uses the "latest" dist-tag.
func FetchAndExtract(pkgName, version, destDir string) error {
	// 1. Fetch metadata
	metaURL := fmt.Sprintf("https://registry.npmjs.org/%s", pkgName)
	resp, err := http.Get(metaURL)
	if err != nil {
		return fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registry returned status %s for package %s", resp.Status, pkgName)
	}

	var meta PackageMeta
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return fmt.Errorf("failed to decode metadata: %w", err)
	}

	// 2. Resolve version
	targetVersion := version
	if targetVersion == "" {
		latest, ok := meta.DistTags["latest"]
		if !ok {
			return fmt.Errorf("no 'latest' tag found for %s", pkgName)
		}
		targetVersion = latest
	}

	vData, ok := meta.Versions[targetVersion]
	if !ok {
		return fmt.Errorf("version %s not found for %s", targetVersion, pkgName)
	}

	tarballURL := vData.Dist.Tarball
	if tarballURL == "" {
		return fmt.Errorf("no tarball URL found for %s@%s", pkgName, targetVersion)
	}

	fmt.Printf("Fetching %s@%s...\n", pkgName, targetVersion)

	// 3. Download tarball
	tarResp, err := http.Get(tarballURL)
	if err != nil {
		return fmt.Errorf("failed to download tarball: %w", err)
	}
	defer tarResp.Body.Close()

	if tarResp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download tarball, status: %s", tarResp.Status)
	}

	// 4. Extract tarball
	return extractTarGz(tarResp.Body, destDir)
}

func extractTarGz(r io.Reader, destDir string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// npm tarballs usually contain a root "package/" directory. We want to strip it.
		// e.g. "package/index.js" -> "index.js"
		target := header.Name
		if strings.HasPrefix(target, "package/") {
			target = strings.TrimPrefix(target, "package/")
		}

		targetPath := filepath.Join(destDir, target)

		// Defend against Zip Slip
		if !strings.HasPrefix(targetPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", targetPath)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}

	return nil
}
