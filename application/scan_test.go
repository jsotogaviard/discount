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


func TestScanKo(t *testing.T) {
	// Get application
	application := GetApplication()

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	params:= checkout.NewScanParams()
	params.ID = "id"
	items := []string{"MUG", "MUG"}
	params.Items = &models.ScanParamsBody{items}
	params.HTTPRequest = req

	// Execute function
	response := application.Scan(params)

	// Serialize response
	responseRecord := httptest.NewRecorder()
	response.WriteResponse(responseRecord, runtime.JSONProducer())
	log.Print(responseRecord.Code)
	if responseRecord.Code == http.StatusOK {
		t.Error("status code")
	}
}

func TestScanOk(t *testing.T) {
	// Get application
	application := GetApplication()

	id:=DoCheckout(*application, t)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	params:= checkout.NewScanParams()
	params.ID = id
	items := []string{"MUG", "MUG"}
	params.Items = &models.ScanParamsBody{items}
	params.HTTPRequest = req

	// Execute function
	response := application.Scan(params)

	// Serialize response
	responseRecord := httptest.NewRecorder()
	response.WriteResponse(responseRecord, runtime.JSONProducer())
	log.Print(responseRecord.Code)
	if responseRecord.Code != http.StatusOK {
		t.Error("status code")
	}

}


func TestScanKoItem(t *testing.T) {
	// Get application
	application := GetApplication()

	id:=DoCheckout(*application, t)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	params:= checkout.NewScanParams()
	params.ID = id
	items := []string{"MAG", "MUG"}
	params.Items = &models.ScanParamsBody{items}
	params.HTTPRequest = req

	// Execute function
	response := application.Scan(params)

	// Serialize response
	responseRecord := httptest.NewRecorder()
	response.WriteResponse(responseRecord, runtime.JSONProducer())
	log.Print(responseRecord.Code)
	if responseRecord.Code != http.StatusBadRequest {
		t.Error("status code")
	}

}

func DoCheckout(application Application, t *testing.T) string{
	// Create request
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	params:= checkout.NewCartParams()
	params.HTTPRequest = req

	// Execute function
	response:=application.Checkout(params)

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
	return check.ID
}