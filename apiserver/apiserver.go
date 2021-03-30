package apiserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/mpenate/stokkolm/engine"
	"go.mongodb.org/mongo-driver/bson"
)

// StartServer runs an http rest server
func StartServer() {
	e := echo.New()
	e.GET("/stock", getMaxOrder)
	e.POST("/sell", postOrder)
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	e.Logger.Fatal(e.Start(":1323"))
}

// postOrder implements the logic for the router /order
func postOrder(c echo.Context) error {
	amount, err := strconv.Atoi(c.FormValue("amount"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Invalid amount")
	}
	productName := c.FormValue("product")
	available, err := engine.GetMaxProductByName(productName)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Could not handle the request for product %s.", productName))
	}
	if available < amount {
		return c.String(http.StatusNotAcceptable, fmt.Sprintf("Impossible to make an order. Max stock %d. Required %d.", available, amount))
	} else {
		err = engine.RemoveStock(productName, amount)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, bson.M{"name": productName, "available_amount": available})
	}
}

// getMaxOrder gives the max amount that can be built
func getMaxOrder(c echo.Context) error {
	stockResponse := engine.RetrieveFullStock()
	return c.JSONPretty(http.StatusOK, stockResponse, " ")
}
