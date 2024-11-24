package models

type Device struct {
	AvailableMemory float64 `json:"availableMemory"`
	DeviceName      string  `json:"deviceName"`
	EntryList       []Entry `json:"fileList"`
	RouteList       []Route `json:"routeList"`
	TotalMemory     float64 `json:"totalMemory"`
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

type Route struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
