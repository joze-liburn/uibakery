package cmd

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"gitlab.com/joze-liburn/uibakery/cryptlex"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name string
		cmd  *cobra.Command
		flag string
	}{
		{
			name: "normal",
			cmd:  &cobra.Command{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_ = &cobra.Command{}
		})
	}
}

func TestCreatedAt(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "page",
			args: []string{"--page", "1", "--limit", "100"},
			want: []string{"page=1", "limit=100"},
		},
		{
			name: "sort",
			args: []string{"--sort", "name"},
			want: []string{"sort=name"},
		},
		{
			name: "dates",
			args: []string{"--createdAt", "eq 2025-08-19"},
			want: []string{"createdAt=eq%202025-08-19T00:00:00Z"},
		},
		{
			name: "dates-formats",
			args: []string{"--createdAt", "eq 2025-08-19 9:30", "--updatedAt", "eq 2025-08-21 17:34:51"},
			want: []string{"createdAt=eq%202025-08-19T09:30:00Z", "updatedAt=eq%202025-08-21T17:34:51Z"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := &cobra.Command{}
			addFlags(c)
			c.Flags().Parse(test.args)
			opt, _ := parseFlags(c)
			co := &cryptlex.CryptlexOpts{}
			co.Apply(opt...)
			got := co.ToSearchParam()
			head, tail, ok := strings.Cut(got, "?")
			if !ok || len(head) > 0 {
				t.Fatalf("%s: parameters should start with '?': %s", test.name, got)
			}
			params := strings.Split(tail, "&")
			if df := cmp.Diff(params, test.want); df != "" {
				t.Errorf("%s: -want +got\n%s", test.name, df)
			}

		})
	}

}
