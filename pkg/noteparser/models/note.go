package models

type Note struct {
	FileVersion string
	Elements    map[string]string
	Layers      []Layer
}

type Layer struct {
	LayerType   string
	Protocol    string
	Name        string
	Path        int
	Bitmap      int
	VectorGraph int
	Recognition int
	Bytes       []byte
}
