package restclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/version"
)

type Client struct {
	Services map[string]Endpoint
	Token    string
	Debug    bool
}

type Endpoint struct {
	ID    string
	Url   string
	token string
	debug bool
}

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per-page"`
	Pages   int `json:"pages"`
}

type PagResp struct {
	Pagination Pagination    `json:"pagination"`
	Data       []interface{} `json:"data"`
}

func NewClient(endpoints []Endpoint, token string, debug bool) *Client {
	services := map[string]Endpoint{}
	for _, e := range endpoints {
		e.token = token
		e.debug = debug
		services[e.ID] = e
	}

	return &Client{
		Services: services,
		Token:    token,
		Debug:    debug,
	}
}

func (e *Endpoint) Put(pathAction string, params url.Values, body string) (string, int, error) {
	resp, err := restCall(e.Url, e.token, pathAction, "PUT", params, http.Header{}, bytes.NewBufferString(body), e.debug)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	// check response code
	if resp.StatusCode >= 400 {
		return "", resp.StatusCode, errors.New(string(respBody))
	}

	return jsonPrettyPrint(string(respBody)), resp.StatusCode, nil
}

func (e *Endpoint) Post(pathAction string, params url.Values, header http.Header, body string) (string, int, error) {
	resp, err := restCall(e.Url, e.token, pathAction, "POST", params, header, bytes.NewBufferString(body), e.debug)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	// check response code
	if resp.StatusCode >= 400 {
		return "", resp.StatusCode, errors.New(string(respBody))
	}

	return jsonPrettyPrint(string(respBody)), resp.StatusCode, nil
}

func (e *Endpoint) GetList(pathAction string, params url.Values) ([]interface{}, int, error) {
	page := 1
	pages := 1
	per_page := 100
	result := []interface{}{}

	for i := 0; i < pages; i++ {
		// merge orig url values with the pagination
		helpers.MapMerge(params, url.Values{"page": []string{fmt.Sprintf("%d", page)}, "per_page": []string{fmt.Sprintf("%d", per_page)}})

		// get list entry
		pagData, _, err := e.getListEntry(pathAction, params)
		if err != nil {
			return nil, 0, err
		}

		// update pagination data
		if pagData.Pagination.Pages > 0 {
			pages = pagData.Pagination.Pages
		}
		if pagData.Pagination.PerPage > 0 {
			per_page = pagData.Pagination.PerPage
		}
		page++

		// add to the resutls
		result = append(result, pagData.Data...)
	}

	return result, 0, nil
}

func (e *Endpoint) Get(pathAction string, params url.Values, showPagination bool) (string, int, error) {
	resp, err := restCall(e.Url, e.token, pathAction, "GET", params, http.Header{}, nil, e.debug)
	if err != nil {
		return "", 0, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	// check response code
	if resp.StatusCode >= 400 {
		return "", resp.StatusCode, errors.New(string(respBody))
	}

	return jsonPrettyPrint(string(respBody)), resp.StatusCode, nil
}

func (e *Endpoint) Delete(pathAction string, params url.Values) (string, int, error) {
	resp, err := restCall(e.Url, e.token, pathAction, "DELETE", params, http.Header{}, nil, e.debug)
	if err != nil {
		return "", 0, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	// check response code
	if resp.StatusCode >= 400 {
		return "", resp.StatusCode, errors.New(string(respBody))
	}

	return jsonPrettyPrint(string(respBody)), resp.StatusCode, nil
}

// private

func (e *Endpoint) getListEntry(pathAction string, params url.Values) (*PagResp, int, error) {
	resp, err := restCall(e.Url, e.token, pathAction, "GET", params, http.Header{}, nil, e.debug)
	if err != nil {
		return nil, 0, err
	}

	// read body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// check response code
	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, errors.New(string(data))
	}

	// create a paginated response
	pagData := PagResp{}
	pagData.Pagination = extractPagination(resp.Header)
	err = helpers.JSONStringToStructure(string(data), &pagData.Data)
	if err != nil {
		pagData.Data = []interface{}{string(data)}
	}

	return &pagData, resp.StatusCode, nil
}

func restCall(endpoint string, token string, pathAction string, method string, params url.Values, headers http.Header, body *bytes.Buffer, debug bool) (*http.Response, error) {
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
	if debug {
		debugOutput(httputil.DumpRequestOut(req, true))
	}
	if err != nil {
		return nil, err
	}
	for k, entries := range headers {
		for _, v := range entries {
			req.Header.Add(k, v)
		}
	}
	req.Header.Add("User-Agent", fmt.Sprint("lyra-cli/", version.String()))
	req.Header.Add("X-Auth-Token", token)
	req.Header.Add("Content-Type", "application/json")

	// send the request
	resp, err := httpclient.Do(req)
	if debug {
		debugOutput(httputil.DumpResponse(resp, true))
	}
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

// Dump outgoing client requests. It includes any headers that the standard http.Transport adds, such as User-Agent.
func debugOutput(data []byte, err error) {
	var fmtError error
	// debug information will send to the stderr so it does not get mixed with the real output
	if err == nil {
		_, fmtError = fmt.Fprintf(os.Stderr, "%s\n\n", data)
	} else {
		_, fmtError = fmt.Fprintf(os.Stderr, "%s\n\n", err)
	}
	if fmtError != nil {
		log.Fatal(fmtError)
	}
}
