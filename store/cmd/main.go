package main

import (
	"Store/store"
)

func main() {
	repo := store.NewInMemoryRepository()

	storeService := store.NewService(repo)

	storeService.StartServer(":8080")
}
