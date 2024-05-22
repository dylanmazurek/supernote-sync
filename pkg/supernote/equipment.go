package supernote

import (
	"net/http"

	"github.com/dylanmazurek/supernote-sync/pkg/supernote/constants"
	"github.com/dylanmazurek/supernote-sync/pkg/supernote/models"
)

func (c *Client) GetBindStatus() (*models.EquipmentBindStatusResp, error) {
	req, err := c.NewRequest(http.MethodPost, constants.API_EQUIPMENT_BIND_STATUS, nil, nil)
	if err != nil {
		return nil, err
	}

	var equipmentStatus models.EquipmentBindStatusResp
	err = c.Do(req, &equipmentStatus)

	return &equipmentStatus, err
}
