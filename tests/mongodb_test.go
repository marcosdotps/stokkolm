package tests

import (
	"log"
	"testing"

	"github.com/mpenate/stokkolm/dbconnect"
)

var ()

func TestGetInventedProductByName(t *testing.T) {
	_, err := dbconnect.GetProductByName("invented")
	if err == nil {
		t.Errorf("Error on GetProducByName for unexistent product. Expected: %s, Found: %s", "nil", err.Error())
	}
}

func TestGetExistentProductByName(t *testing.T) {
	_, err := dbconnect.GetProductByName("TestingChair")
	if err != nil {
		t.Errorf("Error on GetProducByName. Expected: %s, Found: %s", "nil", err.Error())
	}
}

func TestGetAllProducts(t *testing.T) {
	prod := dbconnect.GetAllProducts()
	if len(prod) != 2 {
		t.Errorf("Error on GetAllProducts. Expected: 2 products, Found: %d", len(prod))
	}
}

func TestRemoveUnexistentAmount(t *testing.T) {
	prod, err := dbconnect.GetProductByName("TestingChair")
	if err != nil {
		t.Errorf("Unexpected error during TestRemoveUnexistentAmount: %s", err.Error())
	}
	err = dbconnect.RemoveProductComponents(prod, 100)
	if err == nil {
		t.Errorf("Expecting error to occur but nothing happened while requesting more products than we have")
	}
}
func TestRemoveProductComponents(t *testing.T) {
	prod, err := dbconnect.GetProductByName("TestingChair")
	if err != nil {
		t.Errorf("Unexpected error during TestRemoveProductComponents: %s", err.Error())
	}
	err = dbconnect.RemoveProductComponents(prod, 1)
	if err != nil {
		t.Errorf("Expecting error to be nil but error raised selling products: %s", err.Error())
	}
	log.Println("Rearranging our warehouse after this test :)")
	defer t.Cleanup(func() { dbconnect.InitializeDB() })

}
