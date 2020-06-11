package httpclient

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"
)

const (
	RequestTimeout time.Duration = 10 * time.Second
)

type HTTPClient struct {
	*http.Client
}

func NewHTTPClient() *HTTPClient {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second,
	}

	return &HTTPClient{
		client,
	}
}

func (c *HTTPClient) UseSocksProxyNoPassword(address string) {
	dialer, err := proxy.SOCKS5("tcp", address, nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}

	c.Client = &http.Client{
		Transport: &http.Transport{
			Dial:            dialer.Dial,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second,
	}
}

func (c *HTTPClient) DoGet(URL string, header http.Header) (int, []byte, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return 0, nil, err
	}

	if header != nil {
		req.Header = header
	}

	resp, err := c.Do(req)
	if err != nil {
		return 0, nil, err
	}
	if resp.StatusCode != 200 {
		return resp.StatusCode, nil, fmt.Errorf("Fail with code %d", resp.StatusCode)
	}

	ct, err := readResponse(resp)
	return resp.StatusCode, ct, err
}

func (c *HTTPClient) DoGetWithParams(URL string, header http.Header, params map[string]interface{}) ([]byte, error) {
	var newURL = URL

	if params != nil {
		var arg = url.Values{}
		for k, v := range params {
			arg.Add(k, fmt.Sprintf("%v", v))
		}
		newURL += "?" + arg.Encode()
	}

	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Fail with code %d", resp.StatusCode)
	}

	ct, err := readResponse(resp)
	return ct, err
}

func (c *HTTPClient) DoGetWithParamsAndUnpackRespJson(URL string, header http.Header, params map[string]interface{}, respObjPointer interface{}) (err error) {
	data, err := c.DoGetWithParams(URL, header, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, respObjPointer)
	return
}

func GetMapValues(c interface{}) *url.Values {
	if m, ok := c.(map[string]interface{}); ok {
		structMap := &url.Values{}
		for k, v := range m {
			structMap.Add(k, fmt.Sprintf("%v", v))
		}
		return structMap
	}

	var t reflect.Type
	var v reflect.Value

	if reflect.TypeOf(c).Kind() != reflect.Struct {
		t = reflect.TypeOf(c).Elem()
		v = reflect.ValueOf(c).Elem()
	} else {
		t = reflect.TypeOf(c)
		v = reflect.ValueOf(c)
	}

	structMap := &url.Values{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key := field.Tag.Get("form")
		value := v.Field(i).Interface()
		structMap.Add(key, fmt.Sprintf("%v", value))

	}

	return structMap
}

func cloneHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

func (c *HTTPClient) DoPostForm(URL string, header http.Header, params interface{}) ([]byte, error) {
	var req *http.Request
	var resp *http.Response
	var err error

	var vars = GetMapValues(params)
	body := strings.NewReader(vars.Encode())
	req, err = http.NewRequest("POST", URL, body)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = cloneHeader(header)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		ct, _ := readResponse(resp)
		err = fmt.Errorf("Fail to request(%s): [%d] %s", URL, resp.StatusCode, string(ct))
		return nil, err
	}

	return readResponse(resp)
}

func (c *HTTPClient) DoPostFormAndUnpackRespJson(URL string, header http.Header, params interface{}, respObjPointer interface{}) (err error) {
	data, err := c.DoPostForm(URL, header, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, respObjPointer)
	return
}

func (c *HTTPClient) DoPostJSON(URL string, header http.Header, params interface{}) ([]byte, error) {
	var req *http.Request
	var resp *http.Response
	var err error
	var content string

	switch params.(type) {
	case []byte:
		content = string(params.([]byte))
	case string:
		content = params.(string)
	default:

		buf := bytes.NewBuffer(nil)
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)

		err = encoder.Encode(params)
		if err != nil {
			return nil, err
		}

		content = buf.String()
	}

	body := strings.NewReader(content)
	req, err = http.NewRequest("POST", URL, body)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		ct, _ := readResponse(resp)
		err = fmt.Errorf("Fail to request(%s): [%d] %s", URL, resp.StatusCode, string(ct))
		return nil, err
	}

	return readResponse(resp)
}

func (c *HTTPClient) DoPutJSON(URL string, header http.Header, params interface{}) ([]byte, error) {
	var req *http.Request
	var resp *http.Response
	var err error
	var content string

	switch params.(type) {
	case []byte:
		content = string(params.([]byte))
	case string:
		content = params.(string)
	default:

		buf := bytes.NewBuffer(nil)
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)

		err = encoder.Encode(params)
		if err != nil {
			return nil, err
		}

		content = buf.String()
	}

	body := strings.NewReader(content)
	req, err = http.NewRequest("PUT", URL, body)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		ct, _ := readResponse(resp)
		err = fmt.Errorf("HTTP [%d] Fail to request(%s) with req %v, resp %s", resp.StatusCode, URL, req, string(ct))
		return nil, err
	}

	return readResponse(resp)
}

func (c *HTTPClient) DoPostJsonWithRespAndData(URL string, header http.Header, params interface{}) (resp *http.Response, data []byte, err error) {
	var req *http.Request
	var content string

	switch params.(type) {
	case []byte:
		content = string(params.([]byte))
	case string:
		content = params.(string)
	default:
		buf := bytes.NewBuffer(nil)
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)

		err = encoder.Encode(params)
		if err != nil {
			return
		}

		content = buf.String()
	}

	body := strings.NewReader(content)
	req, err = http.NewRequest("POST", URL, body)
	if err != nil {
		return
	}

	if header != nil {
		req.Header = header
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err = c.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		ct, _ := readResponse(resp)
		err = fmt.Errorf("Fail to request(%s): [%d] %s", URL, resp.StatusCode, string(ct))
		return
	}

	data, err = readResponse(resp)
	return
}

func (c *HTTPClient) DoPostJsonAndUnpackRespJson(URL string, header http.Header, params interface{}, respObjPointer interface{}) (err error) {
	data, err := c.DoPostJSON(URL, header, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, respObjPointer)
	return
}

func (c *HTTPClient) DoDelete(URL string, header http.Header) ([]byte, error) {
	var req *http.Request
	var resp *http.Response
	var err error

	req, err = http.NewRequest("DELETE", URL, nil)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	}

	resp, err = c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		ct, _ := readResponse(resp)
		err = fmt.Errorf("Fail to HTTP DELETE %s: [%d] %s", URL, resp.StatusCode, string(ct))
		return nil, err
	}

	return readResponse(resp)
}

func readResponse(resp *http.Response) ([]byte, error) {
	var reader io.ReadCloser
	var err error
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = resp.Body
	}

	defer reader.Close()
	return ioutil.ReadAll(reader)
}

func paramsToString(params map[string]interface{}) string {
	values := url.Values{}
	for k, v := range params {
		if strV, ok := v.(string); ok {
			values.Set(k, strV)
		} else {
			values.Set(k, fmt.Sprintf("%v", v))
		}
	}

	return values.Encode()
}
