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

func TestWithFlagStr(t *testing.T) {
	tests := []struct {
		name string
		opt  func(string) (CallOptions, error)
		flag string
		want CryptlexOpts
		err  error
	}{
		{
			name: "supported",
			opt:  WithIdFlag,
			flag: "in test",
			want: CryptlexOpts{Id: OptCmpString{Cmp: "in", Operand: "test"}},
		},
		{
			name: "unsupported",
			opt:  WithNameFlag,
			flag: "vk test",
			err:  errCallOption,
		},
		{
			name: "single-word",
			opt:  WithEmailFlag,
			flag: "test",
			err:  errCallOption,
		},
		{
			name: "multiple-word",
			opt:  WithNameFlag,
			flag: "ne to be or not to be",
			want: CryptlexOpts{Name: OptCmpString{Cmp: "ne", Operand: "to be or not to be"}},
		},
		{
			name: "empty",
			opt:  WithEmailFlag,
			flag: "test",
			err:  errCallOption,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fnc, err := test.opt(test.flag)
			if !errors.Is(err, test.err) {
				t.Errorf("%s: parsing error, got %v, want %v", test.name, err, test.err)
			}
			if fnc == nil {
				return
			}
			co := CryptlexOpts{}
			err = fnc(&co)
			if !errors.Is(err, test.err) {
				t.Errorf("%s: WithFlag error, got %v, want %v", test.name, err, test.err)
			}
			if df := cmp.Diff(co, test.want); df != "" {
				t.Errorf("%s: -want +got\n%s", test.name, df)
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
			want: []string{"page=12"},
		},
		{
			name: "page",
			opts: []CallOptions{WithPage(12), WithPageSize(20)},
			want: []string{"page=12", "limit=20"},
		},
		{
			name: "set",
			opts: []CallOptions{WithId("in", strings.Join([]string{"1", "2", "4"}, ","))},
			want: []string{"id=in%201,2,4"},
		},
		{
			name: "time",
			opts: []CallOptions{WithCreatedAt("gt", time.Date(2025, 8, 18, 12, 23, 45, 0, time.UTC))},
			want: []string{"createdAt=gt%202025-08-18T12:23:45Z"},
		},
		{
			name: "all",
			opts: []CallOptions{
				WithPage(12),
				WithPageSize(20),
				WithName("eq", "Donald"),
				WithEmail("ne", "djt@wh.gov"),
				WithSearch("current tariff"),
				WithId("in", strings.Join([]string{"1", "2", "4"}, ",")),
				WithSort("key"),
				WithCreatedAt("gt", time.Date(2025, 8, 18, 12, 23, 45, 0, time.UTC)),
				WithUpdatedAt("gt", time.Date(2025, 8, 18, 13, 03, 15, 0, time.UTC)),
			},
			want: []string{
				"page=12",
				"limit=20",
				"name=eq%20Donald",
				"email=ne%20djt@wh.gov",
				"current tariff",
				"id=in%201,2,4",
				"sort=key",
				"createdAt=gt%202025-08-18T12:23:45Z",
				"updatedAt=gt%202025-08-18T13:03:15Z",
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
