package zendesk

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"
)

const ZendeskApi = "https://lightburnsoftware.zendesk.com/api/v2/"

func TestGeturl(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		api       string
		endpoint  string
		opts      []GetOptions
		wantaddr  string
		wantparts []string
	}{
		{name: "api", host: "host.example.org", api: "v2", endpoint: "org", wantaddr: "https://host.example.org/v2/org", wantparts: []string{}},
		{name: "page", host: "host.example.org", api: "v2", endpoint: "org", opts: []GetOptions{WithPage(12)}, wantaddr: "https://host.example.org/v2/org", wantparts: []string{"page[size]=12"}},
		{name: "after", host: "host.example.org", api: "v2", endpoint: "org", opts: []GetOptions{StartAfter("42")}, wantaddr: "https://host.example.org/v2/org", wantparts: []string{}},
		{name: "page and after", host: "host.example.org", api: "v2", endpoint: "org", opts: []GetOptions{WithPage(33), StartAfter("42")}, wantaddr: "https://host.example.org/v2/org", wantparts: []string{"page[size]=33", "page[after]=42"}},
		{name: "byExternalId", host: "host.example.org", api: "v2", endpoint: "org", opts: []GetOptions{ByExternalId(159)}, wantaddr: "https://host.example.org/v2/org", wantparts: []string{"external_id=159"}},
		{name: "byName", host: "host.example.org:1", api: "v2", endpoint: "org", opts: []GetOptions{ByName("159")}, wantaddr: "https://host.example.org:1/v2/org", wantparts: []string{"name=159"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := geturl("https://", test.host, test.api, test.endpoint, test.opts...)
			t.Log(got)
			spl := strings.Split(got, "?")
			if len(spl) > 2 {
				t.Errorf("%s: Unexpected two or more ? in the url: %s", test.name, got)
			}
			gotaddr := spl[0]
			gotparts := []string{}
			if len(spl) > 1 {
				gotparts = strings.Split(spl[1], "&")
			}
			if gotaddr != test.wantaddr {
				t.Errorf("%s: expected address %s, got %s", test.name, gotaddr, test.wantaddr)
			}
			if len(gotparts) != len(test.wantparts) {
				t.Errorf("%s: expected %d http options, got %d", test.name, len(test.wantparts), len(gotparts))
				return
			}
			for _, part := range test.wantparts {
				if !slices.Contains(gotparts, part) {
					t.Errorf("%s: expected part %s", test.name, part)
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokens := strings.Split(r.Header.Get("Authorization"), " ")
		if len(tokens) < 2 || tokens[0] != "Bearer" {
			t.Errorf("Expected Authorization: Bearer ..., got: %s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"value":"fixed"}`))
	}))
	defer server.Close()

	t.Log(server.URL)
	u, err := url.ParseRequestURI(server.URL)
	if err != nil {
		t.Fatalf("TestGet: error parsing URL %s; error was %s", server.URL, err)
	}
	zd := &Zendesk{zdProtocol: "http://", zdHost: u.Host, token: "token"}
	data, err := zd.Get("")
	if err != nil {
		t.Errorf("%s: Error %s", "TestGet", err)
	}
	if string(data) != `{"value":"fixed"}` {
		t.Errorf("%s: unexpected response body %s", "TestGet", string(data))
	}
}
