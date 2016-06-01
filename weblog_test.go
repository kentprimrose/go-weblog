package weblog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

type MockHandler struct{}

func (mh MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Output for URL: %s", r.URL.Path)
}

func TestMockHandler(t *testing.T) {
	var mockWriter bytes.Buffer
	mockLogger := log.New(&mockWriter, "", log.Ldate|log.Ltime)

	handler := MockHandler{}
	server := httptest.NewServer(Handler(handler, mockLogger))
	defer server.Close()

	resp, err := http.Get(server.URL)

	if err != nil {
		t.Fatalf("http.Get error: %s", err.Error())
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d", resp.StatusCode)
	}

	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll error: %s", err.Error())
	}

	expected := "Output for URL: /"
	received := string(actual)
	if received != expected {
		t.Errorf("Expected '%s' but received '%s'", expected, received)
	}

	m := "\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} GET \\\"/\\\" 200 17 \\d{1,3}.\\d{1,3}"
	r, err := regexp.Compile(m)
	if err != nil {
		t.Fatalf("regexp.MatchString error: %s", err.Error())
	}

	received = mockWriter.String()
	if !r.MatchString(received) {
		t.Errorf("Logged '%s' doesn't match '%s'", received, m)
	}
}
