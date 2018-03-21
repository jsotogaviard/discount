package application

import (
	"github.com/jsotogaviard/discount/application/model"
	"time"
	"github.com/jsotogaviard/discount/configuration"
	"log"
	"sync"
)


type Application struct {
	Channels interface {
		GetCart(threadIdx int) (chan *model.Cart)
		GetScan(threadIdx int) (chan *model.Scan )
		GetDelete(threadIdx int) (chan *model.Cart)
		GetPrice(threadIdx int) (chan *model.Price)
	}
	Parallelism int
	Delay int
}


type Channels struct {
	Cart        []chan *model.Cart
	Scan        []chan *model.Scan
	Delete      []chan *model.Cart
	Price       []chan *model.Price
	Parallelism int
}

func (t Channels) GetCart(threadIdx int) chan *model.Cart {
	return t.Cart[threadIdx]
}

func (t Channels) GetScan(threadIdx int) chan *model.Scan {
	return t.Scan[threadIdx]
}

func (t Channels) GetDelete(threadIdx int) chan *model.Cart {
	return t.Delete[threadIdx]
}

func (t Channels) GetPrice(threadIdx int) chan *model.Price {
	return t.Price[threadIdx]
}

// Create and start the application
func GetApplication() (app *Application) {

	config, err := configuration.GetConfig()
	if err != nil {
		log.Fatal("Cannot get config")
		return nil
	}
	carts := make([]chan *model.Cart, 0)
	scans := make([]chan *model.Scan, 0)
	deletes := make([]chan *model.Cart, 0)
	prices := make([]chan *model.Price, 0)

	for idx := 0; idx < config.Parallelism; idx++ {
		carts = append(carts, make(chan *model.Cart))
		scans = append(scans, make(chan *model.Scan))
		deletes = append(deletes, make(chan *model.Cart))
		prices = append(prices, make(chan *model.Price))
	}

	// Create channels
	channels := Channels{
		Cart:   carts,
		Scan:   scans,
		Delete: deletes,
		Price:  prices,
	}

	itemMap := createItems()

	// Create the application
	a := Application{channels, config.Parallelism, config.Delay}

	// Sync group to make sure
	var wg sync.WaitGroup
	
	// Start threads
	for idx := 0; idx < config.Parallelism; idx++ {
		wg.Add(1)
		go Start(a, idx, itemMap, &wg)
	}

	wg.Wait()

	// Return the application
	return &a
}

// Create the items
func createItems() map[string]model.Item {
	var OneForTwo = model.Rule{"1FOR2", 2, 0.5, "(quantity - quantity % ruleQuantity) / ruleQuantity * price + (quantity % ruleQuantity) * price"}
	var Voucher = model.Item{"VOUCHER", "", 5, OneForTwo}

	var MoreThanThree = model.Rule{"MORETHAN3", 3, 0.95, "quantity < ruleQuantity ? quantity * price : quantity * price * discount"}
	var TShirt = model.Item{"TSHIRT", "", 20, MoreThanThree}

	var UnitPrice = model.Rule{"UNITPRICE", 2, 0.5, "quantity * price"}
	var Mug = model.Item{"MUG", "", 7.5, UnitPrice}

	var itemMap = make(map[string]model.Item)

	itemMap[Voucher.Code] = Voucher
	itemMap[TShirt.Code] = TShirt
	itemMap[Mug.Code] = Mug

	return itemMap
}

func Start(application Application, threadIdx int, itemMap map[string]model.Item, wg *sync.WaitGroup) {

	CartChannel := application.Channels.GetCart(threadIdx)
	scanChannel := application.Channels.GetScan(threadIdx)
	deleteChannel := application.Channels.GetDelete(threadIdx)
	priceChannel := application.Channels.GetPrice(threadIdx)

	delay := application.Delay

	var state = make(map[string][]string)
	
	wg.Done()
	
	for {
		select {
		case cart := <-CartChannel:

			// Simulate sleep
			sleep(delay)

			// Do the cart
			doCheckout(state, cart)

		case scan := <-scanChannel:

			// Simulate sleep
			sleep(delay)

			// Do the scan
			doScan(state, itemMap, scan)

		case del := <-deleteChannel:

			// Simulate sleep
			sleep(delay)

			// Do the deletion
			doDelete(state, del)

		case price := <-priceChannel:

			// Simulate sleep
			sleep(delay)

			// Do the pricing
			doPrice(state, itemMap, price)
		}
	}
}

func sleep(delay int){
	if delay != 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}