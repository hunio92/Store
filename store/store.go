package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Name         string         `json: name`
	ShoppingList map[string]int `json: shoppinglist`
	Money        int            `json: money`
	Port         string         `json: port`
}

type Provider struct {
	Name     string             `json: name`
	Supply   map[string]float64 `json: supply`
	Address  string             `json: port`
	Products []string
}

type ProductInfo struct {
	Name             string
	Amount           int
	Price            float64
	BestProviderAddr string
}

type ProductStock struct {
	ProductInfo
	PretAchizitie float64
	Expires       time.Time
	Qty           int
}

type service struct {
	repo *repository
}

func NewService(repo *repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) RunDaemons() {
	go s.VerifyAmount()
	// go s.VerifyPrices()
	// go s.VerifyExpiration()
}

func (s *service) VerifyAmount() {

	tick := time.Tick(5 * time.Second)

	for {

		select {
		case <-tick:
			for _, product := range s.repo.GetProductsBelowStock() {
				productOrdered := orderProduct(product)
				//order
				s.repo.stockChan <- productOrdered
				//s.repo.AddStock(ProductStock{})
			}
		}

	}
}

// func getCheapest(seekProduct string) ProductInfo {
// 	var cheapestProduct ProductInfo
// 	cheapestProduct.Price = math.MaxFloat32
// 	product := products[seekProduct]

// 	for _, provider := range respository.GetProviders() {

// 		res := sendRequest(provider.Address, "/price", ProductInfo{Name: seekProduct})
// 		defer res.Body.Close()
// 		if res.StatusCode == http.StatusOK {
// 			var order ProductInfo
// 			body, _ := ioutil.ReadAll(res.Body)
// 			json.Unmarshal(body, &order)
// 			if order.Price < cheapestProduct.Price {
// 				cheapestProduct.Price = order.Price
// 				cheapestProduct.BestProvider = provider.Address
// 			}
// 		}
// 	}

// 	return cheapestProduct
// }

func (s *service) orderProduct(product ProductInfo) ProductStock {
	const hostAddr = "http://localhost:"
	const route = "/orders"
	jsonBytes, err := json.Marshal(product)
	handleError(err, "Could not marshal json config: ")

	jsonReader := bytes.NewReader(jsonBytes)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	addr := s.getProviderAddress(product)
	res, err := client.Post(hostAddr+addr+route, "json", jsonReader)
	handleError(err, "Could not send data to address: ")

	return res
}

func (s *service) getProviderAddress(product ProductInfo) string {
	var addr string
	for _, provider := range s.repo.providers {
		if provider.Name == product.BestProvider {
			addr = provider.Address
			break
		}
	}
	return addr
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
