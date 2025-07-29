package shopify

type (
	// https://shopify.dev/docs/api/admin-graphql/2024-10/objects/Count
	Count struct {
		Count     int
		Precision string
	}

	// https://shopify.dev/docs/api/admin-graphql/2024-10/objects/MoneyV2
	MoneyV2 struct {
		Amount       float32
		CurrencyCode string
	}
)
