package controllers

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testRawData = []struct {
	Method   string
	URL      string
	Expected string
	Body     string
}{
	{
		Method:   "GET",
		URL:      "/v1/api/images",
		Expected: `{"status":"Ok"}`,
		Body:     "",
	},
	{
		Method:   "GET",
		URL:      "/v1/api/images/upload/invalid-content-job-id",
		Expected: `{"status":"Ok"}`,
		Body:     "",
	},
	{
		Method:   "POST",
		URL:      "/v1/api/images/upload",
		Expected: `{"status":"Ok"}`,
		Body:     "",
	},
	{
		Method:   "GET",
		URL:      "/v1/api/credentials/xcode",
		Expected: `{"status":"Ok"}`,
		Body:     "",
	},
}

/*
1. router := mux.NewRouter() //initialise the router
2. testServer := httptest.NewServer(router) //setup the testing server
3. request,error := http.NewRequest("METHOD","URL",Body)
4. //We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
        resp := httptest.NewRecorder()
5. handler := http.HandlerFunc(functionname)
// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
// directly and pass in our Request and ResponseRecorder.
        handler.ServeHTTP(resp, req)
*/
func TestHandlers(t *testing.T) {

	svc, _ := NewService(WithSvcOptHandler(&ApiHandler{}))

	/**
	r := chi.NewRouter()
	r.Get("/v1/api/images", api.GetAllImages)
	r.Get("/v1/api/images/upload/{id}", api.GetOneImage)
	r.Post("/v1/api/images/upload", api.UploadImage)
	r.Get("/v1/api/credentials/{code}", api.SetUserCode)
	**/

	ts := httptest.NewServer(svc.Router)

	for _, rec := range testRawData {
		_, body := testRequest(t, ts, rec.Method, rec.URL, bytes.NewBufferString(rec.Body))

		if body != rec.Expected {
			t.Fatalf("expected:%s got:%s", rec.Expected, body)
		}

	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, strings.TrimSpace(string(respBody))
}
