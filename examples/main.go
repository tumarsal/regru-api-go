package main

import (
	"context"
	"fmt"
	"log"
	"time"

	regru "github.com/example/regru-api-go"
)

func main() {
	ctx := context.Background()

	// === Example 1: Test API availability (no auth required) ===
	fmt.Println("=== Testing API availability ===")
	testClient := regru.NewClient("test", "test")
	if err := testClient.Nop(ctx); err != nil {
		log.Printf("Test nop failed: %v", err)
	} else {
		fmt.Println("API is accessible with test credentials")
	}

	// === Example 2: Initialize client with real credentials ===
	// Replace with your actual REG.RU username and password
	// Or use signature authentication for better security
	username := "your_username"
	password := "your_password"

	client := regru.NewClient(username, password,
		regru.WithOutputFormat(regru.OutputJSON),
	)

	// === Example 3: Check domain availability ===
	fmt.Println("\n=== Checking domain availability ===")
	domainService := client.Domain()
	result, err := domainService.Check(ctx, "example.ru", false)
	if err != nil {
		log.Printf("Domain check error: %v", err)
	} else {
		fmt.Printf("Domain: %s, Status: %s\n", result.DomainName, result.Result)
		if result.IsPremium {
			fmt.Printf("Premium domain price: %.2f\n", result.Price)
		}
	}

	// === Example 4: Check multiple domains ===
	fmt.Println("\n=== Checking multiple domains ===")
	domains := []string{"example.ru", "example.com", "example.net"}
	results, err := domainService.CheckBatch(ctx, domains, false)
	if err != nil {
		log.Printf("Batch check error: %v", err)
	} else {
		for _, r := range results {
			status := r.Result
			if r.ErrorCode != "" {
				status = fmt.Sprintf("%s (%s)", r.Result, r.ErrorCode)
			}
			fmt.Printf("  %s: %s\n", r.DomainName, status)
		}
	}

	// === Example 5: Get domain prices ===
	fmt.Println("\n=== Getting domain prices ===")
	prices, err := domainService.GetPrices(ctx, regru.CurrencyRUR)
	if err != nil {
		log.Printf("Get prices error: %v", err)
	} else {
		count := 0
		for zone, price := range prices {
			if count >= 5 {
				fmt.Println("  ...")
				break
			}
			fmt.Printf("  .%s: register=%.2f, renew=%.2f\n",
				zone, price.PriceRegister, price.PriceRenew)
			count++
		}
	}

	// === Example 6: Get user balance ===
	fmt.Println("\n=== Getting user balance ===")
	userService := client.User()
	balance, err := userService.GetBalance(ctx)
	if err != nil {
		log.Printf("Get balance error: %v", err)
	} else {
		fmt.Printf("Balance: %.2f %s\n", balance.Balance, balance.Currency)
	}

	// === Example 7: Get service list ===
	fmt.Println("\n=== Getting service list ===")
	serviceHandler := client.Service()
	services, err := serviceHandler.GetList(ctx, "")
	if err != nil {
		log.Printf("Get service list error: %v", err)
	} else {
		fmt.Printf("Found %d services\n", len(services))
		for i, svc := range services {
			if i >= 5 {
				fmt.Println("  ...")
				break
			}
			fmt.Printf("  %s (ID: %s, Type: %s)\n", svc.DomainName, svc.ServiceID, svc.ServiceType)
		}
	}

	// === Example 8: DNS zone management ===
	fmt.Println("\n=== DNS zone management ===")
	zoneService := client.Zone()

	if err := zoneService.Add_Alias(ctx, "example.ru", "www", "1.2.3.4"); err != nil {
		log.Printf("Add A record error: %v", err)
	} else {
		fmt.Println("Added A record")
	}

	if err := zoneService.Add_MX(ctx, "example.ru", "@", "mail.example.ru", "10"); err != nil {
		log.Printf("Add MX record error: %v", err)
	} else {
		fmt.Println("Added MX record")
	}

	// === Example 9: Folder management ===
	fmt.Println("\n=== Folder management ===")
	folderService := client.Folder()
	folders, err := folderService.GetList(ctx)
	if err != nil {
		log.Printf("Get folders error: %v", err)
	} else {
		fmt.Printf("Found %d folders\n", len(folders))
		for _, f := range folders {
			fmt.Printf("  %s (ID: %d)\n", f.FolderName, f.FolderID)
		}
	}

	// === Example 10: Get bills ===
	fmt.Println("\n=== Getting bills ===")
	billService := client.Bill()
	now := time.Now()
	bills, err := billService.GetForPeriod(ctx, regru.GetForPeriodRequest{
		StartDate: now.AddDate(0, -3, 0).Format("2006-01-02"),
		EndDate:   now.Format("2006-01-02"),
	})
	if err != nil {
		log.Printf("Get bills error: %v", err)
	} else {
		fmt.Printf("Found %d bills\n", len(bills))
		for i, bill := range bills {
			if i >= 5 {
				fmt.Println("  ...")
				break
			}
			fmt.Printf("  Bill %s: %.2f %s (%s)\n",
				bill.BillID, bill.TotalPayment, bill.Currency, bill.PayStatus)
		}
	}

	// === Example 11: Using signature authentication ===
	fmt.Println("\n=== Signature authentication (advanced) ===")
	// To use signature auth, you need an SSL certificate uploaded to REG.RU API settings
	// and its corresponding private key file.
	//
	// privateKeyData, err := os.ReadFile("/path/to/private.key")
	// if err != nil {
	//     log.Fatal(err)
	// }
	// sigClient, err := regru.ClientWithSignature("username", privateKeyData)
	// if err != nil {
	//     log.Fatal(err)
	// }
	fmt.Println("Signature authentication is available via regru.ClientWithSignature()")

	fmt.Println("\n=== All examples completed ===")
}
