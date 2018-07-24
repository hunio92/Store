package store

import (
	"net/http"
)

func (s *service) StartServer(addr string) {
	s.RunDaemons()
	http.HandleFunc("/client", s.handleClient)
	http.HandleFunc("/provider", s.handleProvider)

	http.ListenAndServe(addr, nil)
}

func (s *service) handleClient(w http.ResponseWriter, r *http.Request) {
	var client Client
	decodeJSON(r, &client)

	// blah

}

func (s *service) handleProvider(w http.ResponseWriter, r *http.Request) {
	var provider Provider
	decodeJSON(r, &provider)

	s.repo.AddProvider(provider)

	for product, price := range provider.Supply {
		s.repo.AddProduct(ProductInfo{Name: product, Price: price})
	}

}
