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

const maxAge = 86400 // 24 hours in seconds

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

// buildAnalyzeURL constructs the SSL Labs API URL for analyzing a domain's TLS/SSL configuration.
// It builds the query string with appropriate parameters based on whether a new analysis should be
// started or cached results should be used.
//
// Parameters:
//   - hostName: The domain name to analyze (e.g., "example.com")
//   - startNew: If true, starts a new analysis; if false, retrieves existing results
//   - fromCache: If true, allows using cached results from the server (max age: 24 hours)
//
// Returns:
//   - string: The complete SSL Labs API URL with encoded query parameters
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

// parseJSONtoHost parses a JSON response from the SSL Labs API into a Host struct.
// It deserializes the byte array containing the API response into a structured format.
//
// Parameters:
//   - jsonResponse: A byte array containing the JSON response from SSL Labs API
//
// Returns:
//   - Host: The parsed host information with TLS/SSL analysis details
//   - error: An error if JSON unmarshaling fails, nil otherwise
func parseJSONtoHost(jsonResponse []byte) (Host, error) {

	var host Host

	result := json.Unmarshal(jsonResponse, &host)
	if result != nil {
		return Host{}, result
	}
	return host, nil

}

// analyze performs TLS/SSL security analysis of a domain using the SSL Labs API.
// It initiates the analysis, polls the API until the analysis is complete, and returns the results.
//
// Parameters:
//   - hostName: The domain name to analyze (e.g., "example.com")
//   - fromCache: If true, attempts to retrieve cached results; if false, forces a new analysis
//
// Returns:
//   - Host: Host struct containing all analysis results including endpoint details
//   - error: An error if the API request fails, parsing fails, or assessment limits are reached
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

// makeRequest performs an HTTP GET request to the SSL Labs API with retry logic.
// The function also extracts rate limit information from response headers.
//
// Parameters:
//   - apiURL: The complete SSL Labs API URL to request
//
// Returns:
//   - body: The response body as a byte array
//   - maxAssessments: Maximum number of concurrent assessments allowed
//   - currentAssessments: Current number of active assessments
//   - err: An error if the request fails or returns a non-retryable error status
func makeRequest(apiURL string) (body []byte, maxAssessments int, currentAssessments int, err error) {
	const maxRetries = 3 // Maximum number of retries for 503 and 529 status codes
	retryCount := 0

	for {
		resp, err := http.Get(apiURL)
		if err != nil {
			return nil, 0, 0, err
		}

		body, errBody := io.ReadAll(resp.Body)
		resp.Body.Close()
		if errBody != nil {
			return nil, 0, 0, errBody
		}

		statusCode := resp.StatusCode
		maxAssessments, _ = strconv.Atoi(resp.Header.Get("X-Max-Assessments"))
		currentAssessments, _ = strconv.Atoi(resp.Header.Get("X-Current-Assessments"))

		if statusCode == 503 && retryCount < maxRetries {
			retryCount++
			time.Sleep(15 * time.Second)
			continue
		}
		if statusCode == 529 && retryCount < maxRetries {
			retryCount++
			time.Sleep(30 * time.Second)
			continue
		}

		switch statusCode {
		case 400:
			return body, maxAssessments, currentAssessments, &HTTPError{StatusCode: 400, Message: "invocation error (e.g., invalid parameters)"}
		case 429:
			return body, maxAssessments, currentAssessments, &HTTPError{StatusCode: 429, Message: "client request rate too high or too many new assessments too fast"}
		case 500:
			return body, maxAssessments, currentAssessments, &HTTPError{StatusCode: 500, Message: "internal error"}
		case 503:
			return body, maxAssessments, currentAssessments, &HTTPError{StatusCode: 503, Message: "the service is not available (e.g., down for maintenance)"}
		case 529:
			return body, maxAssessments, currentAssessments, &HTTPError{StatusCode: 529, Message: "the service is overloaded"}
		}

		return body, maxAssessments, currentAssessments, nil
	}
}
