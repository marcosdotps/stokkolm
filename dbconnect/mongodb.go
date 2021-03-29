package dbconnect

import (
	"context"
	"encoding/json"	
	"io/ioutil"
	"log"
	"time"

	"github.com/mpenate/stokkolm/schemas"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


//InitializeDB sets the db on startup for demo purposes
func InitializeDB() {
	initializeInventory()

}

//GetProductByName uses a name to query the db
func GetProductByName(name string) schemas.Product{
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://172.17.0.2:27017/"))
	if err != nil {
		log.Panicf("Error starting mongoClient: %s", err.Error())
	}

	mongoResp := client.Database("testing").Collection("products").FindOne(ctx, bson.M{"name": name})
	if err != nil {
		log.Panicf("Error getting products: %s", err.Error())
	}	
	var result schemas.Product

	err = mongoResp.Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result

}

//GetStock returns the components stock stored in mongodb
func GetStock(item int) int {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://172.17.0.2:27017/"))
	if err != nil {
		log.Panicf("Error starting mongoClient: %s", err.Error())
	}

	cursor := client.Database("testing").Collection("inventory").FindOne(ctx, bson.M{"artid": item})
	if err != nil {
		log.Panicf("Error getting stock: %s", err.Error())
	}
	var result schemas.Article
	err = cursor.Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result.Stock
}

//GetAllProducts returns the products stores in mongodb
func GetAllProducts() []schemas.Product {	
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://172.17.0.2:27017/"))
	if err != nil {
		log.Panicf("Error starting mongoClient: %s", err.Error())
	}

	cursor, err := client.Database("testing").Collection("products").Find(ctx, bson.M{})
	if err != nil {
		log.Panicf("Error getting products: %s", err.Error())
	}
	var plist []schemas.Product

	for cursor.Next(ctx) {
		var result schemas.Product
		err = cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		plist = append(plist, result)
	}
		
	return plist	
}


func RemoveProductComponents(prod schemas.Product, amount int) {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://172.17.0.2:27017/"))
	if err != nil {
		log.Panicf("Error starting mongoClient: %s", err.Error())
	}	
	collection := client.Database("testing").Collection("inventory")
	for _, article := range(prod.ContainArticles){
		newAmount := GetStock(article.ArtID) - (article.AmountOf * amount)
		log.Printf("Removing %d", article.ArtID)
		ures, err := collection.UpdateOne(
			ctx, 
			bson.M{ "artid": article.ArtID }, 
			bson.D{{"$set",bson.D{{"stock", newAmount}}}})
		if err!=nil {
			log.Fatal(err)
		}
		log.Printf("Matched %v docs and updated %v documents.\n", ures.MatchedCount, ures.ModifiedCount)
	}	
}

func initializeInventory(){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://172.17.0.2:27017/"))
	collection := client.Database("testing").Collection("inventory")

	log.Println("Generating screws, tops and so on...")

	stockfile, _ := ioutil.ReadFile("resources/inventory.json")
	stockdata := schemas.Inventory{}
	log.Println("Swipe previous items...")

	_, err = collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Error during stock clean up: %s", err.Error())
	}
	err = json.Unmarshal([]byte(stockfile), &stockdata)
	if err != nil {
		log.Fatalf("Error generating stock: %s", err.Error())
	}

	for i := 0; i < len(stockdata.Articles); i++ {
		log.Printf("Cutting some %s. Added %d with id %d.\n", stockdata.Articles[i].Name, stockdata.Articles[i].Stock, stockdata.Articles[i].ArtID)
		_, err := collection.InsertOne(ctx, stockdata.Articles[i])
		if err != nil {
			log.Fatalf("Insert stock ERROR: %s", err.Error())
		}
	}
}

func initializeProducts(){
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://172.17.0.2:27017/"))

	collection := client.Database("testing").Collection("products")
	log.Println("Swipe previous products...")
	_, err = collection.DeleteMany(ctx, bson.M{})	
	if err != nil {
		log.Fatalf("Error during products clean up: %s", err.Error())
	}

	prodsfile, _ := ioutil.ReadFile("resources/products.json")
	proddata := schemas.Products{}	
	err = json.Unmarshal([]byte(prodsfile), &proddata)
	if err != nil {
		log.Fatalf("Error generating products: %s", err.Error())
	}
	for i := 0; i < len(proddata.Products); i++ {
		log.Printf("Sending %s to our stores.\n", proddata.Products[i].Name)
		_, err := collection.InsertOne(ctx, proddata.Products[i])
		if err != nil {
			log.Fatalf("Insert product ERROR: %s", err.Error())
		}
	}

}