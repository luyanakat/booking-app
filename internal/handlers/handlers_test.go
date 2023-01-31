package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	// GET
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"normal-room", "/normal-room", "GET", []postData{}, http.StatusOK},
	{"luxurious-room", "/luxurious-room", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"search-avai", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"make-reser", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"reservation-summary", "/reservation-summary", "GET", []postData{
		{key: "first_name", value: "nai"},
		{key: "last_name", value: "nải"},
		{key: "email", value: "nai@proton.me"},
		{key: "phone", value: "0123-456-789"},
	}, http.StatusMethodNotAllowed},
	// POST
	{"search-avai-post", "/search-availability", "POST", []postData{
		{key: "start", value: "16-10-2020"},
		{key: "end", value: "16-10-2020"},
	}, http.StatusOK},

	{"search-avai-json-post", "/search-availability-json", "POST", []postData{
		{key: "start", value: "16-10-2020"},
		{key: "end", value: "16-10-2020"},
	}, http.StatusOK},

	{"make-reser", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "nai"},
		{key: "last_name", value: "nải"},
		{key: "email", value: "nai@proton.me"},
		{key: "phone", value: "0123-456-789"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		// test reservation summary
		if e.url == "/reservation-summary" {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
		// end test reservation summary
		if e.method == "GET" && e.url != "/reservation-summary" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
