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
	Hash             string  `json:"hash,omitempty"`
	WorkTitle        string  `json:"workTitle,omitempty"`
	MediaStreamUrl   string  `json:"mediaStreamUrl,omitempty"`
	MediaDownloadUrl string  `json:"mediaDownloadUrl,omitempty"`
	ImgUrl           string  `json:"imgUrl,omitempty"`
	LrcUrl           string  `json:"lrcUrl,omitempty"`
	Duration         float64 `json:"duration,omitempty"`
	Abs              string  `json:"abs,omitempty"`
	Children         []*Node `json:"children,omitempty"`
}
