package restclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

var AUTOMATION_URI = "https://automation-staging.***REMOVED***/api/v1/"

type Client struct {
	Endpoint string
	Token    string
}

func NewClient(endpoint, token string) *Client {
	return &Client{
		Endpoint: endpoint,
		Token:    token,
	}
}

func (c *Client) Post(pathAction string, params url.Values, body string) (string, error) {
	u, err := url.Parse(c.Endpoint) //	u, err := url.Parse("http://localhost:3003/api/v1/")
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, pathAction)
	u.RawQuery = params.Encode()

	httpclient := &http.Client{}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBufferString(body))
	if err != nil {
		return "", err
	}
	req.Header.Add("X-Auth-Token", c.Token)
	req.Header.Add("Content-Type", "application/json")

	// log.Infof(fmt.Sprintf("RestClient Request: %+v", req))

	resp, err := httpclient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return jsonPrettyPrint(string(respBody)), nil
}

func (c *Client) Get(pathAction string, params url.Values) (string, error) {
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, pathAction)
	u.RawQuery = params.Encode()

	httpclient := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-Auth-Token", c.Token)

	// log.Infof(fmt.Sprintf("RestClient Request: %+v", req))

	resp, err := httpclient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return jsonPrettyPrint(string(body)), nil
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}
