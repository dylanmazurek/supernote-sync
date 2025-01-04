package models

type Device struct {
	DeviceName  string  `json:"deviceName"`
	EntryList   []Entry `json:"fileList"`
	TotalMemory float64 `json:"totalMemory"`
	UsedMemory  int64   `json:"usedMemory"`
}
