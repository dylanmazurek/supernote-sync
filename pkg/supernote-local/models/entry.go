package models

type EntryList struct {
	DeviceName string  `json:"deviceName"`
	FileList   []Entry `json:"fileList"`
}

type EntryType int

const (
	Directory EntryType = iota
	File
)

type Entry struct {
	Date        string  `json:"date"`
	Extension   string  `json:"extension"`
	IsDirectory bool    `json:"isDirectory"`
	Name        string  `json:"name"`
	Size        float64 `json:"size"`
	Uri         string  `json:"uri"`

	Depth int
}

func (e Entry) GetType() EntryType {
	if e.IsDirectory {
		return Directory
	}

	return File
}
