// +build integration

package helper

import (
	"sync"

	"github.com/go-resty/resty/v2"
)

type HttpRequestResult struct {
	Response *resty.Response
	Error    error
}

type httpRequest struct {
	req    *resty.Request
	method string
	url    string
}

func NewHttpConcurrentHandler() *HttpConcurrentHandler {
	var wg sync.WaitGroup
	return &HttpConcurrentHandler{
		requests:  []*httpRequest{},
		startFlag: make(chan bool),
		wg:        &wg,
	}
}

type HttpConcurrentHandler struct {
	requests  []*httpRequest
	startFlag chan bool
	wg        *sync.WaitGroup
	results   chan *HttpRequestResult
}

func (h *HttpConcurrentHandler) RegisterRequest(req *resty.Request, method string, url string) {
	h.requests = append(h.requests, &httpRequest{
		req:    req,
		method: method,
		url:    url,
	})
}

func (h *HttpConcurrentHandler) RunAndWait() []*HttpRequestResult {
	n := len(h.requests)
	h.results = make(chan *HttpRequestResult, n)

	// register routines
	for _, r := range h.requests {
		h.wg.Add(1)
		rObj := r
		go func() {
			<-h.startFlag // wait until channel is closed
			defer h.wg.Done()
			resp, err := rObj.req.Execute(rObj.method, rObj.url)
			h.results <- &HttpRequestResult{
				Response: resp,
				Error:    err,
			}
		}()
	}

	// start
	close(h.startFlag)

	// wait
	h.wg.Wait()

	close(h.results)

	res := []*HttpRequestResult{}
	for r := range h.results {
		res = append(res, r)
	}
	return res
}
