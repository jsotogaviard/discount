package application

import (
	"github.com/jsotogaviard/discount/models"
	"github.com/jsotogaviard/discount/application/model"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jsotogaviard/discount/restapi/operations/checkout"
	"net/http"
	"github.com/Knetic/govaluate"
)

func (app *Application) Price(params checkout.PriceParams) middleware.Responder{

	// Create price object
	priceCom := &model.Price{
		Id:   params.ID,
		Answer: make(chan model.ResponseErrorFloat),
	}

	// Redirect to correct thread
	threadIdx := ComputeThreadIdx(priceCom.Id, app.Parallelism)

	// Send it and wait for answer
	totalChannel := app.Channels.GetPrice(threadIdx)
	totalChannel <- priceCom
	result := <-priceCom.Answer

	// Send http answer
	if result.Response == http.StatusOK {

		price := &models.Price{
			Price:result.Price,
			Currency: "EUR",
		}
		return checkout.NewPriceOK().WithPayload(price)
	} else if result.Response == http.StatusNotFound {

		payload := &models.NotFound{result.Error}
		return checkout.NewPriceNotFound().WithPayload(payload)
	} else {

		payload := &models.InternalServerError{result.Error}
		return checkout.NewDeleteInternalServerError().WithPayload(payload)
	}
}

func doPrice(state map[string][]string, itemMap map[string]model.Item, price *model.Price) {

	responseError := model.ResponseErrorFloat{}

	// Make sure it already exists
	if _, ok := state[price.Id]; ok {

		var items = state[price.Id]
		var price, err = computePrice(items, itemMap)

		if err != nil {

			responseError.Response = http.StatusInternalServerError
			responseError.Error = err.Error()
		} else {

			responseError.Response = http.StatusOK
			responseError.Price = price
		}

	} else {

		// Log an error
		responseError.Response = http.StatusNotFound
		responseError.Error = "id does not exist"
	}

	// Return the answer
	price.Answer <- responseError
}

func computePrice(items []string, itemMap map[string]model.Item) (float64, error) {

	var price float64 = 0

	itemCountMap := countItems(items)

	for itemKey, count := range itemCountMap {

		item := itemMap[itemKey]
		expression, err := govaluate.NewEvaluableExpression(item.Rule.Formula);
		if err != nil {
			return 0, err
		}

		parameters := make(map[string]interface{})
		parameters["quantity"] = count;
		parameters["price"] = item.Price;
		parameters["ruleQuantity"] = item.Rule.Quantity;
		parameters["discount"] = item.Rule.Discount;
		result, err := expression.Evaluate(parameters) ;
		if err != nil {
			return 0, err
		}

		price += result.(float64)
	}

	return price,nil
}

func countItems(items []string) map[string]int {
	itemCountMap := make(map[string]int)
	for _, item := range items {
		if _, ok := itemCountMap[item]; ok {
			itemCountMap[item]++
		} else {
			itemCountMap[item] = 1
		}
	}
	return itemCountMap
}
