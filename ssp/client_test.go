package ssp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/jsonapi"
)

func ExampleClient() {
	// Default client uses either the $HOME/.dashboard.env, or environment variable overrides.
	ssp, err := NewClient(nil)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	env, _ := ssp.GetEnvironment("mystack", "myenv")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	fmt.Printf("Just fetched environment %s", env.Name)

	d, _ := ssp.ApproveDeployment("mystack", "myenv", &ApproveDeployment{
		ID: 123,
	})
	fmt.Printf("Approved deployment %d", d.ID)
}

func TestNewApi(t *testing.T) {

	api, _ := NewClient(&Config{
		BaseURL: "http://localhost",
	})

	if api.baseURL.Host != "localhost" {
		t.Error("BaseURL not parsed correctly")
	}
}

func TestRequestErrorMessage(t *testing.T) {
	api, ts := newMockDashboard(nil, http.StatusBadGateway)
	defer ts.Close()
	_, err := api.request("GET", "/some/path", nil)
	if err == nil {
		t.Errorf("Expected error, got nil error")
	}

	actual := err.Error()
	expected := fmt.Sprintf("GET /some/path | HTTP 502 - '502 Bad Gateway'")
	if actual != expected {
		t.Errorf("Expected message \"%s\", got \"%s\"", expected, actual)
	}
}

func newMockDashboard(in interface{}, responseCode int) (*Client, *httptest.Server) {

	responseHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", jsonapi.MediaType)
		w.WriteHeader(responseCode)
		jsonapi.MarshalPayload(w, in)
	}

	ts := httptest.NewServer(http.HandlerFunc(responseHandler))

	api, err := NewClient(&Config{BaseURL: ts.URL})
	if err != nil {
		panic("Something bad happened while creating the API")
	}
	return api, ts
}
