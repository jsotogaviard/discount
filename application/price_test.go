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
	"strconv"
)


func TestPriceKo404(t *testing.T) {
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
	if responseRecord.Code != http.StatusNotFound {
		t.Error("status code")
	}
}

func TestPriceOk1(t *testing.T) {
	price := Price(t, []string{"MUG", "MUG"})
	if price.Price != 15 {
		t.Error("wrong in price expected 15 but is " + strconv.FormatFloat(price.Price, 'f', 6, 64))
	}
}

func TestPriceOk2(t *testing.T) {
	price := Price(t, []string{"VOUCHER", "TSHIRT", "MUG"})
	if price.Price != 32.5 {
		t.Error("wrong in price expected 32.5 but is " + strconv.FormatFloat(price.Price, 'f', 6, 64))
	}
}

func TestPriceOk3(t *testing.T) {
	price := Price(t, []string{"VOUCHER", "TSHIRT", "VOUCHER"})
	if price.Price != 25 {
		t.Error("wrong in price expected 35 but is " + strconv.FormatFloat(price.Price, 'f', 6, 64))
	}
}

func TestPriceOk4(t *testing.T) {
	price := Price(t, []string{"TSHIRT", "TSHIRT", "TSHIRT", "VOUCHER", "TSHIRT"})
	if price.Price != 81 {
		t.Error("wrong in price expected 81 but is " + strconv.FormatFloat(price.Price, 'f', 6, 64))
	}
}

func TestPriceOk5(t *testing.T) {
	price := Price(t, []string{"VOUCHER", "TSHIRT", "VOUCHER", "VOUCHER", "MUG", "TSHIRT", "TSHIRT"})
	if price.Price != 74.5 {
		t.Error("wrong in price expected 74.5 but is " + strconv.FormatFloat(price.Price, 'f', 6, 64))
	}
}

func Price(t *testing.T, items []string) *models.Price {
	// Get application
	application := GetApplication()

	id := DoCheckout(*application, t)

	ScanItems(*application, t, items, id)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	params:= checkout.NewPriceParams()
	params.ID = id
	params.HTTPRequest = req

	// Execute function
	response := application.Price(params)

	// Serialize response
	record := httptest.NewRecorder()
	response.WriteResponse(record, runtime.JSONProducer())
	log.Print(record.Code)
	if record.Code != http.StatusOK {
		t.Error("status code")
	}
	var price models.Price
	err := json.Unmarshal(record.Body.Bytes(), &price)
	if err != nil {
		t.Error("unmarshal")
	}
	log.Print(price.Price)
	return &price
}

func ScanItems(application Application, t *testing.T, items []string, id string){
	// Create request
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	params:= checkout.NewScanParams()
	params.ID = id
	params.Items = &models.ScanParamsBody{items}

	params.HTTPRequest = req

	// Execute function
	response := application.Scan(params)

	record := httptest.NewRecorder()
	response.WriteResponse(record, runtime.JSONProducer())
	if record.Code != http.StatusOK {
		t.Error("status code")
	}
}