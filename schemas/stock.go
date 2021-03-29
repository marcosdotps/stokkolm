package schemas

//Inventory represent the stock
type Inventory struct {
	Articles []Article `json:"inventory"`
}

//Article definition
type Article struct {
	ArtID int    `json:"art_id,string"`
	Name  string `json:"name"`
	Stock int    `json:"stock,string"`
}

//Products is the list of products available
type Products struct {
	Products []Product `json:"products"`
}

//Product struct is an abstraction for the different items that our store could sell.
type Product struct {
	Name            string            `json:"name"`
	ContainArticles []ContainArticles `json:"contain_articles"`
}

//ContainArticles is a set of items needed to build a given product that together represent all the components to create a given product (unsigned integer of 8 bits as stock >= 0)
type ContainArticles struct {
	ArtID    int `json:"art_id,string"`
	AmountOf int `json:"amount_of,string"`
}
