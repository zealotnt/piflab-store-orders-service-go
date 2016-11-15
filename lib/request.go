package lib

import (
	"github.com/parnurzeal/gorequest"

	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"
	"regexp"
)

type ResponseError struct {
	Error string `json:"error"`
}

func ParseError(body string) error {
	var err_parsed ResponseError
	json.Unmarshal([]byte(body), &err_parsed)
	return errors.New(err_parsed.Error)
}

func getUrlInfo(route string) (scheme string, host string) {
	re, _ := regexp.Compile(`(\w*):\/\/(.*)`)
	result := re.FindStringSubmatch(route)
	if result == nil {
		return "", ""
	}

	return result[1], result[2]
}

func RequestForwarder(r *http.Request, route string, data_unmarshal interface{}) (*http.Response, string, error) {
	r.URL.Scheme, r.URL.Host = getUrlInfo(route)
	r.RequestURI = ""

	response, err := (&http.Client{}).Do(r)
	if err != nil {
		return nil, "", err
	}

	all_data, _ := httputil.DumpResponse(response, true)
	header, _ := httputil.DumpResponse(response, false)

	body_bytes := all_data[len(header):]
	body_str := string(body_bytes)

	if err := json.Unmarshal([]byte(body_str), &data_unmarshal); err != nil {
		return nil, "", err
	}

	return response, body_str, nil
}

func HttpRequest(method string, route string, body interface{}) (*http.Response, string) {
	var resp *http.Response
	var resp_body string
	var errs []error

	request := gorequest.New()
	if method == "GET" {
		resp, resp_body, errs = request.Get(route).End()
	} else if method == "POST" {
		resp, resp_body, errs = request.Post(route).SendStruct(body).End()
	} else if method == "PUT" {
		resp, resp_body, errs = request.Put(route).SendStruct(body).End()
	} else if method == "DELETE" {
		resp, resp_body, errs = request.Delete(route).End()
	}
	if errs != nil {
		PR_DUMP(errs)
		return nil, ""
	}

	return resp, resp_body
}
