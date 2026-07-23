package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type VendorYAML struct {
	Version  int                    `yaml:"version"`
	Packages map[string]*PackageYML `yaml:"packages"`
}

type PackageYML struct {
	Version string `yaml:"version"`
	Source  string `yaml:"source"`
	Path    string `yaml:"path"`
	Import  string `yaml:"import"`
	Patched bool   `yaml:"patched"`
	Notes   string `yaml:"notes"`
}

type VendorLock struct {
	LockfileVersion int                     `json:"lockfileVersion"`
	Packages        map[string]*PackageLock `json:"packages"`
}

type PackageLock struct {
	Version      string `json:"version"`
	Resolved     string `json:"resolved"`
	Integrity    string `json:"integrity,omitempty"`
	BaselineHash string `json:"baselineHash,omitempty"`
	VendoredAt   string `json:"vendoredAt"`
	PatchFile    string `json:"patchFile,omitempty"`
}

// RecordManifest writes entries to vendor.yaml and vendor-lock.json.
func RecordManifest(pkgName, version, tarballURL, destDir string) error {
	importAlias := fmt.Sprintf("#vendor/%s", pkgName)

	if err := ensureVendorYAML(pkgName, version, destDir, importAlias); err != nil {
		return fmt.Errorf("failed to write vendor.yaml: %w", err)
	}

	if err := ensureVendorLock(pkgName, version, tarballURL); err != nil {
		return fmt.Errorf("failed to write vendor-lock.json: %w", err)
	}

	return nil
}

func ensureVendorYAML(pkgName, version, path, importAlias string) error {
	yamlPath := "vendor.yaml"

	var config VendorYAML
	b, err := os.ReadFile(yamlPath)
	if err == nil {
		if err := yaml.Unmarshal(b, &config); err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		config = VendorYAML{
			Version:  1,
			Packages: make(map[string]*PackageYML),
		}
	} else {
		return err
	}

	if config.Packages == nil {
		config.Packages = make(map[string]*PackageYML)
	}

	config.Packages[pkgName] = &PackageYML{
		Version: version,
		Source:  "npm",
		Path:    path,
		Import:  importAlias,
		Patched: false,
		Notes:   "",
	}

	out, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	return os.WriteFile(yamlPath, out, 0644)
}

func ensureVendorLock(pkgName, version, tarballURL string) error {
	lockPath := "vendor-lock.json"

	var lock VendorLock
	b, err := os.ReadFile(lockPath)
	if err == nil {
		if err := json.Unmarshal(b, &lock); err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		lock = VendorLock{
			LockfileVersion: 1,
			Packages:        make(map[string]*PackageLock),
		}
	} else {
		return err
	}

	if lock.Packages == nil {
		lock.Packages = make(map[string]*PackageLock)
	}

	lock.Packages[pkgName] = &PackageLock{
		Version:    version,
		Resolved:   tarballURL,
		VendoredAt: time.Now().UTC().Format(time.RFC3339),
	}

	out, err := json.MarshalIndent(&lock, "", "  ")
	if err != nil {
		return err
	}

	// Add trailing newline
	out = append(out, '\n')
	return os.WriteFile(lockPath, out, 0644)
}

func GetPackageConfig(pkgName string) (*PackageYML, error) {
	b, err := os.ReadFile("vendor.yaml")
	if err != nil {
		return nil, err
	}
	var config VendorYAML
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	pkg, ok := config.Packages[pkgName]
	if !ok {
		return nil, fmt.Errorf("package %s not found in vendor.yaml", pkgName)
	}
	return pkg, nil
}

func GetPackageLock(pkgName string) (*PackageLock, error) {
	b, err := os.ReadFile("vendor-lock.json")
	if err != nil {
		return nil, err
	}
	var lock VendorLock
	if err := json.Unmarshal(b, &lock); err != nil {
		return nil, err
	}
	pkg, ok := lock.Packages[pkgName]
	if !ok {
		return nil, fmt.Errorf("package %s not found in vendor-lock.json", pkgName)
	}
	return pkg, nil
}

func MarkPatched(pkgName, patchFile string) error {
	// Update vendor.yaml
	b, err := os.ReadFile("vendor.yaml")
	if err == nil {
		var config VendorYAML
		if err := yaml.Unmarshal(b, &config); err == nil {
			if pkg, ok := config.Packages[pkgName]; ok {
				pkg.Patched = true
				out, _ := yaml.Marshal(&config)
				os.WriteFile("vendor.yaml", out, 0644)
			}
		}
	}

	// Update vendor-lock.json
	b2, err := os.ReadFile("vendor-lock.json")
	if err == nil {
		var lock VendorLock
		if err := json.Unmarshal(b2, &lock); err == nil {
			if pkg, ok := lock.Packages[pkgName]; ok {
				pkg.PatchFile = patchFile
				out, _ := json.MarshalIndent(&lock, "", "  ")
				out = append(out, '\n')
				os.WriteFile("vendor-lock.json", out, 0644)
			}
		}
	}

	return nil
}
