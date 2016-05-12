package restclient

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/sapcc/lyra-cli/helpers"
)

var AUTOMATION_URI = "https://automation-staging.***REMOVED***/api/v1/"

type Client struct {
	Services Services
	Token    string
}

type Endpoint struct {
	Url   string
	token string
}

type Services struct {
	Automation Endpoint
	Arc        Endpoint
}

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per-page"`
	Pages   int `json:"pages"`
}

type PagResp struct {
	Pagination Pagination  `json:"pagination"`
	Data       interface{} `json:"data"`
}

func NewClient(services Services, token string) *Client {
	services.Automation.token = token
	services.Arc.token = token
	return &Client{
		Services: services,
		Token:    token,
	}
}

func (e *Endpoint) Put(pathAction string, params url.Values, body string) (string, int, error) {
	resp, err := restCall(e.Url, e.token, pathAction, "PUT", params, bytes.NewBufferString(body))
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	return jsonPrettyPrint(string(respBody)), resp.StatusCode, nil
}

func (e *Endpoint) Post(pathAction string, params url.Values, body string) (string, int, error) {
	resp, err := restCall(e.Url, e.token, pathAction, "POST", params, bytes.NewBufferString(body))
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	return jsonPrettyPrint(string(respBody)), resp.StatusCode, nil
}

func (e *Endpoint) Get(pathAction string, params url.Values, showPagination bool) (string, int, error) {
	resp, err := restCall(e.Url, e.token, pathAction, "GET", params, nil)
	if err != nil {
		return "", 0, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	//extract pagination data
	if showPagination {
		// get pag params
		pagination := extractPagination(resp.Header)
		// create response
		pagData := PagResp{}
		// save pagination
		pagData.Pagination = pagination
		// save response
		err := helpers.JSONStringToStructure(string(data), &pagData.Data)
		if err != nil {
			return "", 0, err
		}
		data, err := json.Marshal(pagData)
		if err != nil {
			return "", 0, err
		}
		return jsonPrettyPrint(string(data)), resp.StatusCode, nil
	}

	return jsonPrettyPrint(string(data)), resp.StatusCode, nil
}

func restCall(endpoint string, token string, pathAction string, method string, params url.Values, body *bytes.Buffer) (*http.Response, error) {
	// set up the rest url
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, pathAction)
	u.RawQuery = params.Encode()

	// set up body
	var reqBody io.Reader
	if body != nil && body.Len() > 0 {
		reqBody = body
	}

	// set up the request
	httpclient := &http.Client{}
	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", token)
	req.Header.Add("Content-Type", "application/json")

	// send the request
	resp, err := httpclient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func extractPagination(header map[string][]string) Pagination {
	pag := Pagination{}
	if len(header["Pagination-Page"]) > 0 {
		i, err := strconv.Atoi(header["Pagination-Page"][0])
		if err == nil {
			pag.Page = i
		}
	}
	if len(header["Pagination-Per-Page"]) > 0 {
		i, err := strconv.Atoi(header["Pagination-Per-Page"][0])
		if err == nil {
			pag.PerPage = i
		}
	}
	if len(header["Pagination-Pages"]) > 0 {
		i, err := strconv.Atoi(header["Pagination-Pages"][0])
		if err == nil {
			pag.Pages = i
		}
	}
	return pag
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}
