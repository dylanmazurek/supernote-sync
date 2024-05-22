package models

type UserResp struct {
	Address       string `json:"address"`
	Avatar        string `json:"avatar"`
	Birthday      string `json:"birthday"`
	Education     string `json:"education"`
	Email         string `json:"email"`
	Hobby         string `json:"hobby"`
	Job           string `json:"job"`
	PersonalSign  string `json:"personalSign"`
	Telephone     string `json:"telephone"`
	CountryCode   string `json:"countryCode"`
	Sex           string `json:"sex"`
	TotalCapacity string `json:"totalCapacity"`
	UserName      string `json:"userName"`
	FileServer    string `json:"fileServer"`
	UserId        int64  `json:"userId"`

	Success bool `json:"success"`
	Error
}
