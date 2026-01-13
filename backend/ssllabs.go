package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const maxAge = 86400

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func buildAnalyzeURL(hostName string, startNew bool, fromCache bool) string {

	basePath := "https://api.ssllabs.com/api/v2/analyze"
	u, _ := url.Parse(basePath)

	query := u.Query()
	query.Set("host", hostName)
	query.Set("all", "done")
	if startNew {
		query.Set("startNew", "on")
	}
	if fromCache {
		query.Set("fromCache", "on")
		query.Set("maxAge", strconv.Itoa(maxAge))
	}
	u.RawQuery = query.Encode()

	return u.String()
}

func parseJSONtoHost(jsonResponse []byte) (Host, error) {

	var host Host

	result := json.Unmarshal(jsonResponse, &host)
	if result != nil {
		fmt.Println("Error parsing JSON:", result)
		return Host{}, result
	}
	return host, nil

}

func analyze(hostName string, fromCache bool) (Host, error) {

	startNew := !fromCache
	apiURL := buildAnalyzeURL(hostName, startNew, fromCache)
	body, maxAssessments, currentAssessments, err := makeRequest(apiURL)
	if err != nil {
		return Host{}, fmt.Errorf("error fetching data: %w", err)
	}

	host, errHost := parseJSONtoHost(body)
	if errHost != nil {
		return Host{}, fmt.Errorf("error parsing host JSON: %w", errHost)
	}
	status := host.Status

	for status != "READY" && status != "ERROR" {
		if maxAssessments > 0 && currentAssessments >= maxAssessments {
			return Host{}, fmt.Errorf("maximum number of assessments reached. Please try again later")
		}

		apiURLNext := buildAnalyzeURL(hostName, false, fromCache)

		if status == "IN_PROGRESS" {
			time.Sleep(10 * time.Second)
		} else {
			time.Sleep(5 * time.Second)
		}

		body, maxAssessments, currentAssessments, err = makeRequest(apiURLNext)
		if err != nil {
			return Host{}, fmt.Errorf("error fetching data: %w", err)
		}

		host, errHost = parseJSONtoHost(body)
		if errHost != nil {
			return Host{}, fmt.Errorf("error parsing host JSON: %w", errHost)
		}
		status = host.Status
	}

	return host, nil
}

func makeRequest(apiURL string) (body []byte, maxAssessments int, currentAssessments int, err error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, 0, 0, err
	}
	defer resp.Body.Close()

	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		return nil, 0, 0, errBody
	}

	statusCode := resp.StatusCode

	maxAssessments, _ = strconv.Atoi(resp.Header.Get("X-Max-Assessments"))
	currentAssessments, _ = strconv.Atoi(resp.Header.Get("X-Current-Assessments"))

	if statusCode == 503 {
		time.Sleep(15 * time.Second)
		return body, maxAssessments, currentAssessments, &HTTPError{StatusCode: 503, Message: "the service is not available (e.g., down for maintenance)"}
	}
	if statusCode == 529 {
		time.Sleep(30 * time.Second)
		return body, maxAssessments, currentAssessments, &HTTPError{StatusCode: 529, Message: "the service is overloaded"}
	}

	switch statusCode {
	case 400:
		return body, maxAssessments, currentAssessments, &HTTPError{StatusCode: 400, Message: "invocation error (e.g., invalid parameters)"}
	case 429:
		return body, maxAssessments, currentAssessments, &HTTPError{StatusCode: 429, Message: "client request rate too high or too many new assessments too fast"}
	case 500:
		return body, maxAssessments, currentAssessments, &HTTPError{StatusCode: 500, Message: "internal error"}
	}

	return body, maxAssessments, currentAssessments, nil
}
