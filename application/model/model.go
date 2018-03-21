package model

import "github.com/jsotogaviard/discount/models"

// The item
type Item struct {
	Code        string
	Description string
	Price      	float64
	Rule      	Rule
}

// The rule associated to an item
type Rule struct {
	Code        string
	Quantity 	int32
	Discount    float64
	Formula 	string
}

// Response with error
type ResponseError struct{
	Success int
	Error	string
}

//TODO Extends success error
type ResponseErrorFloat struct{
	Response int
	Error    string
	Price    float64
}

// Cart id
type Cart struct {
	Id   	string
	Answer  chan ResponseError
}

// Scan entity
type Scan struct {
	Id     string
	Items  models.ScanParamsBody
	Answer chan ResponseError
}

// Price entity
type Price struct {
	Id     string
	Answer chan ResponseErrorFloat
}


