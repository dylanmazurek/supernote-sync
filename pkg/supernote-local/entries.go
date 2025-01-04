package supernotelocal

import (
	"encoding/json"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/dylanmazurek/supernote-sync/pkg/supernote-local/models"
)

func (c *Client) ListEntries(path string, depth *int) ([]models.Entry, error) {
	entries, err := c.getPage(path, 0)
	if err != nil {
		return nil, err
	}

	for i, entry := range entries {
		if depth == nil || entry.Depth == *depth {
			subEntries, err := c.getPage(entry.Uri, i+1)
			if err != nil {
				return nil, err
			}

			entries = append(entries, subEntries...)
		}
	}

	return entries, nil
}

func (c *Client) getPage(path string, depth int) ([]models.Entry, error) {
	req, err := c.NewRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var respDoc goquery.Document
	err = c.Do(req, &respDoc)
	if err != nil {
		return nil, err
	}

	var scriptRegex = regexp.MustCompile(`(?m)const json = '(.*)'$`)
	respObj := respDoc.Find("script").First().Text()
	scriptMatches := scriptRegex.FindStringSubmatch(respObj)
	jsonStr := scriptMatches[1]

	var entries models.EntryList
	err = json.Unmarshal([]byte(jsonStr), &entries)

	for i, _ := range entries.FileList {
		entries.FileList[i].Depth = depth
	}

	return entries.FileList, err
}
