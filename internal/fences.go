package internal

import "fmt"

func checkFences(repo *Repository, metafilePath string) error {
	var metadata struct {
		Inputs map[string]struct {
			Imports []struct {
				Path string `json:"path"`
			} `json:"imports"`
		} `json:"inputs"`
		Outputs map[string]struct {
			Inputs map[string]struct{} `json:"inputs"`
		}
	}
	if err := ReadJSON(metafilePath, &metadata); err != nil {
		return err
	}

	// Check for import cycles.
	{
		var breadcrumbs []string
		visited := make(map[string]bool)
		inside := make(map[string]bool)
		var rec func(key string) bool
		rec = func(key string) bool {
			if inside[key] {
				breadcrumbs = append(breadcrumbs, key)
				return true
			}
			if visited[key] {
				return false
			}
			depth := len(breadcrumbs)
			breadcrumbs = append(breadcrumbs, key)
			visited[key] = true
			inside[key] = true
			input := metadata.Inputs[key]
			for _, imprt := range input.Imports {
				if rec(imprt.Path) {
					return true
				}
			}
			breadcrumbs = breadcrumbs[:depth]
			delete(inside, key)
			return false
		}
	outputsLoop:
		for outputKey, output := range metadata.Outputs {
			breadcrumbs = append(breadcrumbs, outputKey)
			for inputKey := range output.Inputs {
				if rec(inputKey) {
					break outputsLoop
				}
			}
			breadcrumbs = breadcrumbs[:0]
		}

		n := len(breadcrumbs)
		if n > 0 {
			end := breadcrumbs[n-1]
			for i := 0; i < n; i++ {
				if breadcrumbs[i] == end {
					return fmt.Errorf("import cycle: %v", breadcrumbs[i:n-1])
				}
			}
		}
	}

	return nil
}
