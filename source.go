package kasi

import (
	"io/ioutil"
	"net/http"
)

func GetFromSource(url string, r *http.Request) *HTTPResponse {
	req, err := http.NewRequest("GET", url, r.Body)
	response := HTTPResponse{}
	if err != nil {
		response.StatusCode = 0
		response.Error = err
		return &response
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		response.StatusCode = 0
		response.Error = err
		log.Error("%v", err)

		return &response
	}

	defer resp.Body.Close()

	response.StatusCode = resp.StatusCode
	response.Proto = resp.Proto
	response.Header = *&resp.Header
	response.Body = func() []byte {
		body, _ := ioutil.ReadAll(resp.Body)
		return body
	}()
	response.TransferEncoding = resp.TransferEncoding
	response.Trailer = resp.Trailer

	return &response
}
