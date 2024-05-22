package supernote

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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
	req.Header.Add("X-Access-Token", adt.Session.GetToken())
	req.Header.Add("User-Agent", constants.REPO_URL)

	return adt.T.RoundTrip(req)
}

func (c *AuthClient) InitTransportSession() (*http.Client, error) {
	c.session = models.Session{}
	c.session.SetCredentials(c.opts.username, c.opts.password)

	err := c.RefreshSession()
	if err != nil {
		return nil, err
	}

	authTransport, err := c.createAuthTransport()

	return authTransport, err
}

func (c *AuthClient) UpdateRandomCode() error {
	codeUrl := fmt.Sprintf("%s%s", constants.API_BASE_URL, constants.API_USER_RANDOM_CODE)

	var reqBody models.RandomCodeReq
	reqBody.CountryCode = "1"
	reqBody.Account = c.opts.username

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, codeUrl, bytes.NewBuffer(reqBodyJson))

	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", constants.REPO_URL)
	req.Header.Add("Withcredentials", "true")
	req.Header.Add("Origin", "https://cloud.supernote.com")
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.internalClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to get code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var respJson models.RandomCodeResp
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		return err
	}

	c.session.SetMetadata(respJson.RandomCode, respJson.Timestamp)

	log.Debug().Str("random_code", respJson.RandomCode).Int64("timestamp", respJson.Timestamp).Msg("metadata refreshed")

	return nil
}

func (c *AuthClient) RefreshSession() error {
	err := c.UpdateRandomCode()
	if err != nil {
		return err
	}

	codeUrl := fmt.Sprintf("%s%s", constants.API_BASE_URL, constants.API_USER_LOGIN)

	var reqBody models.LoginReq
	reqBody.Account = c.opts.username
	reqBody.Browser = "Chrome125"
	reqBody.CountryCode = "1"
	reqBody.Equipment = "1"
	reqBody.LoginMethod = "1"
	reqBody.Language = "en"

	password, timestamp := c.session.GetMetadata()
	reqBody.Password = password
	reqBody.Timestamp = timestamp

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, codeUrl, bytes.NewBuffer(reqBodyJson))

	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", constants.REPO_URL)
	req.Header.Add("Withcredentials", "true")
	req.Header.Add("Origin", "https://cloud.supernote.com")
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.internalClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to refresh session")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var respJson models.LoginResp
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		return err
	}

	c.session.SetToken(respJson.Token)

	log.Debug().Msg("session refreshed")

	return nil
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
