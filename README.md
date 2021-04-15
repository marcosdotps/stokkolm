# Stokkolm

## About
This project is a time-boxed exercise that wants to generate an scalable and highly decoupled api that could handle stock and orders in an easy way. 

## General Overview
#

```                                                                                                   
+-------------------------------------------------------------------------------+                  
|                                                                               |                  
|                                                                               |                  
|                                   Stokkolm                                    |                  
|                                                                               |                  
|                                                                               |                  
|                                                                               |                  
|     +--------+                           +------------+      +-------------+  |                  
|     |        |        +-----------+      |            |      |             |  |                  
|     |        |        |           |      |            |      |             |  |     +---------+  
|     |        |        |           |      |            |      |             |  |     |         |  
|     |        |        |           |      |            |      |             |  |     |         |  
|     |  Main  ---------- apiserver --------   engine   --------  dbconnect  ---|------ mongodb |  
|     |        |        |           |      |            |      |             |  |     |         |  
|     |        |        |           |      |            |      |             |  |     |         |  
|     |        |        |           |      |            |      |             |  |     +---------+  
|     |        |        +-----------+      |            |      |             |  |                  
|     +--------+                           +------------+      +-------------+  |                  
|                                                                               |                  
|                                                                               |                  
|                                                                               |                  
|                                                                               |                  
|                                                                               |                  
+-------------------------------------------------------------------------------+                  
```

Stokkolm has been designed to have highly decoupled apiserver, engine and dbconnect layers so it would be easy to change any of the parts involved.

- Apiserver implements an echov4 rest server  - with  **:1323** as listening port -  that handles requests for `/stock` and `/sell`.
- Engine deals with requests complexity and translates into dbconnect orders
- DBconnect abstraction layer is to avoid to make engine dependant of the mongodb, so in the future, if mongodb is replaced, every db operation is located under dbconnect package to simplify this kind of refactors

MongoDB is a good documental DB which, talking about stock, makes life easier to deal with update queries and lookup queries and schema is flexible.

An schemas package is present to ensure that objects are used and treated consistently into our microservice.

You can also find Prometheus metrics by going to `/metrics`

## How to run
#
Just clone this repo and from project root execute

```
$ docker-compose up --build -d
```

It will generate the docker image compiling sources inside it and start a mongo server. Also the data inside schemas will be loaded by a method on application startup. Every run will delete previous data (done for demo purposes).

## Consuming the API
#
Once you have started your docker-compose stack you can query http://localhost:1323/stock and you will get the current stock.
```
curl -XGET 'http://localhost:1323/stock'
```

Afer that you can also make a POST request like this one to place an order:
```
curl -XPOST 'http://localhost:1323/sell?product=Dining%20Chair&amount=1'
```



## How to run tests
#

For testing we need to set up some env variables that will be used to ensure that our tests uses mock data and a different mongodb.

0- Make sure you have modules updated:
```shell
$ go mod download
```

1- Run mongodb in a container as we will need it for integration tests:

```shell
$ docker run -d --name mongotest mongo:4.4.4
```

2- Get your container ip:
```shell
docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mongotest
```
3- Run the tests using the given IP inside your MONGO_URL variable and setting the invetory and product paths to testing schemas:
```
INVENTORY_PATH="test_schemas/inventory.json" PRODUCTS_PATH="test_schemas/products.json" MONGO_URL="mongodb://YOURCONTAINERIP:27017/" go test -v github.com/mpenate/stokkolm/tests
```

**WARNING!**

In case you are running this in docker desktop, please make sure to expose mongo to your host interface by binding your port:

```
$ docker run -d --name mongotest -p 27017:27017 mongo:4.4.4
```
And make sure to set test config to point localhost instead of running step 2!!


## Improvements

More decoupled functionality, better testing and a lot of fancy things could be done to make this a little bit better, but I had not infinite time :D 
