package supernotecloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dylanmazurek/supernote-sync/pkg/supernote-cloud/constants"
	"github.com/dylanmazurek/supernote-sync/pkg/supernote-cloud/models"
	"github.com/rs/zerolog/log"
)

func (c *Client) GetFileList(directoryId string, page int, pageSize int) (*models.FileListResp, error) {
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

func (c *Client) DownloadAllFiles(directoryId string, dest string) error {
	page := 1
	pageSize := 10

	topFolder, err := c.GetFileList(directoryId, page, pageSize)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get root folder")
	}

	for _, file := range topFolder.UserFiles {
		destFolder := fmt.Sprintf("%s/%s", dest, file.FileName)

		if !file.IsFolder {
			c.DownloadFile(&file, dest)
		} else {
			log.Info().Msgf("folder: %s", destFolder)

			c.DownloadAllFiles(file.ID, destFolder)
		}
	}

	return nil
}

func (c *Client) getFileUrl(fileId string) (*models.FileDownloadResp, error) {
	fileDownloadReq := models.FileDownloadReq{
		ID:   fileId,
		Type: 0,
	}

	reqJson, err := json.Marshal(fileDownloadReq)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest(http.MethodPost, constants.API_FILE_URL, bytes.NewReader(reqJson), nil)
	if err != nil {
		return nil, err
	}

	var resp models.FileDownloadResp
	err = c.Do(req, &resp)

	return &resp, err
}

func (c *Client) DownloadFile(file *models.UserFile, destination string) error {
	start := time.Now().UnixMilli()

	fileDownloadResp, err := c.getFileUrl(file.ID)
	if err != nil {
		return err
	}

	url, err := url.Parse(fileDownloadResp.URL)
	if err != nil {
		return err
	}

	resp, _ := http.DefaultClient.Do(
		&http.Request{
			Method: http.MethodGet,
			URL:    url,
		},
	)
	if resp.StatusCode != 200 {
		log.Error().Msgf("error while downloading")
	}
	defer resp.Body.Close()

	outPath := fmt.Sprintf("%s/%s", destination, file.FileName)

	tempPath := fmt.Sprintf("%s.tmp", outPath)
	f, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	progressReader := &models.ProgressReader{
		Reader: resp.Body,
		Size:   resp.ContentLength,
	}

	if _, err := io.Copy(f, progressReader); err != nil {
		log.Error().Err(err).Msgf("error while downloading")
	}

	os.Rename(tempPath, outPath)

	log.Info().Msgf("took: %.2fs\n", float64(time.Now().UnixMilli()-start)/1000)

	return nil
}
