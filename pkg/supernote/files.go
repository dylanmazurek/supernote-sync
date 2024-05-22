package supernote

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/dylanmazurek/supernote-sync/pkg/supernote/constants"
	"github.com/dylanmazurek/supernote-sync/pkg/supernote/models"
)

func (c *Client) GetFileList(directoryId int64, page int, pageSize int) (*models.FileListResp, error) {
	fileListReq := models.FileListReq{
		DirectoryId: directoryId,
		PageNo:      page,
		PageSize:    pageSize,
		Sequence:    "desc",
		Order:       "time",
	}

	reqJson, err := json.Marshal(fileListReq)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest(http.MethodPost, constants.API_FILE_LIST_QUERY, bytes.NewReader(reqJson), nil)
	if err != nil {
		return nil, err
	}

	var resp models.FileListResp
	err = c.Do(req, &resp)

	return &resp, err
}
