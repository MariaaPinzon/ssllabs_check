package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func buildAnalyzeURL(hostName string, startNew bool) string {

	basePath := "https://api.ssllabs.com/api/v2/analyze"
	u, _ := url.Parse(basePath)

	query := u.Query()
	query.Set("host", hostName)
	query.Set("all", "done")
	if startNew {
		query.Set("startNew", "on")
	}
	u.RawQuery = query.Encode()

	return u.String()
}

func parseJSONtoInfo(jsonResponse []byte) (Info, error) {

	var info Info

	result := json.Unmarshal(jsonResponse, &info)
	if result != nil {
		fmt.Println("Error parsing JSON:", result)
		return Info{}, result
	}
	return info, nil
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

func analyze(hostName string) (Host, error) {

	// HTTP request and response into Info struct
	infoURL := "https://api.ssllabs.com/api/v2/info"
	bodyInfo, _, errInfo := httpHelper(infoURL)
	if errInfo != nil {
		return Host{}, fmt.Errorf("error fetching info data: %w", errInfo)
	}

	info, errInfo := parseJSONtoInfo(bodyInfo)
	if errInfo != nil {
		return Host{}, fmt.Errorf("error parsing info JSON: %w", errInfo)
	}

	// HTTP request and response into Host struct
	apiURL := buildAnalyzeURL(hostName, true)
	body, _, err := httpHelper(apiURL)
	if err != nil {
		return Host{}, fmt.Errorf("error fetching data: %w", err)
	}

	host, errHost := parseJSONtoHost(body)
	if errHost != nil {
		return Host{}, fmt.Errorf("error parsing host JSON: %w", errHost)
	}
	status := host.Status

	for status != "READY" && status != "ERROR" {

		if info.CurrentAssessments >= info.MaxAssessments {
			return Host{}, fmt.Errorf("maximum number of assessments reached. Please try again later")
		}
		apiURLNext := buildAnalyzeURL(hostName, false)

		if status == "IN_PROGRESS" {
			time.Sleep(10 * time.Second)
		} else {
			time.Sleep(5 * time.Second)
		}
		body, _, err = httpHelper(apiURLNext)
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

func httpHelper(apiURL string) (body []byte, statusCode int, err error) {
	resp, err := http.Get(apiURL)

	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, errBody := io.ReadAll(resp.Body)
	statusCode = resp.StatusCode

	if errBody != nil {
		return nil, statusCode, errBody
	}

	switch statusCode {
	case 400:
		return body, statusCode, fmt.Errorf("invocation error (e.g., invalid parameters)")
	case 429:
		return body, statusCode, fmt.Errorf("client request rate too high or too many new assessments too fast")
	case 500:
		return body, statusCode, fmt.Errorf("internal error")
	case 503:
		return body, statusCode, fmt.Errorf("the service is not available (e.g., down for maintenance)")
	case 529:
		return body, statusCode, fmt.Errorf("the service is overloaded")
	}

	return body, statusCode, nil

}

func printInfo(host Host) {

	fmt.Printf("Host: %s\n", host.Host)
	fmt.Printf("Port: %d\n", host.Port)
	fmt.Printf("Protocol: %s\n", host.Protocol)
	fmt.Printf("IsPublic: %t\n", host.IsPublic)
	fmt.Printf("Status: %s\n", host.Status)
	fmt.Printf("StatusMessage: %s\n", host.StatusMessage)

	for _, endpoint := range host.Endpoints {
		fmt.Printf("  Endpoint IPAddress: %s\n", endpoint.IPAddress)
		fmt.Printf("  Endpoint StatusMessage: %s\n", endpoint.StatusMessage)
		fmt.Printf("  Endpoint Grade: %s\n", endpoint.Grade)
	}
}
