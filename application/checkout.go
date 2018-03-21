package application

import (
	"github.com/jsotogaviard/discount/models"
	"github.com/jsotogaviard/discount/application/model"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jsotogaviard/discount/restapi/operations/checkout"
	"github.com/satori/go.uuid"
	"hash/fnv"
	"net/http"
)

func (app *Application) Checkout(params checkout.CartParams) middleware.Responder{

	// Generate uuid
	uid := uuid.NewV4()

	// Create checkout object
	cartCom := &model.Cart{
		Id:   uid.String(),
		Answer: make(chan model.ResponseError),
	}

	// Redirect to correct thread
	threadIdx := ComputeThreadIdx(cartCom.Id, app.Parallelism)

	// Send it and wait for answer
	cartChannel := app.Channels.GetCart(threadIdx)
	cartChannel <- cartCom
	result := <-cartCom.Answer

	// Send http answer
	if result.Success == http.StatusOK {

		payload := &models.Cart{cartCom.Id}
		return checkout.NewCartOK().WithPayload(payload)
	} else {

		payload := &models.InternalServerError{result.Error}
		return checkout.NewCartInternalServerError().WithPayload(payload)
	}
}

func doCheckout(state map[string][]string, checkout *model.Cart) {

	responseError := model.ResponseError{}

	// Make sure it does not already exist
	if _, ok := state[checkout.Id]; ok {

		// Log an error
		responseError.Success = http.StatusBadRequest
		responseError.Error = "id already used. Try again"

	} else {

		// Add the id
		state[checkout.Id] = make([]string, 0)
		responseError.Success = http.StatusOK
	}

	// Return the answer
	checkout.Answer <- responseError
}

// Compute thread index from uuid
func ComputeThreadIdx(id string, parallelism int) int{
	h := fnv.New32a()
	h.Write([]byte(id))
	return int(h.Sum32()) % parallelism
}
