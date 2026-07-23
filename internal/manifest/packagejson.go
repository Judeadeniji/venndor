package manifest

import (
	"encoding/json"
	"fmt"
	"os"
)

// EnsureWorkspace ensures that "vendor/*" is in the workspaces array of package.json.
func EnsureWorkspace(pkgJSONPath string) error {
	data, err := readJSON(pkgJSONPath)
	if err != nil {
		if os.IsNotExist(err) {
			// No package.json, nothing to do or create one? Let's just return error for now.
			return fmt.Errorf("no package.json found at %s", pkgJSONPath)
		}
		return err
	}

	workspacesRaw, ok := data["workspaces"]
	var workspaces []interface{}

	if ok {
		// workspaces can be an array or an object (e.g., yarn workspaces with 'packages')
		// We'll handle the array case for simplicity here.
		if wArr, isArr := workspacesRaw.([]interface{}); isArr {
			workspaces = wArr
		} else {
			return fmt.Errorf("unsupported workspaces format in package.json")
		}
	}

	// Check if "vendor/*" is already there
	found := false
	for _, w := range workspaces {
		if str, isStr := w.(string); isStr && str == "vendor/*" {
			found = true
			break
		}
	}

	if !found {
		workspaces = append(workspaces, "vendor/*")
		data["workspaces"] = workspaces
		return writeJSON(pkgJSONPath, data)
	}

	return nil
}

// EnsureImport ensures that an import map entry exists for the vendored package.
// e.g. "#vendor/pkg": "./vendor/pkg"
func EnsureImport(pkgJSONPath, pkgName string) error {
	data, err := readJSON(pkgJSONPath)
	if err != nil {
		return err
	}

	importsRaw, ok := data["imports"]
	var imports map[string]interface{}

	if ok {
		if impMap, isMap := importsRaw.(map[string]interface{}); isMap {
			imports = impMap
		} else {
			return fmt.Errorf("unsupported imports format in package.json")
		}
	} else {
		imports = make(map[string]interface{})
	}

	alias := fmt.Sprintf("#vendor/%s", pkgName)
	
	// We point it to the vendor directory. Node.js might require an exact file if not using a bundler,
	// but pointing to the directory works perfectly with modern bundlers and TS. 
	// To be perfectly Node-compliant, we might need subpath exports, but this is a good start.
	target := fmt.Sprintf("./vendor/%s", pkgName)

	if imports[alias] == target {
		return nil // Already exists
	}

	imports[alias] = target
	data["imports"] = imports

	return writeJSON(pkgJSONPath, data)
}

func readJSON(path string) (map[string]interface{}, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func writeJSON(path string, data map[string]interface{}) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	return enc.Encode(data)
}
