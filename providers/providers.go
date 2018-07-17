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
	Supply map[string]float32 `json: supply`
	Port   string             `json: port`
}

var prdName string
var providerInfo Provider

func main() {
	flag.StringVar(&prdName, "prd", "cl1", "Input a provider from configs")
	flag.Parse()

	getProductsInfo()

	sendRequest(providerInfo)

	fmt.Println("Port: ", providerInfo.Port)
	http.HandleFunc("/orders", handlerOrders)
	log.Fatal(http.ListenAndServe(":"+providerInfo.Port, nil))
}

func handlerOrders(w http.ResponseWriter, r *http.Request) {
	var product string
	decodeJSON(r, &product)

	getProductsInfo()

	fmt.Println("product, price: ", product, providerInfo.Supply[product])
	jData, err := json.Marshal(providerInfo.Supply[product])
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
