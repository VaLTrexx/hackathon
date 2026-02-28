package main

import (
	"fmt"
	"log"

	"github.com/ValTrexx/hackathon/internal/services"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env")
	}

	market, err := services.FetchMarket("IBM")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", market)
}
