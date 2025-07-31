package shopify

import "time"

type (
	// https://shopify.dev/docs/api/admin-graphql/2024-10/objects/Metafield
	Metafield struct {
		Id        string    `json:"id"`
		Key       string    `json:"key"`
		Type      string    `json:"type"`
		UpdatedAt time.Time `json:"updatedAt"`
		Value     string    `json:"value"`
	}

	// https://shopify.dev/docs/api/admin-graphql/2024-10/connections/MetafieldConnection
	MetafieldConnection struct {
		Nodes    []Metafield `json:"nodes"`
		PageInfo *PageInfo   `json:"pageInfo"`
	}
)
