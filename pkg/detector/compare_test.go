package detector

import "testing"

func TestValuesEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		expected bool
	}{
		{"equal strings", "test", "test", true},
		{"different strings", "test1", "test2", false},
		{"equal string slices", []string{"a", "b"}, []string{"a", "b"}, true},
		{"different string slices", []string{"a", "b"}, []string{"a", "c"}, false},
		{"equal maps", map[string]interface{}{"key": "value"}, map[string]interface{}{"key": "value"}, true},
		{"different maps", map[string]interface{}{"key": "v1"}, map[string]interface{}{"key": "v2"}, false},
		{"nil values", nil, nil, true},
		{"one nil", "test", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := valuesEqual(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetNestedValue(t *testing.T) {
	data := map[string]interface{}{
		"simple": "value",
		"nested": map[string]interface{}{
			"key": "nested_value",
		},
	}

	tests := []struct {
		name     string
		path     string
		expected interface{}
	}{
		{"simple path", "simple", "value"},
		{"nested path", "nested.key", "nested_value"},
		{"non-existent", "non_existent", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getNestedValue(data, tt.path)
			if !valuesEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSplitPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected []string
	}{
		{"simple", "simple", []string{"simple"}},
		{"nested", "nested.key", []string{"nested", "key"}},
		{"deep", "a.b.c", []string{"a", "b", "c"}},
		{"empty", "", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitPath(tt.path)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}
		})
	}
}