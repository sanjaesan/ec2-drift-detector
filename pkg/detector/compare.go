package detector

import "reflect"

// getNestedValue retrieves a value from a nested map using dot notation
func getNestedValue(data map[string]any, path string) any {
	if data == nil {
		return nil
	}

	// Handle simple case
	if val, ok := data[path]; ok {
		return val
	}

	// Handle nested case (e.g., "network.subnet_id")
	keys := splitPath(path)
	current := any(data)

	for _, key := range keys {
		switch v := current.(type) {
		case map[string]any:
			if val, ok := v[key]; ok {
				current = val
			} else {
				return nil
			}
		default:
			return nil
		}
	}

	return current
}

// splitPath splits a path by dots
func splitPath(path string) []string {
	result := []string{}
	current := ""

	for _, char := range path {
		if char == '.' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}

// valuesEqual compares two values for equality
func valuesEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Handle slices
	aSlice, aIsSlice := a.([]any)
	bSlice, bIsSlice := b.([]any)

	if aIsSlice && bIsSlice {
		return slicesEqual(aSlice, bSlice)
	}

	// Handle string slices
	aStrSlice, aIsStrSlice := a.([]string)
	bStrSlice, bIsStrSlice := b.([]string)

	if aIsStrSlice && bIsStrSlice {
		return stringSlicesEqual(aStrSlice, bStrSlice)
	}

	// Handle maps
	aMap, aIsMap := a.(map[string]any)
	bMap, bIsMap := b.(map[string]any)

	if aIsMap && bIsMap {
		return mapsEqual(aMap, bMap)
	}

	// Use reflect for deep equality
	return reflect.DeepEqual(a, b)
}

func slicesEqual(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !valuesEqual(a[i], b[i]) {
			return false
		}
	}

	return true
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func mapsEqual(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}

	for key, aVal := range a {
		bVal, exists := b[key]
		if !exists || !valuesEqual(aVal, bVal) {
			return false
		}
	}

	return true
}