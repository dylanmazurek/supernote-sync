package supernotecloud

import (
	"net/http"

	"github.com/dylanmazurek/supernote-sync/pkg/supernote-cloud/constants"
	"github.com/dylanmazurek/supernote-sync/pkg/supernote-cloud/models"
)

func (c *Client) GetUserInfo() (*models.UserResp, error) {
	req, err := c.NewRequest(http.MethodPost, constants.API_USER_QUERY, nil, nil)
	if err != nil {
		return nil, err
	}

	var userResp models.UserResp
	err = c.Do(req, &userResp)

	return &userResp, err
}
