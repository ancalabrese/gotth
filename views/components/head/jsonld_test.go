package head

import (
	"encoding/json"
	"reflect"
	"testing"
)

// assertJSONEqual compares the marshaled JSON of a GenericJSONLDNode with an expected JSON string.
// It unmarshals both into maps for a more robust comparison that ignores key order.
func assertJSONEqual(t *testing.T, node JSONLDNode, expectedJSON string) {
	t.Helper()

	actualBytes, err := json.Marshal(node)
	if err != nil {
		t.Fatalf("Failed to marshal GenericJSONLDNode: %v", err)
	}

	var actualMap map[string]any
	if err := json.Unmarshal(actualBytes, &actualMap); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON bytes: %v\nActual JSON: %s", err, string(actualBytes))
	}
	// Ensure actualMap is non-nil if actualBytes was "{}", otherwise json.Unmarshal makes it nil.
	if string(actualBytes) == "{}" {
		actualMap = make(map[string]any)
	}

	var expectedMap map[string]any
	if err := json.Unmarshal([]byte(expectedJSON), &expectedMap); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON string: %v\nExpected JSON: %s", err, expectedJSON)
	}
	if expectedJSON == "{}" {
		expectedMap = make(map[string]any)
	}

	if !reflect.DeepEqual(actualMap, expectedMap) {
		// For better diff, re-marshal to indented JSON strings for output
		actualJSONIndented, _ := json.MarshalIndent(actualMap, "", "  ")
		expectedJSONIndented, _ := json.MarshalIndent(expectedMap, "", "  ")
		t.Errorf("JSON mismatch:\nGOT:\n%s\n\nWANT:\n%s", string(actualJSONIndented), string(expectedJSONIndented))
	}
}

func TestMarshal_BasicNode(t *testing.T) {
	node := JSONLDNode{
		Context: "http://schema.org/",
		ID:      "http://example.com/product/widget",
		Type:    "Product",
		Properties: map[string]any{
			"name":        "Awesome Widget",
			"description": "The best widget ever created.",
		},
	}
	expected := `{
		"@context": "http://schema.org/",
		"@id": "http://example.com/product/widget",
		"@type": "Product",
		"name": "Awesome Widget",
		"description": "The best widget ever created."
	}`
	assertJSONEqual(t, node, expected)
}

func TestMarshal_OmitEmptyFields(t *testing.T) {
	tests := []struct {
		name         string
		node         JSONLDNode
		expectedJSON string
	}{
		{
			name:         "Omit nil Context",
			node:         JSONLDNode{Context: nil, ID: "id1", Type: "Type1", Properties: map[string]any{"prop": "val"}},
			expectedJSON: `{"@id": "id1", "@type": "Type1", "prop": "val"}`,
		},
		{
			name:         "Omit empty string Context",
			node:         JSONLDNode{Context: "", ID: "id1", Type: "Type1"},
			expectedJSON: `{"@id": "id1", "@type": "Type1"}`,
		},
		{
			name:         "Omit empty string ID",
			node:         JSONLDNode{Context: "ctx", ID: "", Type: "Type1"},
			expectedJSON: `{"@context": "ctx", "@type": "Type1"}`,
		},
		{
			name:         "Omit nil Type",
			node:         JSONLDNode{Context: "ctx", ID: "id1", Type: nil},
			expectedJSON: `{"@context": "ctx", "@id": "id1"}`,
		},
		{
			name:         "Omit empty string Type",
			node:         JSONLDNode{Context: "ctx", ID: "id1", Type: ""},
			expectedJSON: `{"@context": "ctx", "@id": "id1"}`,
		},
		{
			name:         "Omit empty slice string Type",
			node:         JSONLDNode{Context: "ctx", ID: "id1", Type: []string{}},
			expectedJSON: `{"@context": "ctx", "@id": "id1"}`,
		},
		{
			name:         "Omit empty slice interface Type",
			node:         JSONLDNode{Context: "ctx", ID: "id1", Type: []any{}},
			expectedJSON: `{"@context": "ctx", "@id": "id1"}`,
		},
		{
			name:         "Omit nil Graph",
			node:         JSONLDNode{Graph: nil, Properties: map[string]any{"prop": "val"}},
			expectedJSON: `{"prop": "val"}`,
		},
		{
			name:         "Omit empty Graph",
			node:         JSONLDNode{Graph: []JSONLDNode{}, Properties: map[string]any{"prop": "val"}},
			expectedJSON: `{"prop": "val"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertJSONEqual(t, tc.node, tc.expectedJSON)
		})
	}
}

func TestMarshal_ComplexContextAndType(t *testing.T) {
	node := JSONLDNode{
		Context: []any{
			"https://schema.org",
			map[string]any{"ex": "http://example.com/vocab#"},
		},
		ID:   "http://example.com/item",
		Type: []string{"Thing", "ex:Custom"},
		Properties: map[string]any{
			"ex:property": "value",
		},
	}
	expected := `{
		"@context": ["https://schema.org", {"ex": "http://example.com/vocab#"}],
		"@id": "http://example.com/item",
		"@type": ["Thing", "ex:Custom"],
		"ex:property": "value"
	}`
	assertJSONEqual(t, node, expected)
}

func TestMarshal_GraphNode(t *testing.T) {
	node := JSONLDNode{
		Context: "http://schema.org/",
		Graph: []JSONLDNode{
			{
				ID:   "http://example.com/person/1",
				Type: "Person",
				Properties: map[string]any{
					"name": "Alice",
				},
			},
			{
				ID:   "http://example.com/person/2",
				Type: "Person",
				Properties: map[string]any{
					"name": "Bob",
				},
			},
		},
	}
	expected := `{
		"@context": "http://schema.org/",
		"@graph": [
			{
				"@id": "http://example.com/person/1",
				"@type": "Person",
				"name": "Alice"
			},
			{
				"@id": "http://example.com/person/2",
				"@type": "Person",
				"name": "Bob"
			}
		]
	}`
	assertJSONEqual(t, node, expected)
}

func TestMarshal_PropertiesOnly(t *testing.T) {
	node := JSONLDNode{
		Properties: map[string]any{
			"custom1": "value1",
			"custom2": 123,
		},
	}
	expected := `{
		"custom1": "value1",
		"custom2": 123
	}`
	assertJSONEqual(t, node, expected)
}

func TestMarshal_EmptyNode(t *testing.T) {
	node := JSONLDNode{} // All fields zero/nil, Properties is nil
	expected := `{}`
	assertJSONEqual(t, node, expected)

	nodeWithEmptyProperties := JSONLDNode{
		Properties: map[string]any{},
	}
	assertJSONEqual(t, nodeWithEmptyProperties, expected)
}

func TestMarshal_PropertiesWithCoreKeywords(t *testing.T) {
	// Core keywords in Properties should be ignored if the corresponding struct field is set,
	// or ignored entirely by the property copying logic.
	// The explicit struct fields always take precedence.
	node := JSONLDNode{
		ID:   "explicit-id", // This should win
		Type: "ExplicitType",
		Properties: map[string]any{
			"@id":      "id-in-properties",   // Should be ignored by property copy
			"@type":    "type-in-properties", // Should be ignored
			"name":     "Test Name",
			"@context": "context-in-properties", // Should be ignored
		},
	}
	// Since @context is not set on the struct field, it should not appear.
	// @id and @type from struct fields should be used.
	expected := `{
		"@id": "explicit-id",
		"@type": "ExplicitType",
		"name": "Test Name"
	}`
	assertJSONEqual(t, node, expected)

	// Case: explicit field is empty, property has core keyword (should still be ignored)
	node2 := JSONLDNode{
		ID: "", // Explicitly empty
		Properties: map[string]any{
			"@id":  "id-in-properties", // Should be ignored
			"name": "Another Name",
		},
	}
	expected2 := `{
		"name": "Another Name" 
	}` // No @id because explicit ID is empty, and properties @id is ignored
	assertJSONEqual(t, node2, expected2)
}

func TestMarshal_ContextMap(t *testing.T) {
	node := JSONLDNode{
		Context: map[string]any{
			"name": "http://schema.org/name",
			"xsd":  "http://www.w3.org/2001/XMLSchema#",
			"age":  map[string]any{"@id": "http://schema.org/age", "@type": "xsd:integer"},
		},
		ID: "http://example.com/thing",
	}
	expected := `{
		"@context": {
			"name": "http://schema.org/name",
			"xsd": "http://www.w3.org/2001/XMLSchema#",
			"age": {
				"@id": "http://schema.org/age",
				"@type": "xsd:integer"
			}
		},
		"@id": "http://example.com/thing"
	}`
	assertJSONEqual(t, node, expected)
}

func TestMarshal_ContextEmptyMap(t *testing.T) {
	// JSON-LD spec: "An inline context SHOULD NOT be empty."
	// Current MarshalJSON includes it if provided and not nil/empty string.
	node := JSONLDNode{
		Context: map[string]any{}, // Empty map context
		ID:      "http://example.com/item",
	}
	// It will be included by the current MarshalJSON logic as it's not nil and not an empty string
	expected := `{
        "@context": {}, 
        "@id": "http://example.com/item"
    }`
	assertJSONEqual(t, node, expected)
}

func TestMarshal_ContextEmptySlice(t *testing.T) {
	// Similar to empty map context
	node := JSONLDNode{
		Context: []any{}, // Empty slice context
		ID:      "http://example.com/item2",
	}
	expected := `{
        "@context": [],
        "@id": "http://example.com/item2"
    }`
	assertJSONEqual(t, node, expected)
}
