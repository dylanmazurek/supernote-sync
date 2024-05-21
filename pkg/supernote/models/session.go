package models

type Session struct {
	username string
	password string

	randomCode string

	token string
}

func (s *Session) GetCredentials() (string, string) {
	return s.username, s.password
}

func (s *Session) SetCredentials(username, password string) {
	s.username = username
	s.password = password
}

func (s *Session) GetToken() string {
	return s.token
}

func (s *Session) SetToken(token string) {
	s.token = token
}

func (s *Session) GetRandomCode() string {
	return s.randomCode
}

func (s *Session) SetRandomCode(randomCode string) {
	s.randomCode = randomCode
}
