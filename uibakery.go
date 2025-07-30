package main

import (
	"fmt"
	"os"

	"gitlab.com/joze-liburn/uibakery/lbqueue"
	"gitlab.com/joze-liburn/uibakery/shopify"
)

func TestLbqueue() {
	fmt.Println("Test LbDatabase.")
	host := os.Getenv("PG_DEV_HOST")
	secret := os.Getenv("PG_DEV_PWD")
	db := lbqueue.LbDb{}
	err := db.Open("lb_ap_uibakery", secret, host, 5432, "lightburn")
	if err != nil {
		fmt.Printf("Opened with error %s\n", err)
		return
	}
	fmt.Println("Database opened.")
	g, c, err := db.ClaimRecrds(10)
	if err != nil {
		fmt.Printf("ClaimRecords() retured an error %s\n", err)
		return
	}
	fmt.Printf("%d records claimed.\n", c)
	x, err := db.GetClaimedRecords(g)
	if err != nil {
		fmt.Printf("GetClaimedRecords() retured an error %s\n", err)
		return
	}
	fmt.Printf("%d records got\n", len(x))
	fmt.Println("Test LbDatabase ended.")
}

func TestShopify(limit int) {
	fmt.Println("Test Shopify.")
	secret := os.Getenv("SHOPIFY_DEV_PWD")
	client := shopify.New("lightburn-software-llc.myshopify.com", secret)
	fmt.Println("Client created")
	ids, err := client.GetCompaniesIds(limit)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Printf("Found %d ids\n", len(ids.Nodes))
	for i, id := range ids.Nodes {
		fmt.Printf("%d: %s\n", i, id.Id)
	}
	if ids.PageInfo.HasNextPage {
		fmt.Println("There's more")
	}
	fmt.Println("Test Shopify ended.")
}

func main() {
	TestLbqueue()
}
