package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const storeAdress = "http://localhost:8080/client"

type Client struct {
	Name         string         `json: name`
	ShoppingList map[string]int `json: shoppinglist`
	Money        int            `json: money`
	Port         string         `json: port`
}

func main() {
	var cl string
	flag.StringVar(&cl, "cl", "cl1", "Input a client from configs")
	flag.Parse()

	viper.SetConfigName(cl)
	viper.AddConfigPath("./clientsConfig/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var clientConfigs Client
	err = viper.Unmarshal(&clientConfigs)
	handleError(err, "Could not unmarshal config file: ")

	sendRequest(clientConfigs)
}

func sendRequest(clientConfigs Client) {
	jsonBytes, err := json.Marshal(clientConfigs)
	handleError(err, "Could not marshal json config: ")

	jsonReader := bytes.NewReader(jsonBytes)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	client.Post(storeAdress, "json", jsonReader)
}

func handleError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
	}
}
