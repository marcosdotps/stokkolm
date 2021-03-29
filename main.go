package main

import (
	"github.com/mpenate/stokkolm/apiserver"
	"github.com/mpenate/stokkolm/dbconnect"
)

func main() {
	dbconnect.InitializeDB()
	apiserver.StartServer()
}
