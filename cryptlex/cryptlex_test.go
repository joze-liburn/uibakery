package cryptlex

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGetUrl(t *testing.T) {
	tests := []struct {
		name string
		cryp *Cryptlex
		api  []string
		want string
	}{
		{name: "normal", cryp: &Cryptlex{major: 3, host: "api.cryptlex.com"}, api: []string{"test1"}, want: "https://api.cryptlex.com/v3/test1"},
		{name: "normal-id", cryp: &Cryptlex{major: 3, host: "api.cryptlex.com"}, api: []string{"test1", "sedem"}, want: "https://api.cryptlex.com/v3/test1/sedem"},
		{name: "with-port", cryp: &Cryptlex{major: 3, host: "api.cryptlex.com:123"}, api: []string{"test2"}, want: "https://api.cryptlex.com:123/v3/test2"},
		{name: "http", cryp: &Cryptlex{major: 3, host: "http://localhost:42"}, api: []string{"test3"}, want: "http://localhost:42/v3/test3"},
		{name: "empty", cryp: &Cryptlex{major: 3, host: ""}, api: []string{"test4"}, want: "https:///v3/test4"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.cryp.geturl(test.api...)
			if got != test.want {
				t.Errorf("%s: geturl(%s), got %s, want %s", test.name, test.api, got, test.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		major      uint
		api        []string
		status     int
		body       []byte
		wantPath   string
		wantStatus int
		wantBody   []byte
	}{
		{
			name:       "normal",
			major:      7,
			api:        []string{"krneki"},
			status:     http.StatusOK,
			body:       []byte(`[{ "id": "1" }]`),
			wantPath:   "/v7/krneki",
			wantStatus: http.StatusOK,
			wantBody:   []byte(`[{ "id": "1" }]`),
		},
		{
			name:       "with-id",
			major:      7,
			api:        []string{"krneki", "external-id"},
			status:     http.StatusOK,
			body:       []byte(`[{ "id": "1" }]`),
			wantPath:   "/v7/krneki/external-id",
			wantStatus: http.StatusOK,
			wantBody:   []byte(`[{ "id": "1" }]`),
		},
		{
			name:       "404",
			major:      12,
			api:        []string{"123"},
			status:     http.StatusNotFound,
			wantPath:   "/v12/123",
			wantStatus: http.StatusNotFound,
			wantBody:   []byte(""),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != test.wantPath {
					t.Errorf("Expected to request %q, got: %q", test.wantPath, r.URL.Path)
				}
				w.WriteHeader(test.status)
				w.Write([]byte(test.body))
			}))
			defer server.Close()

			c := NewCryptlex(test.major, server.URL, "***")
			body, code, err := c.Get(test.api...)
			if code != test.wantStatus {
				t.Errorf("%s: http status: got %d, want %d", test.name, code, test.wantStatus)
			}
			if err != nil {
				t.Errorf("%s: error: got %v, want none", test.name, err)
			}
			if df := cmp.Diff(test.wantBody, body); df != "" {
				t.Errorf("%s: response body: -want +got\n%s", test.name, df)
			}
		})
	}
}

func TestPost(t *testing.T) {
	tests := []struct {
		name       string
		major      uint
		api        string
		token      string
		status     int
		body       []byte
		wantPath   string
		wantStatus int
		wantBody   []byte
		wantToken  []string
	}{
		{
			name:       "normal",
			major:      7,
			api:        "krneki",
			token:      "psst!",
			status:     http.StatusOK,
			body:       []byte(`[{ "id": "1" }]`),
			wantPath:   "/v7/krneki",
			wantStatus: http.StatusOK,
			wantBody:   []byte(`[{ "id": "1" }]`),
			wantToken:  []string{"Bearer psst!"},
		},
		{
			name:       "404",
			major:      12,
			api:        "123",
			token:      "***",
			status:     http.StatusNotFound,
			wantPath:   "/v12/123",
			wantStatus: http.StatusNotFound,
			wantBody:   []byte(""),
			wantToken:  []string{"Bearer ***"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != test.wantPath {
					t.Errorf("Expected to request %q, got: %q", test.wantPath, r.URL.Path)
				}
				tok := r.Header["Authorization"]
				lessFnc := func(p, q string) bool { return p < q }
				if df := cmp.Diff(test.wantToken, tok, cmpopts.SortSlices(lessFnc)); df != "" {
					t.Errorf("%s: http headers: Authorization: -want +got\n%s", test.name, df)
				}

				w.WriteHeader(test.status)
				w.Write([]byte(test.body))
			}))
			defer server.Close()

			c := NewCryptlex(test.major, server.URL, test.token)
			body, code, err := c.Post(test.api, test.body)
			if code != test.wantStatus {
				t.Errorf("%s: http status: got %d, want %d", test.name, code, test.wantStatus)
			}
			if err != nil {
				t.Errorf("%s: error: got %v, want none", test.name, err)
			}
			if df := cmp.Diff(test.wantBody, body); df != "" {
				t.Errorf("%s: response body: -want +got\n%s", test.name, df)
			}
		})
	}
}

func TestPatch(t *testing.T) {
	tests := []struct {
		name       string
		major      uint
		api        string
		id         string
		token      string
		status     int
		body       []byte
		wantPath   string
		wantStatus int
		wantBody   []byte
		wantToken  []string
	}{
		{
			name:       "normal",
			major:      7,
			api:        "krneki",
			id:         "id-1",
			token:      "psst!",
			status:     http.StatusOK,
			body:       []byte(`[{ "id": "1" }]`),
			wantPath:   "/v7/krneki/id-1",
			wantStatus: http.StatusOK,
			wantBody:   []byte(`[{ "id": "1" }]`),
			wantToken:  []string{"Bearer psst!"},
		},
		{
			name:       "404",
			major:      12,
			api:        "123",
			id:         "anonymous",
			token:      "***",
			status:     http.StatusNotFound,
			wantPath:   "/v12/123/anonymous",
			wantStatus: http.StatusNotFound,
			wantBody:   []byte(""),
			wantToken:  []string{"Bearer ***"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != test.wantPath {
					t.Errorf("Expected to request %q, got: %q", test.wantPath, r.URL.Path)
				}
				tok := r.Header["Authorization"]
				lessFnc := func(p, q string) bool { return p < q }
				if df := cmp.Diff(test.wantToken, tok, cmpopts.SortSlices(lessFnc)); df != "" {
					t.Errorf("%s: http headers: Authorization: -want +got\n%s", test.name, df)
				}

				w.WriteHeader(test.status)
				w.Write([]byte(test.body))
			}))
			defer server.Close()

			c := NewCryptlex(test.major, server.URL, test.token)
			body, code, err := c.Patch(test.api, test.id, test.body)
			if code != test.wantStatus {
				t.Errorf("%s: http status: got %d, want %d", test.name, code, test.wantStatus)
			}
			if err != nil {
				t.Errorf("%s: error: got %v, want none", test.name, err)
			}
			if df := cmp.Diff(test.wantBody, body); df != "" {
				t.Errorf("%s: response body: -want +got\n%s", test.name, df)
			}
		})
	}
}
