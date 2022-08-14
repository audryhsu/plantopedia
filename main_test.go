package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// test the handler function in main.go
func TestHandler(t *testing.T) {
	// create test HTTP request
	req, err := http.NewRequest("GET", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	// recorder acts as target of http request
	recorder := httptest.NewRecorder()
	// create a http handler from handler function from main to be tested
	hf:= http.HandlerFunc(handler)

	// serve http request to recorder and executes handler
	hf.ServeHTTP(recorder, req)

	// Check response status code
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v instead of %v\n", status, http.StatusOK)
	}

	// check response body
	expected := `hello world`
	actual:= recorder.Body.String()

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v instead of %v\n", actual, expected)
	}
}

// test the /hello route
func TestRouter(t *testing.T)  {
	r:= newRouter()

	// create mock test server
	mockServer := httptest.NewServer(r)

	// mock server exposes location in URL attribute to make requests to
	resp, err := http.Get(mockServer.URL+"/hello")

	if err != nil {
		t.Fatal(err)
	}

	// Check response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be %v, instead go %v\n", http.StatusOK, resp.StatusCode)
	}

	// Check response body - read body from bytes to string a
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	
	expected :="hello world"
	respString := string(respBytes)
	if respString != expected {
		t.Errorf("Response should be %v instead got %v\n", expected, respString)
	}

}

func TestRouterForNonExistentRoute(t *testing.T)  {
	r:= newRouter()
	mockServer := httptest.NewServer(r)

	resp, err := http.Post(mockServer.URL+"/hello", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status should be 405, got %v\n", resp.StatusCode)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	respString := string(b)
	expected := ""
	if respString != expected {
		t.Fatalf("Expected %s, got %s\n", expected, respString)
	}
}

func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)

	// hit `Get /assets/` route to retrieve index.html file resp
	resp, err := http.Get(mockServer.URL+"/assets/")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be 300, got %v\n", resp.StatusCode)
	}

	// test content-type header is html 
	contentType := resp.Header.Get("Content-Type")
	expectedType := "text/html; charset=utf-8"

	if contentType != expectedType {
		t.Fatalf("Wrong content type, expected %s, got %s\n", expectedType, contentType)
	}
}