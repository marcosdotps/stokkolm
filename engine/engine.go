package engine

import (	

	"github.com/mpenate/stokkolm/dbconnect"
	"github.com/mpenate/stokkolm/schemas"	
)

//GetProductMaxItems implements the logic to retrieve max assembled products of a kind with current stock
func GetProductMaxItems(prod schemas.Product) int {	
	matrix := make(map[int]int)			
	for i := 0; i < len(prod.ContainArticles); i++ {		
		matrix[prod.ContainArticles[i].ArtID] = dbconnect.GetStock(prod.ContainArticles[i].ArtID) / prod.ContainArticles[i].AmountOf		
	}			
	

	//Max possible items, otherwise it will be an integer overflow
	smallest := 4294967295

	
	for _, num := range matrix {		
		if num < smallest {			
			smallest = num
		}
	}
		
	return smallest
}

//RetrieveFullStock returs all stock for all the products
func RetrieveFullStock() map[string]int{	
	prods := dbconnect.GetAllProducts()
	stockMap := make(map[string]int)
	for _, prod := range prods {
		stockMap[prod.Name] = GetProductMaxItems(prod)
	}		
	return stockMap
}

//GetProductById returs a product filtered by id
func GetMaxProductByName(name string) int{
	product := dbconnect.GetProductByName(name)
	return GetProductMaxItems(product)
}


//RemoveStock ensures stock is deleted
func RemoveStock(productName string, amount int) {
	prod := dbconnect.GetProductByName(productName)
	dbconnect.RemoveProductComponents(prod, amount)
}