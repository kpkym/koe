package others

type Po struct {
	Type string `json:"type"`
	Ext  string `json:"ext"`
	Path string `json:"path"`
}

// Node represents a node in a directory tree.
type Node struct {
	Type             string  `json:"type,omitempty"`
	Title            string  `json:"title,omitempty"`
	UUID             string  `json:"uuid,omitempty"`
	WorkTitle        string  `json:"workTitle,omitempty"`
	MediaStreamUrl   string  `json:"mediaStreamUrl,omitempty"`
	MediaDownloadUrl string  `json:"mediaDownloadUrl,omitempty"`
	Duration         float64 `json:"duration,omitempty"`
	Path             string  `json:"abs,omitempty"`
	Children         []*Node `json:"children,omitempty"`
}
