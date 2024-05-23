package models

import (
	"encoding/json"
)

type FileListReq struct {
	DirectoryId string `json:"directoryId"`
	PageNo      int    `json:"pageNo"`
	PageSize    int    `json:"pageSize"`
	Sequence    string `json:"sequence"`
	Order       string `json:"order"`
}

type FileListResp struct {
	UserFiles []UserFile `json:"userFileVOList"`

	PageSize   int `json:"size"`
	TotalFiles int `json:"total"`
	TotalPage  int `json:"pages"`

	Error
}

type UserFile struct {
	ID          string `json:"id"`
	DirectoryID string `json:"directoryId"`
	FileName    string `json:"fileName"`
	Size        int64  `json:"size"`
	MD5         string `json:"md5"`
	IsFolder    bool
	CreateTime  int64 `json:"createTime"`
	UpdateTime  int64 `json:"updateTime"`
}

func (u *UserFile) UnmarshalJSON(data []byte) error {
	type Alias UserFile
	aux := &struct {
		*Alias
		IsFolder string `json:"isFolder"`
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	u.IsFolder = aux.IsFolder == "Y"

	return nil
}

type FileDownloadReq struct {
	ID   string `json:"id"`
	Type int    `json:"type"`
}

type FileDownloadResp struct {
	MD5 string `json:"md5"`
	URL string `json:"url"`

	Meta
}
