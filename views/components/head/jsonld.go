package head

import (
	"encoding/json"
)

// JSONLDNode represents a generic JSON-LD object node.
// It includes common JSON-LD keywords and a map for additional properties.
type JSONLDNode struct {
	// Context specifies the JSON-LD context. It can be a URL string,
	// a map (inline context), or a slice of URL strings/maps.
	// Set to nil to omit.
	Context any `json:"-"` // Handled by custom MarshalJSON

	// ID is the IRI (Internationalized Resource Identifier) for the node.
	// Omitted if empty.
	ID string `json:"-"` // Handled by custom MarshalJSON

	// Type specifies the type(s) of the node. It can be a string or a slice of strings.
	// Omitted if nil, an empty string, or an empty slice.
	Type any `json:"-"` // Handled by custom MarshalJSON

	// Graph is used for representing a collection of nodes (a named graph or just a set of nodes).
	// Omitted if empty or nil.
	Graph []JSONLDNode `json:"-"` // Handled by custom MarshalJSON

	// Properties holds all other custom key-value pairs for the JSON-LD object.
	Properties map[string]any `json:"-"`
}

// MarshalJSON provides custom JSON marshaling for GenericJSONLDNode.
// It combines the explicit fields (Context, ID, Type, Graph) with the Properties map
// into a single JSON object, respecting common omitempty-like conditions.
func (n JSONLDNode) MarshalJSON() ([]byte, error) {
	out := make(map[string]any)

	// Copy properties from the Properties map first
	if n.Properties != nil {
		for k, v := range n.Properties {
			// Avoid overwriting core JSON-LD keywords if they are accidentally in Properties
			// and also explicitly set on the struct. Explicit struct fields take precedence.
			if k != "@context" && k != "@id" && k != "@type" && k != "@graph" {
				out[k] = v
			}
		}
	}

	// Add/overwrite with explicit fields, applying omitempty-like logic

	// @context
	if n.Context != nil {
		include := true
		// An empty string is not a valid context IRI.
		if cStr, ok := n.Context.(string); ok && cStr == "" {
			include = false
		}
		// Note: An inline context map or array of contexts could be empty (e.g., {} or []).
		// The JSON-LD spec says "An inline context SHOULD NOT be empty."
		// This marshaler will include them if provided and not nil/empty-string.
		// To omit, ensure n.Context is nil or an empty string.
		if include {
			out["@context"] = n.Context
		}
	}

	// @id
	if n.ID != "" { // Standard omitempty for string
		out["@id"] = n.ID
	}

	// @type
	if n.Type != nil {
		include := true
		if tStr, ok := n.Type.(string); ok && tStr == "" { // Empty string is not a valid @type
			include = false
		} else if tSliceStr, ok := n.Type.([]string); ok && len(tSliceStr) == 0 { // Empty slice of types
			include = false
		} else if tSliceIntf, ok := n.Type.([]any); ok && len(tSliceIntf) == 0 { // Empty slice of types (more generic)
			include = false
		}
		if include {
			out["@type"] = n.Type
		}
	}

	// @graph
	if len(n.Graph) > 0 { // Standard omitempty for slice
		out["@graph"] = n.Graph
	}

	return json.Marshal(out)
}
