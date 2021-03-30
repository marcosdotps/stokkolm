package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/mpenate/stokkolm/dbconnect"
	"github.com/mpenate/stokkolm/engine"
)

func TestGetProductMaxItems(t *testing.T) {
	prod, err := dbconnect.GetProductByName("TestingChair")
	if err != nil {
		t.Errorf("Unexpected error occurred executing TestGetProductMaxItems: %s", err.Error())
	}
	expected := 2
	result := engine.GetProductMaxItems(prod)

	if result != expected {
		t.Errorf("Error executing TestGetProductMaxItems. Expected: %d Result: %d", expected, result)
	}
}

func TestRetrieveFullStockHasAlwaysAllItems(t *testing.T) {
	stockMap := engine.RetrieveFullStock()
	expected := 2
	result := len(stockMap)
	if result != expected {
		t.Errorf("No stock info has been found for any of the products. Expected: %d, Result: %d", expected, result)
	}

}

func TestTooManyProductsRemoveStock(t *testing.T) {
	err := engine.RemoveStock("TestingChair", 200)
	if err == nil {
		t.Errorf("Something must be wrong in TestTooManyProductsRemoveStock while requesting 200 TestingChairs but no error raised!")
	}
}

func TestInventedProductRemoveStock(t *testing.T) {
	err := engine.RemoveStock("InventedChair", 200)
	if err == nil {
		t.Errorf("Something must be wrong in TestInventedProductRemoveStock while requesting InventedChairs but no error raised!")
	}
}

func TestRemoveStock(t *testing.T) {
	err := engine.RemoveStock("TestingChair", 2)
	if err != nil {
		t.Errorf("Something happened running TestRemoveStock while removing stock: %s", err.Error())
	}
}

func TestRetrieveStockAfterBuyAll(t *testing.T) {
	//This test relies on previous TestRemoveStock successful execution.
	stockMap := engine.RetrieveFullStock()
	expected := 2
	expectedItems := 0
	result := len(stockMap)
	resultChairItems := stockMap["TestingChair"]
	resultTableItems := stockMap["TestingTable"]
	if result != expected {
		t.Errorf("No stock has been found for any of the products. Expected: %d, Result: %d", expected, result)
	}
	if resultChairItems != expectedItems {
		t.Errorf("TestingChair stock has not been updated. Expected: %d, Result: %d", expectedItems, resultChairItems)
	}
	if resultTableItems != expectedItems {
		t.Errorf("TestingTable stock has not been updated. Expected: %d, Result: %d", expectedItems, resultTableItems)
	}
	log.Println("Rearranging our warehouse after this test :)")
	defer t.Cleanup(func() { dbconnect.InitializeDB() })
}

func TestGetMaxExistentProductByName(t *testing.T) {
	expected := 2
	result, err := engine.GetMaxProductByName("TestingChair")
	if err != nil {
		t.Errorf("Something went wrong duringTestGetMaxExistentProductByName execution: %s", err)
	}
	if result != expected {
		t.Errorf("TestGetMaxExistentProductByName failed. Expected: %d, Result: %d", expected, result)
	}
}

func TestGetMaxInventedProductByName(t *testing.T) {
	result, err := engine.GetMaxProductByName("NotATestingChair")
	fmt.Println(result)
	if err == nil {
		t.Errorf("TestGetMaxInventedProductByName failed as not error has been thrown. Result: %d", result)
	}
	if result != 0 {
		t.Errorf("TestGetMaxInventedProductByName failed as we expected 0 items retrieved. Result: %d", result)
	}
}
