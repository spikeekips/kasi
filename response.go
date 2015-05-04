package kasi

import "net/http"

type HTTPResponse struct {
	SourceResponse   *http.Response
	StatusCode       int
	Proto            string
	Header           map[string][]string
	Body             []byte
	TransferEncoding []string
	Trailer          map[string][]string
	Error            error
}

func (response *HTTPResponse) HasError() bool {
	return response.Error != nil
}
