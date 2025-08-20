package cryptlex

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	resellerAsJson = `{
"name": "Creativo3D",
"email": "info@creativo3d.com",
"description": "SourceID: gid://shopify/Company/2185199784",
"allowedUsers": 0,
"allowedOrganizations": 0,
"id": "0197624c-be2b-705c-b866-ca4921844c4d",
"createdAt": "2025-06-12T04:01:26.059037Z",
"updatedAt": "2025-06-12T04:01:26.059044Z"
}`
	resellerAsObject = Reseller{
		Name:                 "Creativo3D",
		Email:                "info@creativo3d.com",
		Description:          "SourceID: gid://shopify/Company/2185199784",
		AllowedUsers:         0,
		AllowedOrganizations: 0,
		Id:                   "0197624c-be2b-705c-b866-ca4921844c4d",
		CreatedAt:            time.Date(2025, 6, 12, 4, 1, 26, 59037000, time.UTC),
		UpdatedAt:            time.Date(2025, 6, 12, 4, 1, 26, 59044000, time.UTC),
	}
)

func TestJsonToReseller(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    Reseller
		wanterr error
	}{
		{
			name: "real",
			json: resellerAsJson,
			want: resellerAsObject,
		},
		{
			name:    "bad json",
			json:    `{ 1:2 }`,
			want:    Reseller{},
			wanterr: errResellerJson,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := jsonToReseller([]byte(test.json))
			if !errors.Is(err, test.wanterr) {
				t.Errorf("%s: Reseller Unmarshal error: got %v, want %v", test.name, err, test.wanterr)
			}
			if df := cmp.Diff(test.want, got); df != "" {
				t.Errorf("%s: Reseller -want +got\n%s", test.name, df)
			}
		})
	}
}

func TestJsonToResellers(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    []Reseller
		wanterr error
	}{
		{
			name: "real",
			json: "[" + resellerAsJson + "]",
			want: []Reseller{resellerAsObject},
		},
		{
			name:    "bad json",
			json:    `{ 1:2 }`,
			want:    []Reseller{},
			wanterr: errResellerJson,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := jsonToResellers([]byte(test.json))
			if !errors.Is(err, test.wanterr) {
				t.Errorf("%s: Reseller Unmarshal error: got %v, want %v", test.name, err, test.wanterr)
			}
			if df := cmp.Diff(test.want, got); df != "" {
				t.Errorf("%s: Reseller -want +got\n%s", test.name, df)
			}
		})
	}
}

func TestJsonToError(t *testing.T) {
	tests := []struct {
		name string
		json string
		code int
		want error
	}{
		{
			name: "401",
			json: `{ "message": "unauthorized", "code": "401" }`,
			want: errResellerHttp,
		},
		{
			name: "no code",
			json: `{ "message": "unauthorized" }`,
			want: errResellerHttp,
		},
		{
			name: "no error",
			json: `{  }`,
			want: nil,
		},
		{
			name: "bad json",
			json: `{ 1-:2 }`,
			want: errResellerJson,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := jsonToError([]byte(test.json))
			if !errors.Is(err, test.want) {
				t.Errorf("%s: Reseller Unmarshal error: got %v, want %v", test.name, err, test.want)
			}
			if df := cmp.Diff(Reseller{}, got); df != "" {
				t.Errorf("Reseller -want +got\n%s", df)
			}
		})
	}
}

func TestListResellers(t *testing.T) {
	tests := []struct {
		name       string
		major      uint
		token      string
		rspStatus  int
		rspBody    []byte
		wantPath   string
		wantStatus int
		want       []Reseller
		wantErr    error
		wantToken  []string
	}{
		{
			name:      "normal",
			major:     7,
			token:     "psst!",
			rspStatus: http.StatusOK,
			rspBody:   []byte(`[{ "id": "tower", "name": "Donald" }, { "id": "sevnica", "name": "Melanija" }]`),
			wantPath:  "/v7/resellers",
			want:      []Reseller{{Id: "tower", Name: "Donald"}, {Id: "sevnica", Name: "Melanija"}},
			wantToken: []string{"Bearer psst!"},
		},
		{
			name:      "404",
			major:     12,
			token:     "***",
			rspStatus: http.StatusNotFound,
			rspBody:   []byte(`{"message": "covfefe"}`),
			wantPath:  "/v12/resellers",
			want:      []Reseller{},
			wantErr:   errResellerHttp,
			wantToken: []string{"Bearer ***"},
		},
		{
			name:      "bad-json",
			major:     12,
			token:     "***",
			rspStatus: http.StatusNotFound,
			rspBody:   []byte(`1+1=2`),
			wantPath:  "/v12/resellers",
			want:      []Reseller{},
			wantErr:   errResellerJson,
			wantToken: []string{"Bearer ***"},
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

				w.WriteHeader(test.rspStatus)
				w.Write([]byte(test.rspBody))
			}))
			defer server.Close()

			c := NewCryptlex(test.major, server.URL, test.token)
			got, err := c.ListResellers()
			if !errors.Is(err, test.wantErr) {
				t.Errorf("%s: error: got %v, want %s", test.name, err, test.wantErr)
			}
			if df := cmp.Diff(test.want, got); df != "" {
				t.Errorf("%s: response Reseller: -want +got\n%s", test.name, df)
			}
		})
	}
}

func TestCreateReseller(t *testing.T) {
	tests := []struct {
		name       string
		major      uint
		token      string
		reseller   Reseller
		rspStatus  int
		rspBody    []byte
		wantPath   string
		wantStatus int
		want       Reseller
		wantErr    error
		wantToken  []string
	}{
		{
			name:      "normal",
			major:     7,
			token:     "psst!",
			reseller:  Reseller{Name: "Donald"},
			rspStatus: http.StatusOK,
			rspBody:   []byte(`{ "id": "1", "name": "Donald" }`),
			wantPath:  "/v7/resellers",
			want:      Reseller{Id: "1", Name: "Donald"},
			wantToken: []string{"Bearer psst!"},
		},
		{
			name:      "ignore-id",
			major:     7,
			token:     "psst!",
			reseller:  Reseller{Id: "17", Name: "Donald"},
			rspStatus: http.StatusOK,
			rspBody:   []byte(`{ "id": "1", "name": "Donald" }`),
			wantPath:  "/v7/resellers",
			want:      Reseller{Id: "1", Name: "Donald"},
			wantToken: []string{"Bearer psst!"},
		},
		{
			name:      "404",
			major:     12,
			token:     "***",
			reseller:  Reseller{Name: "Melanija"},
			rspStatus: http.StatusNotFound,
			rspBody:   []byte(`{"message": "covfefe"}`),
			wantPath:  "/v12/resellers",
			want:      Reseller{},
			wantErr:   errResellerHttp,
			wantToken: []string{"Bearer ***"},
		},
		{
			name:      "bad-json",
			major:     12,
			token:     "***",
			reseller:  Reseller{Name: "Melanija"},
			rspStatus: http.StatusNotFound,
			rspBody:   []byte(`1+1=2`),
			wantPath:  "/v12/resellers",
			want:      Reseller{},
			wantErr:   errResellerJson,
			wantToken: []string{"Bearer ***"},
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

				w.WriteHeader(test.rspStatus)
				w.Write([]byte(test.rspBody))
			}))
			defer server.Close()

			c := NewCryptlex(test.major, server.URL, test.token)
			got, err := c.CreateReseller(test.reseller)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("%s: error: got %v, want %s", test.name, err, test.wantErr)
			}
			if df := cmp.Diff(test.want, got); df != "" {
				t.Errorf("%s: response Reseller: -want +got\n%s", test.name, df)
			}
		})
	}
}

func TestUpdateReseller(t *testing.T) {
	tests := []struct {
		name       string
		major      uint
		token      string
		reseller   Reseller
		rspStatus  int
		rspBody    []byte
		wantPath   string
		wantStatus int
		want       Reseller
		wantErr    error
		wantToken  []string
	}{
		{
			name:      "normal",
			major:     7,
			token:     "psst!",
			reseller:  Reseller{Id: "tower", Name: "Donald"},
			rspStatus: http.StatusOK,
			rspBody:   []byte(`{ "id": "tower", "name": "Donald" }`),
			wantPath:  "/v7/resellers/tower",
			want:      Reseller{Id: "tower", Name: "Donald"},
			wantToken: []string{"Bearer psst!"},
		},
		{
			name:      "404",
			major:     12,
			token:     "***",
			reseller:  Reseller{Id: "sevnica", Name: "Melanija"},
			rspStatus: http.StatusNotFound,
			rspBody:   []byte(`{"message": "covfefe"}`),
			wantPath:  "/v12/resellers/sevnica",
			want:      Reseller{},
			wantErr:   errResellerHttp,
			wantToken: []string{"Bearer ***"},
		},
		{
			name:      "bad-json",
			major:     12,
			token:     "***",
			reseller:  Reseller{Id: "sevnica", Name: "Melanija"},
			rspStatus: http.StatusNotFound,
			rspBody:   []byte(`1+1=2`),
			wantPath:  "/v12/resellers/sevnica",
			want:      Reseller{},
			wantErr:   errResellerJson,
			wantToken: []string{"Bearer ***"},
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

				w.WriteHeader(test.rspStatus)
				w.Write([]byte(test.rspBody))
			}))
			defer server.Close()

			c := NewCryptlex(test.major, server.URL, test.token)
			got, err := c.UpdateReseller(test.reseller)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("%s: error: got %v, want %s", test.name, err, test.wantErr)
			}
			if df := cmp.Diff(test.want, got); df != "" {
				t.Errorf("%s: response Reseller: -want +got\n%s", test.name, df)
			}
		})
	}
}

func TestRetreiveReseller(t *testing.T) {
	tests := []struct {
		name       string
		major      uint
		token      string
		id         string
		rspStatus  int
		rspBody    []byte
		wantPath   string
		wantStatus int
		want       Reseller
		wantErr    error
		wantToken  []string
	}{
		{
			name:      "normal",
			major:     7,
			token:     "psst!",
			id:        "tower",
			rspStatus: http.StatusOK,
			rspBody:   []byte(`{ "id": "tower", "name": "Donald" }`),
			wantPath:  "/v7/resellers/tower",
			want:      Reseller{Id: "tower", Name: "Donald"},
			wantToken: []string{"Bearer psst!"},
		},
		{
			name:      "404",
			major:     12,
			token:     "***",
			id:        "sevnica",
			rspStatus: http.StatusNotFound,
			rspBody:   []byte(`{"message": "covfefe"}`),
			wantPath:  "/v12/resellers/sevnica",
			want:      Reseller{},
			wantErr:   errResellerHttp,
			wantToken: []string{"Bearer ***"},
		},
		{
			name:      "bad-json",
			major:     12,
			token:     "***",
			id:        "sevnica",
			rspStatus: http.StatusNotFound,
			rspBody:   []byte(`1+1=2`),
			wantPath:  "/v12/resellers/sevnica",
			want:      Reseller{},
			wantErr:   errResellerJson,
			wantToken: []string{"Bearer ***"},
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

				w.WriteHeader(test.rspStatus)
				w.Write([]byte(test.rspBody))
			}))
			defer server.Close()

			c := NewCryptlex(test.major, server.URL, test.token)
			got, err := c.RetrieveReseller(test.id)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("%s: error: got %v, want %s", test.name, err, test.wantErr)
			}
			if df := cmp.Diff(test.want, got); df != "" {
				t.Errorf("%s: response Reseller: -want +got\n%s", test.name, df)
			}
		})
	}
}
