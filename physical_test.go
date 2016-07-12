package physical

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPhysical(t *testing.T) {
	Convey("Let's get physical", t, func() {
		wr := httptest.NewRecorder()
		path := "/healthcheck"
		router := mux.NewRouter().StrictSlash(true)
		request, _ := http.NewRequest("GET", path, nil)

		Convey("Creating", func() {
			CreateHealthCheck(router, path)

			router.ServeHTTP(wr, request)

			So(router.GetRoute("HealthCheck"), ShouldNotBeNil)
		})

		Convey("Adding a check", func() {
			CreateHealthCheck(router, path)
			check := func() HealthCheckResponse {
				return HealthCheckResponse{
					Actionable: true,
					Healthy:    true,
					Name:       "Sample Check",
					Type:       TypeSelf,
				}
			}
			AddCheck(check)
			router.ServeHTTP(wr, request)

			var resp collectedResponse
			body, _ := ioutil.ReadAll(io.LimitReader(wr.Body, 1048576))
			json.Unmarshal(body, &resp)

			So(wr.Code, ShouldEqual, http.StatusOK)
			So(resp.Healthy[0].Name, ShouldEqual, "Sample Check")
		})

		Convey("Failing check", func() {
			CreateHealthCheck(router, path)
			check := func() HealthCheckResponse {
				return HealthCheckResponse{
					Actionable: false,
					Healthy:    false,
					Name:       "Failing Check",
					Type:       TypeSelf,
				}
			}
			AddCheck(check)
			router.ServeHTTP(wr, request)

			var resp collectedResponse
			body, _ := ioutil.ReadAll(io.LimitReader(wr.Body, 1048576))
			json.Unmarshal(body, &resp)

			So(wr.Code, ShouldEqual, http.StatusInternalServerError)
			So(resp.Unhealthy[0].Name, ShouldEqual, "Failing Check")
		})

		Convey("Check with all response parameters set", func() {
			CreateHealthCheck(router, path)
			check := func() HealthCheckResponse {
				info := make(map[string]interface{})
				info["foo"] = "bar"
				return HealthCheckResponse{
					Actionable:     false,
					Healthy:        false,
					Name:           "Failing Check",
					Type:           TypeSelf,
					Severity:       SeverityCritical,
					Message:        "Something has gone really wrong!",
					Dependency:     Dependency{Name: "Upstream"},
					AdditionalInfo: info,
					Link:           "https://www.wolframalpha.com/input/?i=why+are+firetrucks+red%3F",
				}
			}
			AddCheck(check)
			router.ServeHTTP(wr, request)

			var resp collectedResponse
			body, _ := ioutil.ReadAll(io.LimitReader(wr.Body, 1048576))
			json.Unmarshal(body, &resp)

			So(wr.Code, ShouldEqual, http.StatusInternalServerError)
			unhealthy := resp.Unhealthy[0]
			So(unhealthy.Name, ShouldEqual, "Failing Check")
			So(unhealthy.Type, ShouldEqual, TypeSelf)
			So(unhealthy.Severity, ShouldEqual, SeverityCritical)
			So(unhealthy.Message, ShouldEqual, "Something has gone really wrong!")
			So(unhealthy.Dependency.Name, ShouldEqual, "Upstream")
			So(unhealthy.AdditionalInfo, ShouldContainKey, "foo")
			So(unhealthy.AdditionalInfo["foo"], ShouldEqual, "bar")
			So(unhealthy.Link, ShouldContainSubstring, "wolframalpha")
		})
	})
}
