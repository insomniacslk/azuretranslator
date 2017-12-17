package azuretranslator

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	TOKEN_API     = "https://api.cognitive.microsoft.com/sts/v1.0/issueToken"
	TRANSLATE_API = "https://api.microsofttranslator.com/v2/Http.svc/Translate"
	DETECT_API    = "https://api.microsofttranslator.com/v2/Http.svc/Detect"
)

type DetectResponse struct {
	XMLName  xml.Name `xml:"string"`
	Language string   `xml:",chardata"`
}

type TranslatorClient struct {
	SessionToken []byte
	HttpClient   *http.Client
	Transport    *http.Transport
}

func NewTranslatorClient(apiKey string) (*TranslatorClient, error) {
	c := TranslatorClient{}
	token, err := c.getToken(apiKey)
	if err != nil {
		return nil, err
	}
	c.SessionToken = token
	return &c, nil
}

func (c *TranslatorClient) request(reqType, endpoint string, overrideHeaders map[string]string) ([]byte, error) {
	if c.Transport == nil {
		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    10 * time.Second,
			DisableCompression: false,
		}
		c.Transport = tr
	}
	headers := make(map[string]string, 3)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/jwt"
	for k, v := range overrideHeaders {
		headers[k] = v
	}
	if c.SessionToken != nil {
		headers["Authorization"] = fmt.Sprintf("Bearer %v", string(c.SessionToken))
	}
	if c.HttpClient == nil {
		c.HttpClient = &http.Client{Transport: c.Transport}
	}
	req, err := http.NewRequest(reqType, endpoint, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := c.HttpClient.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *TranslatorClient) getToken(apiKey string) ([]byte, error) {
	hdrs := make(map[string]string, 1)
	hdrs["Ocp-Apim-Subscription-Key"] = apiKey
	resp, err := c.request("POST", TOKEN_API, hdrs)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *TranslatorClient) Detect(text string) (string, error) {
	uri := fmt.Sprintf("%v?text=%v", DETECT_API, url.QueryEscape(text))
	resp, err := c.request("GET", uri, nil)
	if err != nil {
		return "", err
	}
	var ret DetectResponse
	err = xml.Unmarshal(resp, &ret)
	if err != nil {
		return "", err
	}
	return ret.Language, nil
}
