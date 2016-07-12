package physical

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type healthCheck struct {
	fun func() HealthCheckResponse
}

var healthChecks []*healthCheck

const (
	SeverityCritical = "CRITICAL"
	SeverityWarning  = "WARNING"
	SeverityDown     = "DOWN"

	TypeSelf                 = "SELF"
	TypeMetrics              = "METRICS"
	TypeInfrastructure       = "INFRASTRUCTURE"
	TypeInternalDependency   = "INTERNAL_DEPENDENCY"
	TypeExternalDependency   = "EXTERNAL_DEPENDENCY"
	TypeInternetConnectivity = "INTERNET_CONNECTIVITY"
)

// The HealthCheckResponse struct defines the return type of a HealthCheck call
type HealthCheckResponse struct {
	Actionable     bool                   `json:"actionable"`
	Healthy        bool                   `json:"healthy"`
	Name           string                 `json:"name"`
	Type           string                 `json:"type"`
	Severity       string                 `json:"severity,omitempty"`
	Message        string                 `json:"message,omitempty"`
	Dependency     Dependency             `json:"dependent_on,omitempty"`
	AdditionalInfo map[string]interface{} `json:"additional_info,omitempty"`
	Link           string                 `json:"link,omitempty"`
}

type Dependency struct {
	Name string `json:"service_name,omitempty"`
}

type collectedResponse struct {
	Healthy   []HealthCheckResponse `json:"healthy"`
	Unhealthy []HealthCheckResponse `json:"unhealthy"`
}

// CreateHealthCheck is used to create a healthcheck on the given route
func CreateHealthCheck(router *mux.Router, path string) {
	router.
		Methods("GET").
		Path(path).
		Name("HealthCheck").
		HandlerFunc(healthCheckHandler)
	healthChecks = make([]*healthCheck, 0)
}

// AddCheck is used to add more checks
func AddCheck(hc func() HealthCheckResponse) {
	check := &healthCheck{fun: hc}
	healthChecks = append(healthChecks, check)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	healthy, unhealthy := performChecks()

	hcr := &collectedResponse{
		Healthy:   healthy,
		Unhealthy: unhealthy,
	}
	body, err := json.Marshal(hcr)

	if err != nil || len(unhealthy) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Write(body)
}

func performChecks() ([]HealthCheckResponse, []HealthCheckResponse) {
	healthy := make([]HealthCheckResponse, 0)
	unhealthy := make([]HealthCheckResponse, 0)

	for _, hc := range healthChecks {
		h := hc.fun()
		if h.Healthy {
			healthy = append(healthy, h)
		} else {
			unhealthy = append(unhealthy, h)
		}
	}

	return healthy, unhealthy
}
