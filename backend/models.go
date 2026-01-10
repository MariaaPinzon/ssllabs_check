package main

type Host struct {
	Host            string     `json:"host"`
	Port            int        `json:"port"`
	Protocol        string     `json:"protocol"`
	IsPublic        bool       `json:"isPublic"`
	Status          string     `json:"status"`
	StatusMessage   string     `json:"statusMessage"`
	StartTime       int64      `json:"startTime"`
	TestTime        int64      `json:"testTime"`
	EngineVersion   string     `json:"engineVersion"`
	CriteriaVersion string     `json:"criteriaVersion"`
	CacheExpiryTime int64      `json:"cacheExpiryTime"`
	Endpoints       []Endpoint `json:"endpoints"`
	CertHostnames   []string   `json:"certHostnames"`
}

type Endpoint struct {
	IPAddress            string `json:"ipAddress"`
	ServerName           string `json:"serverName"`
	StatusMessage        string `json:"statusMessage"`
	StatusDetails        string `json:"statusDetails"`
	StatusDetailsMessage string `json:"statusDetailsMessage"`
	Grade                string `json:"grade"`
	GradeTrustIgnored    string `json:"gradeTrustIgnored"`
	HasWarnings          bool   `json:"hasWarnings"`
	IsExceptional        bool   `json:"isExceptional"`
	Progress             int    `json:"progress"`
	Duration             int64  `json:"duration"`
	Eta                  int    `json:"eta"`
	Delegation           int    `json:"delegation"`
}


