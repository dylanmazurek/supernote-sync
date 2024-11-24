package supernotelocal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/dylanmazurek/supernote-sync/pkg/supernote-local/models"
	"github.com/rs/zerolog/log"
)

type Client struct {
	internalClient *http.Client
	host           string
	port           int
}

func New(ctx context.Context, opts ...Option) (*Client, error) {
	clientOptions := defaultOptions()
	for _, opt := range opts {
		opt(&clientOptions)
	}

	newServiceClient := &Client{
		internalClient: http.DefaultClient,
		host:           clientOptions.host,
		port:           clientOptions.port,
	}

	return newServiceClient, nil
}

func (c *Client) NewRequest(method string, path string) (*models.Request, error) {
	urlString := fmt.Sprintf("http://%s:%d%s", c.host, c.port, path)
	requestUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, requestUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	request := &models.Request{
		HTTPRequest: req,
	}

	return request, nil
}

func (c *Client) Do(req *models.Request, resp *goquery.Document) error {
	httpResponse, err := c.internalClient.Do(req.HTTPRequest)
	if httpResponse.StatusCode >= 400 || err != nil {
		if httpResponse != nil {
			log.Debug().Msgf("http response error: %s", httpResponse.Status)
		}

		return err
	}
	defer httpResponse.Body.Close()

	doc, err := goquery.NewDocumentFromReader(httpResponse.Body)
	if err != nil {
		return err
	}

	*resp = *doc

	return nil
}
