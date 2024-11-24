package models

type RandomCodeReq struct {
	CountryCode string `json:"country_code"`
	Account     string `json:"account"`
}

type RandomCodeResp struct {
	RandomCode string `json:"randomCode"`
	Timestamp  int64  `json:"timestamp"`

	Meta
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
	Token string `json:"token"`

	Meta
}
