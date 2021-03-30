package dbconnect

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/mpenate/stokkolm/schemas"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoURL      = getEnvOrDefault("MONGO_URL", "mongodb://mongo:27017/")
	inventoryPath = getEnvOrDefault("INVENTORY_PATH", "resources/inventory.json")
	productsPath  = getEnvOrDefault("PRODUCTS_PATH", "resources/products.json")
)

//InitializeDB sets the db on startup for demo purposes
func InitializeDB() {
	log.Printf("Initializing DB at %s with files %s and %s", mongoURL, inventoryPath, productsPath)
	initializeInventory()
	initializeProducts()
}

//GetProductByName uses a name to query the db
func GetProductByName(name string) (schemas.Product, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Printf("Error starting mongoClient: %s", err.Error())
	}

	mongoResp := client.Database("testing").Collection("products").FindOne(ctx, bson.M{"name": name})
	if err != nil {
		log.Printf("Error getting product %s", err.Error())
	}
	var result schemas.Product

	err = mongoResp.Decode(&result)
	if err != nil {
		log.Printf("Error decoding product %s", err.Error())
		return schemas.Product{}, err
	}
	return result, nil
}

//GetStock returns the components stock stored in mongodb
func GetStock(item int) int {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Printf("Error starting mongoClient: %s", err.Error())
	}

	cursor, err := client.Database("testing").Collection("products").Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error getting products: %s", err.Error())
	}
	var plist []schemas.Product

	for cursor.Next(ctx) {
		var result schemas.Product
		err = cursor.Decode(&result)
		if err != nil {
			log.Printf("Error getting items for product: %s", err.Error())
		}
		plist = append(plist, result)
	}
	return plist
}

//RemoveProductComponents deletes the given amount of components from inventory stock
func RemoveProductComponents(prod schemas.Product, amount int) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Printf("Error starting mongoClient: %s", err.Error())
		return err
	}
	collection := client.Database("testing").Collection("inventory")
	for _, article := range prod.ContainArticles {
		newAmount := GetStock(article.ArtID) - (article.AmountOf * amount)
		if newAmount < 0 {
			return errors.Errorf("Impossible to assume that amount of elements")
		}
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"artid": article.ArtID},
			bson.D{{"$set", bson.D{{"stock", newAmount}}}})
		if err != nil {
			log.Printf("Error removing elements: %s" + err.Error())
			return err
		}
	}
	return nil
}

func initializeInventory() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	collection := client.Database("testing").Collection("inventory")

	stockfile, _ := ioutil.ReadFile(inventoryPath)
	stockdata := schemas.Inventory{}

	_, err = collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Warning: Error during inventory clean up.  If it is the first load, ignore this message. Error %s", err.Error())
	}
	err = json.Unmarshal([]byte(stockfile), &stockdata)
	if err != nil {
		log.Printf("Warning: Error during unmarshall inventory clean up response. If it is the first load, ignore this message. Error: %s", err.Error())
	}

	for i := 0; i < len(stockdata.Articles); i++ {
		_, err := collection.InsertOne(ctx, stockdata.Articles[i])
		if err != nil {
			log.Fatalf("Insert stock ERROR: %s", err.Error())
		}
	}
}

func initializeProducts() {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))

	collection := client.Database("testing").Collection("products")
	_, err = collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Warning: Error during products cleaning. If it is the first load, ignore this message. Error: %s", err.Error())
	}

	prodsfile, _ := ioutil.ReadFile(productsPath)
	proddata := schemas.Products{}
	err = json.Unmarshal([]byte(prodsfile), &proddata)
	if err != nil {
		log.Printf("Warning: Error during products unmarshalling. If it is the first load, ignore this message. Error: %s", err.Error())
	}
	for i := 0; i < len(proddata.Products); i++ {
		_, err := collection.InsertOne(ctx, proddata.Products[i])
		if err != nil {
			log.Fatalf("Insert product ERROR: %s", err.Error())
		}
	}
}

func getEnvOrDefault(key string, defaultValue string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defaultValue
	}
	return val
}
