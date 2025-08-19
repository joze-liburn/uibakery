package cryptlex

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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
