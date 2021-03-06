package others

type Po struct {
	Type string `json:"type"`
	// Ext  string `json:"ext"`
	Path string `json:"path"`
}

// Node represents a node in a directory tree.
type Node struct {
	Type     string  `json:"type,omitempty"`
	Title    string  `json:"title,omitempty"`
	Code     string  `json:"code,omitempty"`
	UUID     uint32  `json:"uuid,omitempty"`
	Children []*Node `json:"children,omitempty"`
}
