package supernotelocal

import (
	"github.com/dylanmazurek/supernote-sync/pkg/supernote-local/models"
)

func (c *Client) GetEntries(path string, depth int) ([]models.Entry, error) {
	rootEntries, err := c.getEntry(path, depth)
	if err != nil {
		return nil, err
	}

	for _, entry := range rootEntries {
		if entry.Depth == depth {
			return []models.Entry{entry}, nil
		}
	}

	// return device.EntryList, nil
	return nil, nil
}

func (c *Client) getEntry(path string, depth int) ([]models.Entry, error) {
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
