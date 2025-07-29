package main

import (
	"fmt"
	"os"

	"github.com/oussama4/gopify"
	"gitlab.com/joze-liburn/uibakery/lbqueue"
	"gitlab.com/joze-liburn/uibakery/shopify"
)

func TestLbqueue() {
	host := os.Getenv("PG_DEV_HOST")
	secret := os.Getenv("PG_DEV_PWD")
	db := lbqueue.LbDb{}
	err := db.Open("lb_ap_uibakery", secret, host, 5432, "lightburn")
	if err != nil {
		fmt.Printf("Opened with error %s\n", err)
		return
	}
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
}

func Test(limit int) {
	secret := os.Getenv("SHOPIFY_DEV_PWD")
	client := gopify.NewClient("lightburn-software-llc.myshopify.com", secret)
	ids, err := shopify.GetCompaniesIds(client, limit)
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
}

func main() {
	TestLbqueue()
}
