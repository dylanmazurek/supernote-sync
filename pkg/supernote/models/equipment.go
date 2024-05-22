package models

type EquipmentBindStatusReq struct {
}

type EquipmentBindStatusResp struct {
	Success    bool   `json:"success"`
	BindStatus bool   `json:"bindStatus"`
	BindValue  string `json:"bindValue"`

	Meta
}
