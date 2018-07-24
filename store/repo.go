package store

import "sync"

const productMinStock = 100

type repository struct {
	clients   map[string]Client
	providers []Provider
	products  map[string]ProductInfo
	stock     map[string][]ProductStock
	m         sync.Mutex
	stockChan chan ProductStock
}

func NewInMemoryRepository() *repository {
	repo := &repository{
		clients:   make(map[string]Client),
		providers: make([]Provider, 0),
		products:  make(map[string]ProductInfo),
		stockChan: make(chan ProductStock),
	}

	go repo.WatchStockChanges()

	return repo
}

func (r *repository) AddProvider(p Provider) {
	r.providers = append(r.providers, p)
}

func (r *repository) AddProduct(p ProductInfo) {
	if prod, ok := r.products[p.Name]; !ok {
		p.Price = p.Price * 2.5
		r.products[p.Name] = p
	}
}

func (r *repository) WatchStockChanges() {
	for p := range r.stockChan {
		r.m.Lock()
		r.stock[p.Name] = append(r.stock[p.Name], p)
		r.m.Unlock()
	}
}

func (r *repository) GetProductsBelowStock() []ProductInfo {
	retval := make([]ProductInfo, 0)
	for product, info := range r.products {
		if info.Amount < productMinStock {
			retval := append(retval, info)
		}
	}
	// fore produ
	// if productStock < productMinStock
	// add to retval

	return retval
}
