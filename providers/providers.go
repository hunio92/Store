package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const storeAdress = "http://localhost:8080/provider"

type Provider struct {
	Name   string             `json: name`
	Supply map[string]float64 `json: supply`
	Port   string             `json: port`
}

type ProductInfo struct {
	Name   string
	Amount int
	Price  float64
	Port   string
}

var prdName string
var providerInfo Provider

func main() {
	flag.StringVar(&prdName, "prd", "prd1", "Input a provider from configs (default: prd1)")
	flag.Parse()

	getProductsInfo()
	sendRequest(providerInfo)

	fmt.Println("Port: ", providerInfo.Port)
	http.HandleFunc("/price", handlerPrice)
	http.HandleFunc("/orders", handlerOrders)
	log.Fatal(http.ListenAndServe(":"+providerInfo.Port, nil))
}

func handlerOrders(w http.ResponseWriter, r *http.Request) {
	var order ProductInfo
	decodeJSON(r, &order)

	jData, err := json.Marshal(ProductInfo{Amount: order.Amount, Name: order.Name})
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func handlerPrice(w http.ResponseWriter, r *http.Request) {
	var product ProductInfo
	decodeJSON(r, &product)

	getProductsInfo()

	jData, err := json.Marshal(ProductInfo{Price: providerInfo.Supply[product.Name]})
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func getProductsInfo() {
	viper.SetConfigName(prdName)
	viper.AddConfigPath("./providersConfig/")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.Unmarshal(&providerInfo)
}

func sendRequest(providerInfo Provider) {
	jsonBytes, err := json.Marshal(providerInfo)
	handleError(err, "Could not marshal json config: ")

	jsonReader := bytes.NewReader(jsonBytes)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	client.Post(storeAdress, "json", jsonReader)
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
