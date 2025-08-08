package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/joze-liburn/uibakery/cmd"
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
	g, c, err := db.ClaimRecords(10)
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
	after := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	ids := client.StreamCompaniesIds(limit, 27000, &after)

	count := 0
	for nodeerr := range ids {
		if nodeerr.GetError() != nil {
			fmt.Printf("ERROR: %s\n", nodeerr.GetError())
			return
		}
		fmt.Printf("%3d: %s\n", count, nodeerr.(shopify.CompanyError).Company.Id)
		count++
	}
	fmt.Println("Test Shopify ended.")
}

func main() {
	cmd.Execute()
}
