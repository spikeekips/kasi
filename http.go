package kasi

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/robertkrimen/otto"
	"github.com/spikeekips/kasi/conf"
	"github.com/spikeekips/kasi/util"
)

func HTTPServe(setting *conf.CoreSetting) {
	for bind, services := range setting.GetServicesByBind() {
		log.Debug("> create server for bind: %v", bind)
		serveMux := http.NewServeMux()

		for _, service := range services {
			for _, pattern := range service.GetPatterns() {
				log.Debug("\tregister pattern: %v: %v", service.GetID(), pattern)

				func(service *conf.ServiceSetting) {
					serveMux.HandleFunc(
						pattern,
						func(w http.ResponseWriter, r *http.Request) {
							HTTPServiceHandler(service, w, r)
						},
					)
				}(service)
			}
		}

		go func(bindAddress string, handler *http.ServeMux) {
			log.Debug("\tbind: %v", bindAddress)
			server := &http.Server{Addr: bindAddress, Handler: handler}
			log.Fatal(server.ListenAndServe())
		}(bind, serveMux)
	}
}

func HTTPServiceHandler(service *conf.ServiceSetting, w http.ResponseWriter, r *http.Request) {
	log.Debug("r(%p) %v", r, r.URL)
	log.Debug("r(%p) service id: %v", r, service.GetID())

	matched, err := service.GetMatchedEndpoint(r.URL.Path)
	if err != nil {
		log.Error(fmt.Sprintf("r(%p) %s", r, err))
		return
	}

	log.Debug("r(%p) endpoint id: %v", r, matched.GetID())

	url, err := matched.GetTargetURL(*r.URL)
	if err != nil {
		log.Error(fmt.Sprintf("r(%p) %s", r, err))
		return
	}

	returnedResponse := HTTPMiddlewarePreRequest(matched, r)
	if returnedResponse != nil {
		for header_key, value := range returnedResponse.Header {
			w.Header().Set(string(header_key), string(value))
		}

		w.WriteHeader(returnedResponse.StatusCode)
		w.Write([]byte(returnedResponse.Body))
		return
	}

	log.Debug("> make target url: %v -> %v", r.URL.String(), url)
	response := GetFromSource(url, r)

	returnedResponse = HTTPMiddlewarePostResponse(matched, r, response)
	if returnedResponse != nil {
		for header_key, value := range returnedResponse.Header {
			w.Header().Set(string(header_key), string(value))
		}

		w.WriteHeader(returnedResponse.StatusCode)
		w.Write([]byte(returnedResponse.Body))
		return
	}

	if response.HasError() {
		log.Error("r(%p) response from source has error, %v", response.Error)
		// if has some problems, just send HttpProblem, see
		// https://github.com/spikeekips/kasi/wiki/error
		w.WriteHeader(0)

		w.Write([]byte(""))
		return
	}

	log.Debug("r(%p) send response to client", r)
	for header_key, header_values := range response.Header {
		for _, value := range header_values {
			w.Header().Set(header_key, value)
		}
	}

	w.WriteHeader(response.StatusCode)

	w.Write(response.Body)
}

type MiddlewareResponse struct {
	StatusCode int
	Header     map[string]string
	Body       string
	Error      error
}

var RESPONSE_HEADER_KEY_MUST_BE_SKIPPED []string = []string{
	"Content-Length",
}

func NewMiddlewareResponseFromHTTPResponse(response *HTTPResponse) MiddlewareResponse {
	newHeader := map[string]string{}
	for header_key, value := range response.Header {
		if util.InArray(RESPONSE_HEADER_KEY_MUST_BE_SKIPPED, header_key) {
			continue
		}
		newHeader[header_key] = value[0]
	}

	return MiddlewareResponse{
		StatusCode: response.StatusCode,
		Header:     newHeader,
		Body:       string(response.Body),
		Error:      response.Error,
	}

}

func HTTPMiddlewarePreRequest(endpoint *conf.EndpointSetting, r *http.Request) *MiddlewareResponse {
	for i := 0; i < len(endpoint.Middleware); i++ {
		jsFile, err := ioutil.ReadFile(endpoint.Middleware[i])
		if err != nil {
			log.Error("javascript file, `%v`, not found", i)
			continue
		}
		vm := otto.New()
		vm.Run(jsFile)

		request := *r
		requestJS, _ := vm.ToValue(&request)

		response := MiddlewareResponse{StatusCode: 0, Header: map[string]string{}, Body: ""}
		responseJS, _ := vm.ToValue(&response)

		returnedResponse, err := vm.Call("process_request", nil, requestJS, responseJS)
		if err != nil {
			log.Error("failed to call, %v", err)
			continue
		}

		if !returnedResponse.IsDefined() {
			continue
		}

		exportedResponseReady, err := returnedResponse.Export()
		if err == nil {
			return exportedResponseReady.(*MiddlewareResponse)
		}
		break
	}

	return nil
}

func HTTPMiddlewarePostResponse(
	endpoint *conf.EndpointSetting,
	r *http.Request,
	response *HTTPResponse,
) *MiddlewareResponse {
	for i := len(endpoint.Middleware) - 1; i >= 0; i-- {
		jsFile, err := ioutil.ReadFile(endpoint.Middleware[i])
		if err != nil {
			log.Error("javascript file, `%v`, not found", i)
			continue
		}
		vm := otto.New()
		vm.Run(jsFile)

		request := *r
		requestJS, _ := vm.ToValue(&request)

		response := NewMiddlewareResponseFromHTTPResponse(response)
		responseJS, _ := vm.ToValue(&response)

		returnedResponse, err := vm.Call("process_response", nil, requestJS, responseJS)
		if err != nil {
			log.Error("failed to call, %v", err)
			continue
		}

		if !returnedResponse.IsDefined() {
			continue
		}

		exportedResponseReady, err := returnedResponse.Export()
		if err == nil {
			return exportedResponseReady.(*MiddlewareResponse)
		}
		break
	}

	return nil
}
