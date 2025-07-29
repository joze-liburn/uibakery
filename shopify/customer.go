package shopify

import "time"

type (
	// https://shopify.dev/docs/api/admin-graphql/2024-10/objects/Customer
	Customer struct {
		//AmountSpent             *MoneyV2
		CanDelete               bool
		CreatedAt               time.Time
		DataSaleOptOut          bool
		DisplayName             string
		Email                   string
		FirstName               string
		HasTimelineComment      bool
		Id                      string
		LastName                string
		LegacyResourceId        string
		LifetimeDuration        string
		Locale                  string
		MultipassIdentifier     string
		Note                    string
		NumberOfOrders          string
		Phone                   string
		ProductSubscriberStatus string
		State                   string
		Tags                    []string
		UpdatedAt               time.Time
		ValidEmailAddress       bool
		VerifiedEmail           bool
	}

	// https://shopify.dev/docs/api/admin-graphql/unstable/objects/CustomerEmailAddress
	CustomerEmailAddress struct {
		EmailAddress        string
		MarketingOptInLevel int
	}
)
