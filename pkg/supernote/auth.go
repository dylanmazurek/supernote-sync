package supernote

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dylanmazurek/supernote-sync/pkg/supernote/constants"
	"github.com/dylanmazurek/supernote-sync/pkg/supernote/models"
	"github.com/rs/zerolog/log"
)

type AuthClient struct {
	internalClient *http.Client

	opts options

	session models.Session
}

func NewAuthClient(opts options) (*AuthClient, error) {
	authClient := &AuthClient{
		internalClient: &http.Client{Transport: http.DefaultTransport},

		opts: opts,
	}

	return authClient, nil
}

type addAuthHeaderTransport struct {
	T       http.RoundTripper
	Session models.Session
}

func (adt *addAuthHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", adt.Session.GetToken()))
	req.Header.Add("User-Agent", constants.REPO_URL)

	return adt.T.RoundTrip(req)
}

func (c *AuthClient) InitTransportSession() (*http.Client, error) {
	randomCode, err := c.GetCode()
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("random code: %s", *randomCode)

	token, err := c.RefreshSession()
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("token: %s", *token)

	authTransport, err := c.createAuthTransport()

	return authTransport, err
}

func (c *AuthClient) GetCode() (*string, error) {
	codeUrl := fmt.Sprintf("%s%s", constants.API_BASE_URL, constants.API_USER_RANDOM_CODE)

	var reqBody models.RandomCodeReq
	reqBody.CountryCode = "1"
	reqBody.Account = c.opts.username

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, codeUrl, bytes.NewBuffer(reqBodyJson))

	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", constants.REPO_URL)
	req.Header.Add("Withcredentials", "true")
	req.Header.Add("Origin", "https://cloud.supernote.com")
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.internalClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respJson models.RandomCodeResp
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		return nil, err
	}

	return &respJson.RandomCode, nil
}

func (c *AuthClient) RefreshSession() (*string, error) {
	codeUrl := fmt.Sprintf("%s%s", constants.API_BASE_URL, constants.API_USER_LOGIN)

	var reqBody models.LoginReq
	reqBody.Account = c.opts.username
	reqBody.Browser = "Chrome125"
	reqBody.CountryCode = "1"
	reqBody.Equipment = "1"
	reqBody.LoginMethod = "1"
	reqBody.Language = "en"
	reqBody.Password = c.opts.password
	reqBody.Timestamp = time.Now().Unix()

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, codeUrl, bytes.NewBuffer(reqBodyJson))

	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", constants.REPO_URL)
	req.Header.Add("Withcredentials", "true")
	req.Header.Add("Origin", "https://cloud.supernote.com")
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.internalClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respJson models.LoginResp
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		return nil, err
	}

	return &respJson.Token, nil
}

func (c *AuthClient) createAuthTransport() (*http.Client, error) {
	authClient := &http.Client{
		Transport: &addAuthHeaderTransport{
			T:       http.DefaultTransport,
			Session: c.session,
		},
	}

	return authClient, nil
}
