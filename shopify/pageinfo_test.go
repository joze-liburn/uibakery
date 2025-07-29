package shopify

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPageFromGQL(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]any
		want    PageInfo
		wanterr error
	}{
		{
			name: "normal",
			data: map[string]any{
				"pageInfo": map[string]any{
					"hasNextPage": true,
					"endCursor":   "some-strange-identifier",
				},
			},
			want: PageInfo{HasNextPage: true, EndCursor: "some-strange-identifier"},
		},
		{
			name: "no record",
			data: map[string]any{},
			want: PageInfo{},
		},
		{
			name: "invalid flag",
			data: map[string]any{
				"pageInfo": map[string]any{
					"hasNextPage": 1,
					"endCursor":   "some-strange-identifier",
				},
			},
			wanterr: errPageInfo,
		},
		{
			name: "invalid cursor",
			data: map[string]any{
				"pageInfo": map[string]any{
					"hasNextPage": true,
					"endCursor":   1,
				},
			},
			wanterr: errPageInfo,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pi, err := pageFromGQL(test.data)
			if !errors.Is(err, test.wanterr) {
				t.Errorf("%s: got error %v, want %v", test.name, err, test.wanterr)
			}
			if df := cmp.Diff(test.want, pi); df != "" {
				t.Errorf("%s: -want +got\n%s", test.name, df)
			}
		})
	}
}
