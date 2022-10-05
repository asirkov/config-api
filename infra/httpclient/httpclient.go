package httpclient

import (
	"io/ioutil"
	"net/http"
)

//go:generate mockgen.exe -source=httpclient.go -destination=./mock/mock_httpclient.go
type HttpClient interface {
	DoHttpGet(url string) ([]byte, int, error)
}

type HttpFetcher struct {
}

func NewHttpFetcher() *HttpFetcher {
	return &HttpFetcher{}
}

func (c *HttpFetcher) DoHttpGet(url string) ([]byte, int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var body []byte
	if resp.StatusCode == http.StatusOK {
		body, err = ioutil.ReadAll(resp.Body)
	}

	return body, resp.StatusCode, err
}
