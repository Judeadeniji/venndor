package manifest

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"
)

func writeFormattedJSON(path string, data []byte) error {
	opts := pretty.Options{
		Indent:   "  ",
		SortKeys: false,
	}
	formatted := pretty.PrettyOptions(data, &opts)
	// Add trailing newline to match standard formatting
	if len(formatted) > 0 && formatted[len(formatted)-1] != '\n' {
		formatted = append(formatted, '\n')
	}
	return os.WriteFile(path, formatted, 0644)
}

// EnsureWorkspace ensures that "vendor/*" is in the workspaces array of package.json.
func EnsureWorkspace(pkgJSONPath string) error {
	b, err := os.ReadFile(pkgJSONPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no package.json found at %s", pkgJSONPath)
		}
		return err
	}
	jsonStr := string(b)

	workspaces := gjson.Get(jsonStr, "workspaces")
	if !workspaces.Exists() {
		// Create a new array with "vendor/*"
		newJSON, err := sjson.Set(jsonStr, "workspaces", []string{"vendor/*"})
		if err != nil {
			return err
		}
		return writeFormattedJSON(pkgJSONPath, []byte(newJSON))
	}

	if !workspaces.IsArray() {
		return fmt.Errorf("unsupported workspaces format in package.json (expected array)")
	}

	found := false
	for _, w := range workspaces.Array() {
		if w.String() == "vendor/*" {
			found = true
			break
		}
	}

	if !found {
		// Append to the existing array
		newJSON, err := sjson.Set(jsonStr, "workspaces.-1", "vendor/*")
		if err != nil {
			return err
		}
		return writeFormattedJSON(pkgJSONPath, []byte(newJSON))
	}

	return nil
}

// EnsureImport ensures that an import map entry exists for the vendored package.
// e.g. "#vendor/pkg": "./vendor/pkg"
func EnsureImport(pkgJSONPath, pkgName string) error {
	b, err := os.ReadFile(pkgJSONPath)
	if err != nil {
		return err
	}
	jsonStr := string(b)

	alias := fmt.Sprintf("#vendor/%s", pkgName)
	target := fmt.Sprintf("./vendor/%s", pkgName)

	importsResult := gjson.Get(jsonStr, "imports")

	importsMap := make(map[string]string)
	if importsResult.Exists() {
		if importsResult.Type == gjson.JSON {
			if err := json.Unmarshal([]byte(importsResult.Raw), &importsMap); err != nil {
				return fmt.Errorf("failed to parse existing imports: %v", err)
			}
		} else {
			return fmt.Errorf("imports is not a JSON object")
		}
	}

	if val, ok := importsMap[alias]; ok && val == target {
		return nil // Already exists
	}

	importsMap[alias] = target

	// Marshal the imports map. We use a custom encoding to avoid escaping HTML characters like '&'.
	newImportsJSON, err := json.MarshalIndent(importsMap, "", "  ")
	if err != nil {
		return err
	}

	newJSON, err := sjson.SetRaw(jsonStr, "imports", string(newImportsJSON))
	if err != nil {
		return err
	}

	return writeFormattedJSON(pkgJSONPath, []byte(newJSON))
}

// RemoveImport removes the #vendor/<pkgName> import alias from package.json
func RemoveImport(packageJSONPath, pkgName string) error {
	b, err := os.ReadFile(packageJSONPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	jsonStr := string(b)

	alias := fmt.Sprintf("#vendor/%s", pkgName)

	// sjson requires escaping for special characters in keys like # and /
	// but actually sjson handles paths. Since alias has / and # we should just delete from the map.

	importsResult := gjson.Get(jsonStr, "imports")
	if !importsResult.Exists() || importsResult.Type != gjson.JSON {
		return nil
	}

	importsMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(importsResult.Raw), &importsMap); err != nil {
		return fmt.Errorf("failed to parse existing imports: %v", err)
	}

	if _, exists := importsMap[alias]; !exists {
		return nil // Not found
	}

	delete(importsMap, alias)

	if len(importsMap) == 0 {
		// Remove the whole imports block if empty
		newJSON, err := sjson.Delete(jsonStr, "imports")
		if err != nil {
			return err
		}
		return writeFormattedJSON(packageJSONPath, []byte(newJSON))
	}

	newImportsJSON, err := json.MarshalIndent(importsMap, "", "  ")
	if err != nil {
		return err
	}

	newJSON, err := sjson.SetRaw(jsonStr, "imports", string(newImportsJSON))
	if err != nil {
		return err
	}

	return writeFormattedJSON(packageJSONPath, []byte(newJSON))
}
