package models

type RandomCodeReq struct {
	CountryCode string `json:"country_code"`
	Account     string `json:"account"`
}

type RandomCodeResp struct {
	ErrorCode *string `json:"errorCode"`
	ErrorMsg  *string `json:"errorMsg"`

	Success    bool   `json:"success"`
	RandomCode string `json:"randomCode"`
	Timestamp  int64  `json:"timestamp"`
}

type LoginReq struct {
	Account     string `json:"account"`
	Browser     string `json:"browser"`
	CountryCode string `json:"country_code"`
	Equipment   string `json:"equipment"`
	Language    string `json:"language"`
	LoginMethod string `json:"loginMethod"`
	Password    string `json:"password"`
	Timestamp   int64  `json:"timestamp"`
}

type LoginResp struct {
	ErrorCode *string `json:"errorCode"`
	ErrorMsg  *string `json:"errorMsg"`

	Success bool   `json:"success"`
	Token   string `json:"token"`
}
