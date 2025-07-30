package shopify

import "github.com/oussama4/gopify"

type ShopifyOp struct {
	client *gopify.Client
}

func New(host string, secret string) ShopifyOp {
	return ShopifyOp{client: gopify.NewClient(host, secret)}
}
