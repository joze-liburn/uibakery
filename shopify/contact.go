package shopify

import "time"

type (
	// https://shopify.dev/docs/api/admin-graphql/2024-10/objects/CompanyContact
	Contact struct {
		CreatedAt        time.Time `json:"createdAt"`
		Id               string    `json:"id"`
		IsMainContact    bool      `json:"isMainContact"`
		LifetimeDuration string    `json:"lifetimeDuration"`
		Locale           string    `json:"locale"`
		Title            string    `json:"title"`
		UpdatedAt        time.Time `json:"updatedAt"`
		Customer         Customer  `json:"customer"`
	}

	Contacts struct {
		Nodes    []Contact `json:"nodes"`
		PageInfo *PageInfo `json:"pageInfo"`
	}
)
