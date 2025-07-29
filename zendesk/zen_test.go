package zendesk

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func TestGeturl(t *testing.T) {
	tests := []struct {
		name  string
		api   string
		opts  []GetOptions
		parts []string
	}{
		{name: "api", api: "org", parts: []string{ZendeskApi + "org"}},
		{name: "page", api: "org", opts: []GetOptions{WithPage(12)}, parts: []string{ZendeskApi + "org", "page[size]=12"}},
		{name: "after", api: "org", opts: []GetOptions{StartAfter("42")}, parts: []string{ZendeskApi + "org"}},
		{name: "page and after", api: "org", opts: []GetOptions{WithPage(33), StartAfter("42")}, parts: []string{ZendeskApi + "org", "page[size]=33", "page[after]=42"}},
		{name: "byExternalId", api: "org", opts: []GetOptions{ByExternalId(159)}, parts: []string{ZendeskApi + "org", "external_id=159"}},
		{name: "byName", api: "org", opts: []GetOptions{ByName("159")}, parts: []string{ZendeskApi + "org", "name=159"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := geturl(ZendeskApi, test.api, test.opts...)
			t.Log(got)
			elts := strings.Split(got, "&")
			if len(elts) != len(test.parts) {
				t.Errorf("%s: expected %d http options, got %d", test.name, len(elts), len(test.parts))
				return
			}
			for _, part := range test.parts {
				if !slices.Contains(elts, part) {
					t.Errorf("%s: expected part %s", test.name, part)
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokens := strings.Split(r.Header.Get("Authorization"), " ")
		if len(tokens) < 2 || tokens[0] != "Basic" {
			t.Errorf("Expected Authorization: Basic ..., got: %s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"value":"fixed"}`))
	}))
	defer server.Close()

	zd := NewZendesk(server.URL, "token")
	data, err := zd.Get("")
	if err != nil {
		t.Errorf("%s: Error %s", "TestGet", err)
	}
	if string(data) != `{"value":"fixed"}` {
		t.Errorf("%s: unexpected response body %s", "TestGet", string(data))
	}
}
