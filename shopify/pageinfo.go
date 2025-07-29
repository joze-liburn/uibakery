package shopify

import (
	"errors"
	"fmt"
)

var (
	errPageInfo = errors.New("invalid page info record")
)

type (
	// Shopify Admin GraphQL API
	// PageInfo
	// Returns information about pagination in a connection, in accordance with
	// the Relay specification. For more information, please read our GraphQL
	// Pagination Usage Guide.
	//
	// Fields
	// endCursor        String   The cursor corresponding to the last node in
	//                           edges.
	// hasNextPage      Boolean  Whether there are more pages to fetch following
	//                           the current page.
	// hasPreviousPage  Boolean  Whether there are any pages prior to the
	//                           current page.
	// startCursor      String   The cursor corresponding to the first node in
	//                           edges.
	// https://shopify.dev/docs/api/admin-graphql/latest/objects/PageInfo
	//
	// Only implement forward pagination - it is universally supported.
	PageInfo struct {
		HasNextPage bool
		EndCursor   string
	}
)

func pageFromGQL(r map[string]any) (PageInfo, error) {
	pi, ok := r["pageInfo"]
	if !ok {
		return PageInfo{}, nil
	}
	pit, ok := pi.(map[string]any)
	if !ok {
		return PageInfo{}, fmt.Errorf("%w: expected map, got %T (%v)", errPageInfo, pi, pi)
	}
	hnp, ok := pit["hasNextPage"].(bool)
	if !ok {
		return PageInfo{}, fmt.Errorf("%w: expected bool, got %T (%v)", errPageInfo, pit["hasNextPage"], pit["hasNextPage"])
	}
	ec, ok := pit["endCursor"].(string)
	if !ok {
		return PageInfo{}, fmt.Errorf("%w: expected string, got %T (%v)", errPageInfo, pit["endCursor"], pit["endCursor"])
	}
	return PageInfo{HasNextPage: hnp, EndCursor: ec}, nil
}
