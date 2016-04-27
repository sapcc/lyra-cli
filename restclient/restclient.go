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

func (c *Client) Get(pathAction string) (string, error) {
	u, err := url.Parse(c.Endpoint)
	u.Path = path.Join(u.Path, pathAction)

	httpclient := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("X-Auth-Token", c.Token)
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
