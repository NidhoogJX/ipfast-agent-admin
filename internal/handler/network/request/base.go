package request

import (
	"encoding/json"
	"net/url"

	"github.com/valyala/fasthttp"
)

func Get(urlStr string, queryParams map[string]string) (bodyBytes []byte, header []byte, err error) {
	parsedURL, err1 := url.Parse(urlStr)
	if err1 != nil {
		err = err1
		return
	}
	query := parsedURL.Query()
	for key, value := range queryParams {
		query.Set(key, value)
	}
	parsedURL.RawQuery = query.Encode()
	return MakeRequest(parsedURL.String(), "GET", nil)
}

func Post(url string, body interface{}) (bodyBytes []byte, header []byte, err error) {
	return MakeRequest(url, "POST", body)
}

func MakeRequest(url, method string, body interface{}) (bodyBytes []byte, header []byte, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	if method == "POST" && body != nil {
		jsonBody, err1 := json.Marshal(body)
		if err1 != nil {
			err = err1
			return
		}
		req.SetBody(jsonBody)
		req.Header.SetContentType("application/json")
	}
	err = fasthttp.Do(req, resp)
	bodyBytes = resp.Body()
	header = resp.Header.Header()
	return
}
