package application

import (
	"github.com/jsotogaviard/discount/models"
	"github.com/jsotogaviard/discount/application/model"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jsotogaviard/discount/restapi/operations/checkout"
	"net/http"
)

func (app *Application) Scan(params checkout.ScanParams) middleware.Responder{

	// Create scan object
	scanCom := &model.Scan{
		Id:   params.ID,
		Items : *params.Items,
		Answer: make(chan model.ResponseError),
	}

	// Redirect to correct thread
	threadIdx := ComputeThreadIdx(scanCom.Id, app.Parallelism)

	// Send it and wait for answer
	scanChannel := app.Channels.GetScan(threadIdx)
	scanChannel <- scanCom
	result := <-scanCom.Answer

	// Send http answer
	if result.Success == http.StatusOK {

		return checkout.NewScanOK()
	} else if result.Success == http.StatusNotFound {

		payload := &models.NotFound{result.Error}
		return checkout.NewScanNotFound().WithPayload(payload)
	} else if result.Success == http.StatusBadRequest {

		payload := &models.BadRequest{result.Error}
		return checkout.NewScanBadRequest().WithPayload(payload)
	} else {

		payload := &models.InternalServerError{result.Error}
		return checkout.NewScanInternalServerError().WithPayload(payload)

	}
}

func doScan(state map[string][]string, itemMap map[string]model.Item,  scan *model.Scan) {
	responseError := model.ResponseError{}

	// Make sure id exists
	if _, ok := state[scan.Id]; ok {

		// Make sure items are known
		knownItems := true
		for _, item := range scan.Items.Items {
			if _, ok := itemMap[item]; !ok {
				knownItems = false
				break
			}
		}

		if !knownItems {

			// Log an error
			responseError.Success = http.StatusBadRequest
			responseError.Error = "item does not exist"
		} else {

			// Add the items
			state[scan.Id] = append(state[scan.Id], scan.Items.Items...)
			responseError.Success = http.StatusOK
		}

	} else {

		// Log an error
		responseError.Success = http.StatusNotFound
		responseError.Error = "id does not exist"
	}

	// Return the response
	scan.Answer <- responseError
}
