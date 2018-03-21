package client

import (
	"log"
	"github.com/jsotogaviard/discount/client/checkout"
	"testing"
	"time"
	"sync"
)

func TestWholeSequence(t *testing.T) {
	numConcurrentUsers := 50
	numLoops := 100
	var wg sync.WaitGroup
	start:= time.Now()
	for idx := 0; idx < numConcurrentUsers; idx++ {
		wg.Add(1)
		go singleSequence(t, numLoops, &wg)
	}
	wg.Wait()
	elapsed := time.Since(start)
	log.Print(elapsed)
}

func singleSequence(t *testing.T, num int, wg *sync.WaitGroup){
	defer wg.Done()
	for idx := 0; idx < num; idx++ {

		// Cart test
		tConfig := DefaultTransportConfig()
		tConfig.WithHost("127.0.0.1:8000")
		client := NewHTTPClientWithConfig(nil, tConfig)
		p := &checkout.CheckoutParams{}
		checkResp, checkErr := client.Checkout.Checkout(p.WithTimeout(2 * time.Second))
		if checkErr != nil {
			log.Fatal(checkErr)
		}

		// Scan test
		scanParams := &checkout.ScanParams{
			ID:    checkResp.Payload.ID,
			Items: []string{"MUG"},
		}
		_, scanErr := client.Checkout.Scan(scanParams.WithTimeout(2 * time.Second))
		if scanErr != nil {
			log.Fatal(scanErr)
		}

		// Price test
		priceParams := &checkout.PriceParams{
			ID: checkResp.Payload.ID,
		}
		priceResp, priceErr := client.Checkout.Price(priceParams.WithTimeout(2 * time.Second))
		if priceErr != nil {
			log.Fatal(priceErr)
		}
		if priceResp.Payload.Price != 7.5 {
			t.Error("Wrong price")
		}
		log.Println(priceResp.Payload.Price)

		// Delete test
		deleteParams := &checkout.DeleteParams{
			ID: checkResp.Payload.ID,
		}
		_, deleteErr := client.Checkout.Delete(deleteParams.WithTimeout(2 * time.Second))
		if deleteErr != nil {
			log.Fatal(deleteErr)
		}
	}
}