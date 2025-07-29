package shopify

import "time"

type (
	// https://shopify.dev/docs/api/admin-graphql/2024-10/objects/Metafield
	Metafield struct {
		Id        string
		Key       string
		Type      string
		UpdatedAt time.Time
		Value     string
	}

	// https://shopify.dev/docs/api/admin-graphql/2024-10/connections/MetafieldConnection
	MetafieldConnection struct {
		Nodes    []Metafield
		PageInfo *PageInfo
	}
)
