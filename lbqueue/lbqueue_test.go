package lbqueue

import (
	"database/sql"
	"testing"
	"time"
)

func TestGetTime(t *testing.T) {
	now := time.Now()
	zro := time.Time{}

	tests := []struct {
		name string
		data sql.NullTime
		want *time.Time
	}{
		{name: "time", data: sql.NullTime{Time: now, Valid: true}, want: &now},
		{name: "wierd", data: sql.NullTime{Valid: true}, want: &zro},
		{name: "benign", data: sql.NullTime{Time: now}, want: nil},
		{name: "null", data: sql.NullTime{}, want: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetTime(test.data)
			switch {
			case got == nil && test.want != nil:
				t.Errorf("%s: got nil, want %s", test.name, test.want.Format("2006-01-02 15:04:05"))
			case got != nil && test.want == nil:
				t.Errorf("%s: got %s, want nil", test.name, got.Format("2006-01-02 15:04:05"))
			case (got != nil) && !(*got).Equal(*test.want):
				t.Errorf("%s: got %s, want %s", test.name, got.Format("2006-01-02 15:04:05"), test.want.Format("2006-01-02 15:04:05"))
			}
		})
	}
}

func TestGetString(t *testing.T) {
	tests := []struct {
		name string
		data sql.NullString
		want string
	}{
		{name: "string", data: sql.NullString{String: "data", Valid: true}, want: "data"},
		{name: "wierd", data: sql.NullString{Valid: true}, want: ""},
		{name: "benign", data: sql.NullString{String: "hidden"}, want: ""},
		{name: "null", data: sql.NullString{}, want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetString(test.data)
			if got != test.want {
				t.Errorf("%s: got %s, want %s", test.name, got, test.want)
			}
		})
	}
}
