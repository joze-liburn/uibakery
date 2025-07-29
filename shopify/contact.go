package shopify

import "time"

type (
	// https://shopify.dev/docs/api/admin-graphql/2024-10/objects/CompanyContact
	Contact struct {
		CreatedAt        time.Time
		Id               string
		IsMainContact    bool
		LifetimeDuration string
		Locale           string
		Title            string
		UpdatedAt        time.Time
		Customer         Customer
	}

	Contacts struct {
		Nodes    []Contact
		PageInfo *PageInfo
	}
)
