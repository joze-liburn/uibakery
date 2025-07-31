package shopify

type (
	// https://shopify.dev/docs/api/admin-graphql/2024-10/objects/Count
	Count struct {
		Count     int    `json:"count"`
		Precision string `json:"precision"`
	}

	// https://shopify.dev/docs/api/admin-graphql/2024-10/objects/MoneyV2
	MoneyV2 struct {
		Amount       float32 `json:"amount"`
		CurrencyCode string  `json:"currencyCode"`
	}
)
