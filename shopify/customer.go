package shopify

import "time"

type (
	// https://shopify.dev/docs/api/admin-graphql/2024-10/objects/Customer
	Customer struct {
		//AmountSpent             *MoneyV2
		CanDelete               bool      `json:"canDelete"`
		CreatedAt               time.Time `json:"createdAt"`
		DataSaleOptOut          bool      `json:"dataSaleOptOut"`
		DisplayName             string    `json:"displayName"`
		Email                   string    `json:"email"`
		FirstName               string    `json:"firstName"`
		HasTimelineComment      bool      `json:"hasTimelineComment"`
		Id                      string    `json:"id"`
		LastName                string    `json:"lastName"`
		LegacyResourceId        string    `json:"legacyResourceId"`
		LifetimeDuration        string    `json:"lifetimeDuration"`
		Locale                  string    `json:"locale"`
		MultipassIdentifier     string    `json:"multipassIdentifier"`
		Note                    string    `json:"note"`
		NumberOfOrders          string    `json:"numberOfOrders"`
		Phone                   string    `json:"phone"`
		ProductSubscriberStatus string    `json:"productSubscriberStatus"`
		State                   string    `json:"state"`
		Tags                    []string  `json:"tags"`
		UpdatedAt               time.Time `json:"updatedAt"`
		ValidEmailAddress       bool      `json:"validEmailAddress"`
		VerifiedEmail           bool      `json:"verifiedEmail"`
	}

	// https://shopify.dev/docs/api/admin-graphql/unstable/objects/CustomerEmailAddress
	CustomerEmailAddress struct {
		EmailAddress        string `json:"emailAddress"`
		MarketingOptInLevel int    `json:"marketingOptInLevel"`
	}
)
