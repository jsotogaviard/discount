package application

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"github.com/jsotogaviard/discount/restapi/operations/checkout"
	"github.com/go-openapi/runtime"
	"github.com/jsotogaviard/discount/models"
	"log"
)


func TestCheckoutOk(t *testing.T) {
	// Get application
	application := GetApplication()

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	params:= checkout.NewCartParams()
	params.HTTPRequest = req

	// Execute function
	response := application.Checkout(params)

	// Serialize response
	record := httptest.NewRecorder()
	response.WriteResponse(record, runtime.JSONProducer())
	if record.Code != http.StatusOK {
		t.Error("status code")
	}
	var check models.Cart
	err := json.Unmarshal(record.Body.Bytes(), &check)
	if err != nil {
		t.Error("unmarshal")
	}
	log.Print(check.ID)
}