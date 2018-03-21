package application

import (
	"github.com/jsotogaviard/discount/models"
	"github.com/jsotogaviard/discount/application/model"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jsotogaviard/discount/restapi/operations/checkout"
	"net/http"
)

func (app *Application) Delete(params checkout.DeleteParams) middleware.Responder{

	// Create delete object
	deleteCom := &model.Cart{
		Id:   params.ID,
		Answer: make(chan model.ResponseError),
	}

	/// Redirect to correct thread
	threadIdx := ComputeThreadIdx(deleteCom.Id, app.Parallelism)

	// Send it and wait for answer
	deleteChannel := app.Channels.GetDelete(threadIdx)
	deleteChannel <- deleteCom
	result := <-deleteCom.Answer

	// Send http answer
	if result.Success == http.StatusOK {

		return checkout.NewDeleteOK()
	} else if result.Success == http.StatusNotFound {

		payload := &models.NotFound{result.Error}
		return checkout.NewDeleteNotFound().WithPayload(payload)
	} else {

		payload := &models.InternalServerError{result.Error}
		return checkout.NewDeleteInternalServerError().WithPayload(payload)
	}
}

func doDelete(state map[string][]string, del *model.Cart) {
	responseError := model.ResponseError{}

	// Make sure it already exists
	if _, ok := state[del.Id]; ok {

		// Delete checkout
		delete(state, del.Id)
		responseError.Success = http.StatusOK
	} else {

		// Log an error
		responseError.Success = http.StatusNotFound
		responseError.Error = "id does not exist"
	}

	// Return the response
	del.Answer <- responseError
}
