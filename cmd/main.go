package main

import (
	"github.com/martikan/carrental_cars-api/api"
)

func main() {
	server := api.InitApi()
	server.Start()
}
