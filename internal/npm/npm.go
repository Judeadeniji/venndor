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

// FetchMetadata retrieves package metadata and resolves the target version.
// It returns targetVersion, tarballURL, and error.
func FetchMetadata(pkgName, version string) (string, string, error) {
	metaURL := fmt.Sprintf("https://registry.npmjs.org/%s", pkgName)
	resp, err := http.Get(metaURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("registry returned status %s for package %s", resp.Status, pkgName)
	}

	var meta PackageMeta
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return "", "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	targetVersion := version
	if targetVersion == "" {
		latest, ok := meta.DistTags["latest"]
		if !ok {
			return "", "", fmt.Errorf("no 'latest' tag found for %s", pkgName)
		}
		targetVersion = latest
	}

	vData, ok := meta.Versions[targetVersion]
	if !ok {
		return "", "", fmt.Errorf("version %s not found for %s", targetVersion, pkgName)
	}

	tarballURL := vData.Dist.Tarball
	if tarballURL == "" {
		return "", "", fmt.Errorf("no tarball URL found for %s@%s", pkgName, targetVersion)
	}

	return targetVersion, tarballURL, nil
}

// FetchAndExtract resolves the package version, downloads the tarball, and extracts it to the destination directory.
// Returns the resolved version and the tarball URL.
func FetchAndExtract(pkgName, version, destDir string) (string, string, error) {
	targetVersion, tarballURL, err := FetchMetadata(pkgName, version)
	if err != nil {
		return "", "", err
	}

	fmt.Printf("Fetching %s@%s...\n", pkgName, targetVersion)

	// 3. Download tarball
	tarResp, err := http.Get(tarballURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to download tarball: %w", err)
	}
	defer tarResp.Body.Close()

	if tarResp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to download tarball, status: %s", tarResp.Status)
	}

	// 4. Set up cache file in node_modules/.vendor-cache
	safeName := strings.ReplaceAll(pkgName, "/", "-")
	cacheDir := filepath.Join("node_modules", ".vendor-cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", "", err
	}
	cachePath := filepath.Join(cacheDir, fmt.Sprintf("%s-%s.tgz", safeName, targetVersion))

	cf, err := os.Create(cachePath)
	if err != nil {
		return "", "", err
	}
	defer cf.Close()

	// 5. Extract tarball while streaming to cache file
	tee := io.TeeReader(tarResp.Body, cf)
	if err := extractTarGz(tee, destDir); err != nil {
		return "", "", err
	}

	return targetVersion, tarballURL, nil
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

// GetCachePath returns the expected file path of the cached tarball.
func GetCachePath(pkgName, version string) string {
	safeName := strings.ReplaceAll(pkgName, "/", "-")
	return filepath.Join("node_modules", ".vendor-cache", fmt.Sprintf("%s-%s.tgz", safeName, version))
}

// ExtractLocalTarGz extracts a local .tgz file to the specified directory.
func ExtractLocalTarGz(tarPath, destDir string) error {
	f, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return extractTarGz(f, destDir)
}

// Download fetches a file from a URL and saves it to the specified destination path.
func Download(url, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download tarball, status: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
