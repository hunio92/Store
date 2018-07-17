package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

const hostAddr = "http://localhost:"

type Client struct {
	Name         string         `json: name`
	ShoppingList map[string]int `json: shoppinglist`
	Money        int            `json: money`
	Port         string         `json: port`
}

type Provider struct {
	Name   string             `json: name`
	Supply map[string]float32 `json: supply`
	Port   string             `json: port`
}

type ProductInfo struct {
	Amount int
	Price  float32
	Port   string
}

var clients = map[string]Client{}
var providers = map[string][]reflect.Value{}
var products = map[string]ProductInfo{}

func main() {
	go verifyAmountAndPrice()

	http.HandleFunc("/client", handleClient)
	http.HandleFunc("/provider", handleProvider)

	http.ListenAndServe(":8080", nil)
}

func handleClient(w http.ResponseWriter, r *http.Request) {
	var client Client
	decodeJSON(r, &client)

	clients[client.Port] = client

	fmt.Println("clients:", clients)
}

func handleProvider(w http.ResponseWriter, r *http.Request) {
	var provider Provider
	decodeJSON(r, &provider)

	providers[provider.Port] = reflect.ValueOf(provider.Supply).MapKeys()

	for product, price := range provider.Supply {
		if _, ok := products[product]; ok {
			if price < products[product].Price {
				products[product] = ProductInfo{Price: price, Port: provider.Port}
			}
		} else {
			products[product] = ProductInfo{Price: price, Port: provider.Port}
		}
	}
}

func verifyAmountAndPrice() {
	for {
		for product, info := range products {
			if info.Amount < 100 {
				fmt.Println("old: ", products)
				products[product] = getCheapest(product)
				fmt.Println("new: ", products)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func getCheapest(seekProduct string) ProductInfo {
	var cheapestProduct ProductInfo
	cheapestProduct.Price = math.MaxFloat32
	for port, products := range providers {
		for _, product := range products {
			if product.Interface().(string) == seekProduct {
				res := sendRequest(port, "/orders", seekProduct)
				defer res.Body.Close()
				if res.StatusCode == http.StatusOK {
					bodyBytes, _ := ioutil.ReadAll(res.Body)
					bodyString := string(bodyBytes)
					price, err := strconv.ParseFloat(bodyString, 32)
					handleError(err, "Could not convert price to float: ")

					if float32(price) < cheapestProduct.Price {
						cheapestProduct.Price = float32(price)
						cheapestProduct.Port = port
					}
				}
			}
		}
	}

	return cheapestProduct
}

func sendRequest(port, route, product string) *http.Response {
	jsonBytes, err := json.Marshal(product)
	handleError(err, "Could not marshal json config: ")

	jsonReader := bytes.NewReader(jsonBytes)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Post(hostAddr+port+route, "json", jsonReader)
	handleError(err, "Could not send data to address: ")

	return res
}

func decodeJSON(r *http.Request, container interface{}) {
	rawJSON := json.NewDecoder(r.Body)
	err := rawJSON.Decode(container)
	handleError(err, "Failed to decode the file: ")

	defer r.Body.Close()
}

func handleError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
	}
}
