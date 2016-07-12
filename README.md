# Physical

[![Go Report Card](https://goreportcard.com/badge/github.com/alde/physical)](https://goreportcard.com/report/github.com/alde/physical)

Health Checks for Go
Compatible with DropWizard-style health checks

## Purpose

Make a very simple way of adding custom health checks for your services.

## Usage

Given you use [gorilla-mux](https://github.com/gorilla/mux) you can use this as:

    import "github.com/alde/physical"
    ...

    ...
      path := "/healthcheck"
      router := mux.NewRouter().StrictSlash(true)
      physical.CreateHealthCheck(router, path)

      check := func() physical.HealthCheckResponse {
        return physical.HealthCheckResponse{
          Actionable: true,
          Healthy:    true,
          Name:       "Sample Check",
          Type:       physical.TypeSelf,
        }
      }
      physical.AddCheck(check)
    ...


AddCheck takes a `func() physical.HealthCheckResponse`

## Response Structure
Calls to the healthcheck endpoint will respond with:

Status: `200 OK` or `500 Internal Server Error` depending on whether there are failing checks or not.

Response Object:

    {
      "healthy": [
        ...
      ],
      "unhealthy": [
        ...
      ]
    }

the healthy or unhealthy fields contains:

    {
      // True or False depending on whether an action can be taken to resolve.
      "actionable": true,
      // True or False depending on whether a check is healthy
      "healthy": true,
      // Name of the check
      "name": "A Sample Check",
      // Type of the check. Constants are provided for:
      // SELF, METRICS, INFRASTRUCTURE, INTERNAL_DEPENDENCY, EXTERNAL_DEPENDENCY, INTERNET_CONNECTIVITY
      "type": "The type of check",
      // Severity of the check. Should be included on failing checks.
      // Constants are provided for: CRITICAL,WARNING, DOWN
      "severity": "CRITICAL",
      // A message attached to the check status. Should be included on failing checks.
      "message": "A message with more information",
      // Services this service is dependent on
      "dependent_on": {
        "service_name": "Name of upstream service"
      },
      // Additional information that can be useful
      "additional_info": {
        "foo": "bar"
      },
      // A link with more information, for example an incident matrix
      "link": "https://www.wolframalpha.com/input/?i=why+are+firetrucks+red%3F"
    }


# License

See [LICENSE](LICENSE) file.
