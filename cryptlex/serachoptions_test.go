package cryptlex

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestWithInt(t *testing.T) {
	tests := []struct {
		name string
		opt  func(int) CallOptions
		arg  int
		want error
	}{
		{
			name: "page-above",
			opt:  WithPage,
			arg:  12,
		},
		{
			name: "page-below",
			opt:  WithPage,
			arg:  -7,
			want: errCallOption,
		},
		{
			name: "limit-above",
			opt:  WithPageSize,
			arg:  12,
		},
		{
			name: "limit-below",
			opt:  WithPageSize,
			arg:  -1,
			want: errCallOption,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			co := CryptlexOpts{}
			fcn := test.opt(test.arg)
			if fcn == nil {
				t.Fatalf("%s: nil function", test.name)
			}
			err := fcn(&co)
			if !errors.Is(err, test.want) {
				t.Errorf("%s: error want %v, got %v", test.name, err, test.want)
			}
		})
	}
}

func TestWithString(t *testing.T) {
	tests := []struct {
		name string
		opt  func(string) CallOptions
		arg  string
		want error
	}{
		{
			name: "sort",
			opt:  WithSort,
			arg:  "column",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			co := CryptlexOpts{}
			fcn := test.opt(test.arg)
			if fcn == nil {
				t.Fatalf("%s: nil function", test.name)
			}
			err := fcn(&co)
			if !errors.Is(err, test.want) {
				t.Errorf("%s: error want %v, got %v", test.name, err, test.want)
			}
		})
	}
}

func TestWithStringOp(t *testing.T) {
	tests := []struct {
		name string
		opt  func(string, string) CallOptions
		cmp  string
		arg  string
		want error
	}{
		{
			name: "id",
			opt:  WithId,
			cmp:  "eq",
			arg:  "some-id",
		},
		{
			name: "id-no-args",
			opt:  WithId,
			cmp:  "",
			arg:  "",
			want: errCallOption,
		},
		{
			name: "name",
			opt:  WithName,
			cmp:  "ne",
			arg:  "Janez",
		},
		{
			name: "name-no-args",
			opt:  WithName,
			cmp:  "",
			arg:  "",
			want: errCallOption,
		},
		{
			name: "email",
			opt:  WithEmail,
			cmp:  "in",
			arg:  "a,b,c",
		},
		{
			name: "email-no-args",
			opt:  WithEmail,
			cmp:  "",
			arg:  "",
			want: errCallOption,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			co := CryptlexOpts{}
			fcn := test.opt(test.cmp, test.arg)
			if fcn == nil {
				t.Fatalf("%s: nil function", test.name)
			}
			err := fcn(&co)
			if !errors.Is(err, test.want) {
				t.Errorf("%s: error want %v, got %v", test.name, err, test.want)
			}
		})
	}
}

func TestToSearchPattern(t *testing.T) {
	tests := []struct {
		name string
		opts []CallOptions
		want []string
	}{
		{
			name: "empty",
			want: []string{""},
		},
		{
			name: "error",
			opts: []CallOptions{WithPage(12), WithPageSize(-1)},
			want: []string{"page%2012"},
		},
		{
			name: "page",
			opts: []CallOptions{WithPage(12), WithPageSize(20)},
			want: []string{"page%2012", "limit%2020"},
		},
		{
			name: "set",
			opts: []CallOptions{WithId("in", strings.Join([]string{"1", "2", "4"}, ","))},
			want: []string{"id%20in%201,2,4"},
		},
		{
			name: "time",
			opts: []CallOptions{WithCreatedAt("gt", time.Date(2025, 8, 18, 12, 23, 45, 0, time.UTC))},
			want: []string{"createdAt%20gt%202025-08-18T12:23:45Z"},
		},
		{
			name: "all",
			opts: []CallOptions{
				WithPage(12),
				WithPageSize(20),
				WithName("eq", "Donald"),
				WithEmail("ne", "djt@wh.gov"),
				WithSearch("gte", "tariff"),
				WithId("in", strings.Join([]string{"1", "2", "4"}, ",")),
				WithSort("key"),
				WithCreatedAt("gt", time.Date(2025, 8, 18, 12, 23, 45, 0, time.UTC)),
				WithUpdatedAt("gt", time.Date(2025, 8, 18, 13, 03, 15, 0, time.UTC)),
			},
			want: []string{
				"page%2012",
				"limit%2020",
				"name%20eq%20Donald",
				"email%20ne%20djt@wh.gov",
				"search%20gte%20tariff",
				"id%20in%201,2,4",
				"sort%20key",
				"createdAt%20gt%202025-08-18T12:23:45Z",
				"updatedAt%20gt%202025-08-18T13:03:15Z",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			x := &CryptlexOpts{}
			for _, opt := range test.opts {
				opt(x)
			}
			got := x.ToSearchParam()
			if len(got) > 0 && got[0] != '?' {
				t.Errorf("%s: ToSearchParams(): got %q but doesn't start with '?'", test.name, got)
			}
			if len(got) > 0 && got[0] == '?' {
				got = got[1:]
			}
			gotparts := strings.Split(got, "&")
			if df := cmp.Diff(gotparts, test.want, cmpopts.SortSlices(func(a, b string) bool { return a < b })); df != "" {
				t.Errorf("%s: ToSearchParams(): -want +got\n%s", test.name, df)
			}
		})
	}
}
