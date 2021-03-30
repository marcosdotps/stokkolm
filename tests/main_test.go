package tests

import (
	"os"
	"testing"

	"github.com/mpenate/stokkolm/dbconnect"
)

//TestMain initializes the db in order to make sure we can test
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	os.Setenv("MONGO_URL", "mongodb://172.17.0.2:27017/")
	os.Setenv("INVENTORY_PATH", "tests/test_schemas/inventory.json")
	os.Setenv("PRODUCTS_PATH", "tests/test_schemas/products.json")
	dbconnect.InitializeDB()
	return m.Run()
}
