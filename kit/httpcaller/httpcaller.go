// Package httpcaller provides a way to make http requests.
package httpcaller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// RequestParams holds the parameters to make an HTTP request.
type RequestParams struct {
	Client     *http.Client
	URL        string
	Headers    map[string]string
	Body       io.Reader
	Urlencoded map[string]string
	QueryParam map[string]string
}

// POST allows to execute an http POST request and unmarshal the result.
func POST(ctx context.Context, params RequestParams, dest any) (time.Duration, []byte, int, error) {
	start := time.Now()
	res, status, err := execute(ctx, http.MethodPost, params, dest)
	since := time.Since(start)

	return since, res, status, err
}

// GET allows to execute an http GET request and unmarshal the result.
func GET(ctx context.Context, params RequestParams, dest any) (time.Duration, []byte, int, error) {
	start := time.Now()
	res, status, err := execute(ctx, http.MethodGet, params, dest)
	since := time.Since(start)

	return since, res, status, err
}

// execute contains the logic to create an http request, run it, and validate the result.
func execute(ctx context.Context, method string, params RequestParams, dest any) ([]byte, int, error) {
	// parse URL
	u, err := url.Parse(params.URL)
	if err != nil {
		return nil, 0, err
	}

	query := u.Query()
	for key, value := range params.QueryParam {
		query.Add(key, value)
	}
	u.RawQuery = query.Encode()

	urlenCodeData := url.Values{}
	for key, value := range params.Urlencoded {
		urlenCodeData.Set(key, value)
	}

	// build request
	// nosemgrep: gosec.G107-1
	req, err := newRequest(method, u, params.Body, urlenCodeData)
	if err != nil {
		return nil, 0, err
	}

	// add the header values to the request
	for k, v := range params.Headers {
		req.Header.Set(k, v)
	}

	// do the request
	resp, err := params.Client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// read response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	// unmarshal the response
	if err = json.Unmarshal(data, &dest); err != nil {
		return data, resp.StatusCode, NewHTTPCallerError(err, resp.StatusCode, data)
	}

	return data, resp.StatusCode, nil
}

func Is2xxSuccessful(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func newRequest(method string, u *url.URL, body io.Reader, urlencoded url.Values) (*http.Request, error) {
	var req *http.Request
	var err error
	switch {
	case urlencoded.Encode() != "":
		req, err = http.NewRequest(method, u.String(), strings.NewReader(urlencoded.Encode()))
		if err != nil {
			return nil, err
		}

	default:
		req, err = http.NewRequest(method, u.String(), body)
		if err != nil {
			return nil, err
		}
	}

	return req, nil
}
