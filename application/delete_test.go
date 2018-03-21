package application

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/jsotogaviard/discount/restapi/operations/checkout"
	"github.com/go-openapi/runtime"
	"log"
)


func TestDeleteKo(t *testing.T) {
	// Get application
	application := GetApplication()

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	params:= checkout.NewDeleteParams()
	params.ID = "id"
	params.HTTPRequest = req

	// Execute function
	response := application.Delete(params)

	// Serialize response
	responseRecord := httptest.NewRecorder()
	response.WriteResponse(responseRecord, runtime.JSONProducer())
	log.Print(responseRecord.Code)
	if responseRecord.Code == http.StatusOK {
		t.Error("status code")
	}
}

func TestDeleteOk(t *testing.T) {
	// Get application
	application := GetApplication()

	id:=DoCheckout(*application, t)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	params:= checkout.NewDeleteParams()
	params.ID = id
	params.HTTPRequest = req

	// Execute function
	response := application.Delete(params)

	// Serialize response
	responseRecord := httptest.NewRecorder()
	response.WriteResponse(responseRecord, runtime.JSONProducer())
	log.Print(responseRecord.Code)
	if responseRecord.Code != http.StatusOK {
		t.Error("status code")
	}

}